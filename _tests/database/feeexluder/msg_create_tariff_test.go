package feeexluder

import (
	"testing"

	"git.ooo.ua/vipcoin/lib/filter"
	fe "git.ooo.ua/vipcoin/ovg-chain/x/feeexcluder/types"
	"github.com/brianvoe/gofakeit/v6"

	d "github.com/forbole/bdjuno/v4/_tests/database"
)

// UpdateMsgCreateTariffs() error = sql: transaction has already been committed or rolled back, wantErr false
// DeleteMsgCreateTariffs() error = internal_server_error - sql: transaction has already been committed or rolled back, wantErr false

func TestRepository_InsertToMsgCreateTariffs(t *testing.T) {
	type args struct {
		msg  []fe.MsgCreateTariffs
		hash string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] InsertToMsgCreateTariffs",
			args: args{
				msg: []fe.MsgCreateTariffs{
					{
						Creator: d.TestAddressCreator,
						Denom:   "ovg",
						Tariff: &fe.Tariff{
							Id:            0,
							Amount:        "1",
							Denom:         "stovg",
							MinRefBalance: "10000000000",
							Fees: []*fe.Fees{
								{
									AmountFrom:  "0",
									Fee:         "0.01",
									RefReward:   "0.25",
									StakeReward: "0.5",
									MinAmount:   1000,
									NoRefReward: false,
									Creator:     "ovg1dcftms3rgxvsa2pffedke7jz5np8k4lzp6pet9",
									Id:          0,
								},
								{
									AmountFrom:  "10000000000",
									Fee:         "0.009",
									RefReward:   "0.25",
									StakeReward: "0.5",
									MinAmount:   1000,
									NoRefReward: false,
									Creator:     "ovg1dcftms3rgxvsa2pffedke7jz5np8k4lzp6pet9",
									Id:          1,
								},
							},
						},
					},
					{
						Creator: d.TestAddressCreator,
						Denom:   "ovg",
						Tariff: &fe.Tariff{
							Id:            1,
							Amount:        "100000000000",
							Denom:         "stovg",
							MinRefBalance: "10000000000",
							Fees: []*fe.Fees{
								{
									AmountFrom:  "1000000000000",
									Fee:         "0.002",
									RefReward:   "0.25",
									StakeReward: "0.5",
									MinAmount:   1000,
									NoRefReward: false,
									Creator:     "ovg1dcftms3rgxvsa2pffedke7jz5np8k4lzp6pet9",
									Id:          2,
								},
								{
									AmountFrom:  "10000000000000",
									Fee:         "0.001",
									RefReward:   "0.25",
									StakeReward: "0.5",
									MinAmount:   1000,
									NoRefReward: false,
									Creator:     "ovg1dcftms3rgxvsa2pffedke7jz5np8k4lzp6pet9",
									Id:          3,
								},
							},
						},
					},
				},
				hash: gofakeit.LetterN(64),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := d.Datastore.FeeExluder.InsertToMsgCreateTariffs(tt.args.hash, tt.args.msg...); (err != nil) != tt.wantErr {
				t.Errorf("InsertToMsgCreateTariffs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_GetAllMsgCreateTariffs(t *testing.T) {
	type args struct {
		filter filter.Filter
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] GetAllMsgCreateTariffs",
			args: args{
				filter: filter.NewFilter(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity, err := d.Datastore.FeeExluder.GetAllMsgCreateTariffs(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllMsgCreateTariffs() error = %v, wantErr %v", err, tt.wantErr)
			}

			t.Logf("size: %d", len(entity))
		})
	}
}

func TestRepository_UpdateMsgCreateTariffs(t *testing.T) {
	type args struct {
		msg  []fe.MsgCreateTariffs
		id   uint64
		hash string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] UpdateMsgCreateTariffs",
			args: args{
				msg: []fe.MsgCreateTariffs{
					{
						Creator: d.TestAddressCreator,
						Denom:   "ovg",
						Tariff: &fe.Tariff{
							Id:            0,
							Amount:        "0",
							Denom:         "ovg",
							MinRefBalance: "10000000000",
							Fees: []*fe.Fees{
								{
									AmountFrom:  "0",
									Fee:         "0.01",
									RefReward:   "0.25",
									StakeReward: "0.5",
									MinAmount:   1000,
									NoRefReward: false,
									Creator:     "ovg1dcftms3rgxvsa2pffedke7jz5np8k4lzp6pet9",
									Id:          0,
								},
							},
						},
					},
				},
				id:   1,
				hash: gofakeit.LetterN(64),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := d.Datastore.FeeExluder.UpdateMsgCreateTariffs(tt.args.hash, tt.args.id, tt.args.msg...); (err != nil) != tt.wantErr {
				t.Errorf("UpdateMsgCreateTariffs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_DeleteMsgCreateTariffs(t *testing.T) {
	type args struct {
		id uint64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "[success] DeleteMsgCreateTariffs (1)",
			args: args{
				id: 1,
			},
		},
		{
			name: "[success] DeleteMsgCreateTariffs (2)",
			args: args{
				id: 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := d.Datastore.FeeExluder.DeleteMsgCreateTariffs(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("DeleteMsgCreateTariffs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
