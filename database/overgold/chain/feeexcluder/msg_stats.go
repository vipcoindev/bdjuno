package feeexcluder

import (
	"context"
	"database/sql"
	"errors"

	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	fe "git.ooo.ua/vipcoin/ovg-chain/x/feeexcluder/types"
	"github.com/jmoiron/sqlx"

	"github.com/forbole/bdjuno/v4/database/types"
)

// GetAllStats - method that get data from a db (overgold_feeexcluder_stats).
// TODO: use JOIN and other db model
func (r Repository) GetAllStats(f filter.Filter) ([]fe.Stats, error) {
	q, args := f.Build(tableStats)

	// 1) get stats
	var stats []types.FeeExcluderStats
	if err := r.db.Select(&stats, q, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{What: tableStats}
		}

		return nil, errs.Internal{Cause: err.Error()}
	}
	if len(stats) == 0 {
		return nil, errs.NotFound{What: tableStats}
	}

	// 2) get daily stats
	result := make([]fe.Stats, 0, len(stats))
	for _, s := range stats {
		dailyStats, err := r.GetAllDailyStats(filter.NewFilter().SetArgument(types.FieldID, s.DailyStatsID))
		if err != nil {
			return nil, err
		}
		if len(dailyStats) == 0 {
			return nil, errs.NotFound{What: tableDailyStats}
		}

		result = append(result, toStatsDomain(&dailyStats[0], s))
	}

	return result, nil
}

// InsertToStats - insert new data in a database (overgold_feeexcluder_stats).
func (r Repository) InsertToStats(tx *sqlx.Tx, stats fe.Stats) (lastID string, err error) {
	if tx == nil {
		tx, err = r.db.BeginTxx(context.Background(), &sql.TxOptions{})
		if err != nil {
			return "", errs.Internal{Cause: err.Error()}
		}

		defer commit(tx, err)
	}

	// 1) add daily stats and get unique ids
	dailyStatsID, err := r.InsertToDailyStats(tx, *stats.Stats)
	if err != nil {
		return "", err
	}

	// 2) add stats
	q := `
		INSERT INTO overgold_feeexcluder_stats (
			id, date, daily_stats_id
		) VALUES (
			$1, $2, $3
		) RETURNING	id
	`

	m, err := toStatsDatabase(dailyStatsID, stats)
	if err != nil {
		return "", errs.Internal{Cause: err.Error()}
	}

	if err = r.db.QueryRowx(q, m.ID, m.Date, m.DailyStatsID).Scan(&lastID); err != nil {
		return "", errs.Internal{Cause: err.Error()}
	}

	return lastID, nil
}

// UpdateStats - method that updates in a database (overgold_feeexcluder_stats).
func (r Repository) UpdateStats(tx *sqlx.Tx, stats fe.Stats) (err error) {
	if tx == nil {
		tx, err = r.db.BeginTxx(context.Background(), &sql.TxOptions{})
		if err != nil {
			return errs.Internal{Cause: err.Error()}
		}

		defer commit(tx, err)
	}

	// 1) update stats and get unique id for daily stats
	// 1.a) get daily stats id via stats index
	allStats, err := r.getAllStatsWithUniqueID(filter.NewFilter().SetArgument(types.FieldID, stats.Index))
	if err != nil {
		return err
	}
	if len(allStats) != 1 {
		return errs.Internal{Cause: "expected 1 stats"}
	}
	dailyStatsID := allStats[0].DailyStatsID

	// 1.b) update stats
	q := `UPDATE overgold_feeexcluder_stats SET
				 date = $1,
				 daily_stats_id = $2
			 WHERE id = $3`

	m, err := toStatsDatabase(dailyStatsID, stats)
	if err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	if _, err = tx.Exec(q, m.Date, m.DailyStatsID, m.ID); err != nil {
		return err
	}

	// 2) update daily stats
	if err = r.UpdateDailyStats(tx, dailyStatsID, *stats.Stats); err != nil {
		return err
	}

	return nil
}

// DeleteStats - method that deletes data in a database (overgold_feeexcluder_stats).
func (r Repository) DeleteStats(tx *sqlx.Tx, id string) (err error) {
	if tx == nil {
		tx, err = r.db.BeginTxx(context.Background(), &sql.TxOptions{})
		if err != nil {
			return errs.Internal{Cause: err.Error()}
		}

		defer commit(tx, err)
	}

	// 1) delete stats and get unique id via stats index
	// 1.a) get stats id via stats index
	allStats, err := r.getAllStatsWithUniqueID(filter.NewFilter().SetArgument(types.FieldID, id))
	if err != nil {
		return err
	}
	if len(allStats) != 1 {
		return errs.Internal{Cause: "expected 1 stats"}
	}
	dailyStatsID := allStats[0].DailyStatsID

	// 1.b) delete daily stats
	q := `DELETE FROM overgold_feeexcluder_stats WHERE id IN ($1)`

	if _, err = tx.Exec(q, id); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	// 2) delete daily stats
	if err = r.DeleteDailyStats(tx, dailyStatsID); err != nil {
		return err
	}

	return nil
}

// getAllStatsWithUniqueID - method that get data from a db (overgold_feeexcluder_stats).
func (r Repository) getAllStatsWithUniqueID(f filter.Filter) ([]types.FeeExcluderStats, error) {
	q, args := f.Build(tableStats)

	var stats []types.FeeExcluderStats
	if err := r.db.Select(&stats, q, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{What: tableStats}
		}

		return nil, errs.Internal{Cause: err.Error()}
	}
	if len(stats) == 0 {
		return nil, errs.NotFound{What: tableStats}
	}

	return stats, nil
}
