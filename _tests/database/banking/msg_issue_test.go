package banking

import (
	"reflect"
	"testing"

	assets "git.ooo.ua/vipcoin/chain/x/assets/types"
	bankingtypes "git.ooo.ua/vipcoin/chain/x/banking/types"
	extratypes "git.ooo.ua/vipcoin/chain/x/types"
	"git.ooo.ua/vipcoin/lib/filter"
	"github.com/cosmos/cosmos-sdk/simapp"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/forbole/bdjuno/v2/database/types"
	bankingdb "github.com/forbole/bdjuno/v2/database/vipcoin/chain/banking"
)

func TestRepository_SaveMsgIssue(t *testing.T) {
	db, err := sqlx.Connect("pgx", "host=localhost port=5432 user=postgres dbname=juno password=postgres sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	// Create the codec
	codec := simapp.MakeTestEncodingConfig()

	type args struct {
		issue []*bankingtypes.MsgIssue
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				issue: []*bankingtypes.MsgIssue{
					{
						Creator: "vcg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv2g",
						Wallet:  "vcg1k2gx4u0hwk87ja3wyakne8cl5gytnz0uc27xm4",
						Asset:   assets.AssetVCG,
						Amount:  1_000_000,
						Extras: []*extratypes.Extra{
							{
								Kind: extratypes.EXTRA_KIND_COMMENT,
								Data: "order_id",
							},
							{
								Kind: extratypes.EXTRA_KIND_PHONE,
								Data: "2ef0186765859476750532d110bcaa39568491892edd086f1b810fa5c72db97e",
							},
						},
					},
					{
						Creator: "vcg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv3g",
						Wallet:  "vcg1k2gx4u0hwk87ja3wyakne8cl5gytnz0uc27xm5",
						Asset:   assets.AssetVCG,
						Amount:  1_000_0123,
						Extras: []*extratypes.Extra{
							{
								Kind: extratypes.EXTRA_KIND_COMMENT,
								Data: "order_id",
							},
							{
								Kind: extratypes.EXTRA_KIND_PHONE,
								Data: "2ef0186765859476750532d110bcaa39568491892edd086f1b810fa5c72db67e",
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bankingdb.NewRepository(db, codec.Marshaler)

			if err := r.SaveMsgIssue(tt.args.issue, ""); (err != nil) != tt.wantErr {
				t.Errorf("Repository.SaveMsgIssue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_GetMsgIssue(t *testing.T) {
	db, err := sqlx.Connect("pgx", "host=localhost port=5432 user=postgres dbname=juno password=postgres sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	// Create the codec
	codec := simapp.MakeTestEncodingConfig()

	type args struct {
		filter filter.Filter
	}
	tests := []struct {
		name    string
		args    args
		want    []*bankingtypes.MsgIssue
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				filter: filter.NewFilter().SetArgument(types.FieldCreator, "vcg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv2g"),
			},
			want: []*bankingtypes.MsgIssue{
				{
					Creator: "vcg1ljs7p2p9ae3en8knr3d3ke8srsfcj2zjvefv2g",
					Wallet:  "vcg1k2gx4u0hwk87ja3wyakne8cl5gytnz0uc27xm4",
					Asset:   assets.AssetVCG,
					Amount:  1_000_000,
					Extras: []*extratypes.Extra{
						{
							Kind: extratypes.EXTRA_KIND_COMMENT,
							Data: "order_id",
						},
						{
							Kind: extratypes.EXTRA_KIND_PHONE,
							Data: "2ef0186765859476750532d110bcaa39568491892edd086f1b810fa5c72db97e",
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bankingdb.NewRepository(db, codec.Marshaler)

			got, err := r.GetMsgIssue(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetMsgIssue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.GetMsgIssue() = %v, want %v", got, tt.want)
			}
		})
	}
}
