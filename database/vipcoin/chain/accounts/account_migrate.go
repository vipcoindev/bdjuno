package accounts

import (
	accountstypes "git.ooo.ua/vipcoin/chain/x/accounts/types"
	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v2/database/types"
)

// SaveAccountMigrate - saves the given account migrate inside the database
func (r Repository) SaveAccountMigrate(msg *accountstypes.MsgAccountMigrate, transactionHash string) error {
	query := `INSERT INTO overgold_chain_accounts_account_migrate 
			(transaction_hash, creator, address, hash, public_key) 
		VALUES 
			(:transaction_hash, :creator, :address, :hash, :public_key)`

	if _, err := r.db.NamedExec(query, toAccountMigrateDatabase(msg, transactionHash)); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	return nil
}

// GetAccountMigrate - get the given account migrate from database
func (r Repository) GetAccountMigrate(accountFilter filter.Filter) ([]*accountstypes.MsgAccountMigrate, error) {
	query, args := accountFilter.Build(
		tableAccountMigrate,
		types.FieldCreator, types.FieldAddress,
		types.FieldHash, types.FieldPublicKey,
	)

	var result []types.DBAccountMigrate
	if err := r.db.Select(&result, query, args...); err != nil {
		return []*accountstypes.MsgAccountMigrate{}, errs.Internal{Cause: err.Error()}
	}

	migrates := make([]*accountstypes.MsgAccountMigrate, 0, len(result))
	for _, migrate := range result {
		migrates = append(migrates, toAccountMigrateDomain(migrate))
	}

	return migrates, nil
}
