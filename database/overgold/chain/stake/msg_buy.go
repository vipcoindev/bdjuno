package stake

import (
	"context"
	"database/sql"
	"errors"

	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	"git.ooo.ua/vipcoin/ovg-chain/x/stake/types"

	db "github.com/forbole/bdjuno/v4/database/types"
)

// GetAllMsgBuy - method that get data from a db (overgold_stake_buy).
func (r Repository) GetAllMsgBuy(filter filter.Filter) ([]types.MsgBuyRequest, error) {
	query, args := filter.Build(tableBuy)

	var result []db.StakeMsgBuy
	if err := r.db.Select(&result, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{What: tableBuy}
		}

		return nil, errs.Internal{Cause: err.Error()}
	}
	if len(result) == 0 {
		return nil, errs.NotFound{What: tableBuy}
	}

	return toMsgBuyDomainList(result), nil
}

// InsertMsgBuy - insert a new MsgBuyRequest in a database (overgold_stake_buy).
func (r Repository) InsertMsgBuy(hash string, msgs ...types.MsgBuyRequest) error {
	if len(msgs) == 0 || hash == "" {
		return nil
	}

	tx, err := r.db.BeginTxx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	defer func() {
		_ = tx.Rollback()
	}()

	query := `
		INSERT INTO overgold_stake_buy (
			tx_hash, creator, amount
		) VALUES (
			$1, $2, $3
		) RETURNING
			id, tx_hash, creator, amount
	`

	for _, msg := range msgs {
		m, err := toMsgBuyDatabase(hash, msg)
		if err != nil {
			return err
		}

		if _, err = tx.Exec(query, m.TxHash, m.Creator, m.Amount); err != nil {
			return errs.Internal{Cause: err.Error()}
		}
	}

	return tx.Commit()
}
