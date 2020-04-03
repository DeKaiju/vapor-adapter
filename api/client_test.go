package api

import (
	"reflect"
	"testing"

	"github.com/bytom/vapor/consensus"

	"vapor-adapter/types"
)

var c *ClientAdapter

func init() {
	var err error
	c, err = NewClientAdapter("testnet")
	if err != nil {
		return
	}
}

func TestClientAdapter_Deserialize(t *testing.T) {
	type args struct {
		rawTxHex string
	}
	tests := []struct {
		name    string
		args    args
		want    *types.Tx
		wantErr bool
	}{
		{
			name:    "1",
			args:    args{rawTxHex: "07010001016401628a9415d9ed4bce588b7bdb0208ccf7ed93cdf96266678eaf7e5b9545340bb362ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff7f000116001403a7ab809e80f1d26bcae51c05d3ea01d1bdd3b401000201430041ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffbda8d0ffffffff7f0116001483a69a4dfc19f489aa8aa3d33c5493871d41dc5d00013e003cffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff80c2d72f01160014d9456c6c541e2ef2ea0b00732176ad1d97b1871400"},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.Deserialize(tt.args.rawTxHex)
			if (err != nil) != tt.wantErr {
				t.Errorf("Deserialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Deserialize() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClientAdapter_UnsignedTxHash(t *testing.T) {
	type fields struct {
		netParams *consensus.Params
	}
	type args struct {
		rawUnsignedTxHex string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "1",
			args:    args{rawUnsignedTxHex: "07010001016401628a9415d9ed4bce588b7bdb0208ccf7ed93cdf96266678eaf7e5b9545340bb362ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff7f000116001403a7ab809e80f1d26bcae51c05d3ea01d1bdd3b401000201430041ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffbda8d0ffffffff7f0116001483a69a4dfc19f489aa8aa3d33c5493871d41dc5d00013e003cffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff80c2d72f01160014d9456c6c541e2ef2ea0b00732176ad1d97b1871400"},
			want:    "9160936f300c7b0e3eba18647963fbc8672cd65f0eaec799b6adb100e996e6a3",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ClientAdapter{
				netParams: tt.fields.netParams,
			}
			got, err := c.UnsignedTxHash(tt.args.rawUnsignedTxHex)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnsignedTxHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UnsignedTxHash() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClientAdapter_PubkeyToAddress(t *testing.T) {
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

			got, err := c.PubkeyToAddress(tt.args.pubkey)
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
