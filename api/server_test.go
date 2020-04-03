package api

import (
	"reflect"
	"testing"

	"vapor-adapter/internal"
	"vapor-adapter/types"
)

var s *ServerAdapter

func init() {
	var err error
	s, err = NewServerAdapter("testnet", "http://127.0.0.1:9889", "")
	if err != nil {
		return
	}
}

func TestServerAdapter_BalancesForAddress(t *testing.T) {
	type args struct {
		accountId string
	}
	tests := []struct {
		name    string
		args    args
		want    []*types.Balance
		wantErr bool
	}{
		{name: "1", args: args{accountId: "b64dc4ce-a858-458a-97ca-badf4251d060"}, want: nil, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.BalancesForAddress(tt.args.accountId)
			if (err != nil) != tt.wantErr {
				t.Errorf("BalancesForAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BalancesForAddress() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerAdapter_BuildTransaction(t *testing.T) {
	type args struct {
		accountId       string
		toAddress       string
		tokenIdentifier string
		amount          uint64
	}
	tests := []struct {
		name    string
		args    args
		want    *internal.BuildTransactionResp
		wantErr bool
	}{
		{
			name:    "1",
			args:    args{accountId: "9cc33400-4325-40c4-9438-2233e0de4c19", toAddress: "tp1qm9zkcmz5rch096stqpejza4drktmrpc5dmsfkd", tokenIdentifier: "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", amount: uint64(100000000)},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.BuildTransaction(tt.args.accountId, tt.args.toAddress, tt.args.tokenIdentifier, tt.args.amount)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildTransaction() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerAdapter_CreateAccount(t *testing.T) {
	type args struct {
		rootXPub     string
		accountAlias string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "1",
			args:    args{rootXPub: "ef605dbc27767026b4a408c6c4ecf9fe59c60e7de7c9915e2e26a1e762c19d620c40a1f26af1cd3d749b1b6c9fa2f7f0102961720748d4f8b7eaa136caba0150", accountAlias: "Bob"},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.CreateAccount(tt.args.rootXPub, tt.args.accountAlias)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CreateAccount() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerAdapter_GetBlockCount(t *testing.T) {
	tests := []struct {
		name    string
		want    uint64
		wantErr bool
	}{
		{name: "1", want: uint64(0), wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetBlockCount()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBlockCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetBlockCount() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerAdapter_GetBlockTxs(t *testing.T) {
	type args struct {
		blockNo uint64
	}
	tests := []struct {
		name    string
		args    args
		want    []*types.Tx
		wantErr bool
	}{
		{name: "1", args: args{blockNo: uint64(6876)}, want: nil, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetBlockTxs(tt.args.blockNo)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBlockTxs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBlockTxs() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerAdapter_GetRawMemPool(t *testing.T) {
	tests := []struct {
		name    string
		want    []*types.Tx
		wantErr bool
	}{
		{name: "1", want: nil, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetRawMemPool()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRawMemPool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRawMemPool() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerAdapter_GetTransaction(t *testing.T) {
	type args struct {
		txHash string
	}
	tests := []struct {
		name    string
		args    args
		want    *types.Tx
		wantErr bool
	}{
		{name: "1", args: args{txHash: "11ca540a4e57d5879a11467a93b477fc8ff427bbc674eb46cf6c5498e946fa73"}, want: nil, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetTransaction(tt.args.txHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTransaction() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerAdapter_PubkeyToAddress(t *testing.T) {
	type args struct {
		pubkey string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "1",
			args:    args{pubkey: "1c0c2c75073c438b5612005bacdcbde2352277c44a22c5a31aa35899a3369e5fe61bb70eee5c0de48bcefddca59b14162e411b5f11d1966661a25491d48fcdbf"},
			want:    "tp1q3xjrt7ahef583lckefvvhg3djngq0l3rllkkr9",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.PubkeyToAddress(tt.args.pubkey)
			if (err != nil) != tt.wantErr {
				t.Errorf("PubkeyToAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PubkeyToAddress() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerAdapter_TxsForAddress(t *testing.T) {
	type args struct {
		accountId string
		start     int
		limit     int
	}
	tests := []struct {
		name    string
		args    args
		want    []*types.Tx
		wantErr bool
	}{
		{
			name:    "1",
			args:    args{accountId: "d640b4b3-b7b5-4bbd-85d1-29f7ed643dcb", start: 1, limit: 10},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.TxsForAddress(tt.args.accountId, tt.args.start, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("TxsForAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TxsForAddress() got = %v, want %v", got, tt.want)
			}
		})
	}
}
