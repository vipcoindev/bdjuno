package feeexcluder

import (
	"git.ooo.ua/vipcoin/ovg-chain/x/feeexcluder/types"
	juno "github.com/forbole/juno/v5/types"
)

// handleMsgDeleteAddress allows to properly handle a message
func (m *Module) handleMsgDeleteAddress(tx *juno.Tx, index int, msg *types.MsgDeleteAddress) error {
	// TODO
	return nil
}
