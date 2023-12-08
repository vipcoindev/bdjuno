package feeexcluder

import (
	"context"
	"database/sql"
	"errors"

	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	fe "git.ooo.ua/vipcoin/ovg-chain/x/feeexcluder/types"

	"github.com/forbole/bdjuno/v4/database/types"
)

// GetAllMsgDeleteTariffs - method that get data from a db (overgold_feeexcluder_delete_tariffs).
func (r Repository) GetAllMsgDeleteTariffs(f filter.Filter) ([]fe.MsgDeleteTariffs, error) {
	q, args := f.Build(tableDeleteTariffs)

	var deleteTariffs []types.FeeExcluderDeleteTariffs
	if err := r.db.Select(&deleteTariffs, q, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{What: tableDeleteTariffs}
		}

		return nil, errs.Internal{Cause: err.Error()}
	}
	if len(deleteTariffs) == 0 {
		return nil, errs.NotFound{What: tableDeleteTariffs}
	}

	return toMsgDeleteTariffsDomainList(deleteTariffs), nil
}

// InsertToMsgDeleteTariffs - insert new data in a database (overgold_feeexcluder_delete_tariffs).
func (r Repository) InsertToMsgDeleteTariffs(hash string, dt ...fe.MsgDeleteTariffs) error {
	if len(dt) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	defer func() {
		_ = tx.Rollback()
	}()

	q := `
		INSERT INTO overgold_feeexcluder_delete_tariffs (
			tx_hash, creator, denom, tariff_id, fees_id
		) VALUES (
			$1, $2, $3, $4, $5
		) RETURNING
			id, tx_hash, creator, denom, tariff_id, fees_id
	`

	for _, t := range dt {
		m, err := toMsgDeleteTariffsDatabase(hash, 0, t)
		if err != nil {
			return errs.Internal{Cause: err.Error()}
		}

		if _, err = tx.Exec(q, m.TxHash, m.Creator, m.Denom, m.TariffID, m.FeesID); err != nil {
			return errs.Internal{Cause: err.Error()}
		}
	}

	if err = tx.Commit(); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	return nil
}

// UpdateMsgDeleteTariffs - method that deletes in a database (overgold_feeexcluder_delete_tariffs).
func (r Repository) UpdateMsgDeleteTariffs(hash string, id uint64, ut ...fe.MsgDeleteTariffs) error {
	if len(ut) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	q := `UPDATE overgold_feeexcluder_delete_tariffs SET
				 tx_hash = $1,
				 creator = $2,
            	 tariff_id = $3,
            	 denom = $4,
                 fees_id = $5
			 WHERE id = $6`

	for _, t := range ut {
		m, err := toMsgDeleteTariffsDatabase(hash, id, t)
		if err != nil {
			return errs.Internal{Cause: err.Error()}
		}

		if _, err = tx.Exec(q, m.TxHash, m.Creator, m.TariffID, m.Denom, m.FeesID, m.ID); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// DeleteMsgDeleteTariffs - method that deletes data in a database (overgold_feeexcluder_delete_tariffs).
func (r Repository) DeleteMsgDeleteTariffs(id uint64) error {
	tx, err := r.db.BeginTxx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	defer tx.Rollback()

	q := `DELETE FROM overgold_feeexcluder_delete_tariffs WHERE id IN ($1)`

	if _, err = tx.Exec(q, id); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	if err = tx.Commit(); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	return nil
}
