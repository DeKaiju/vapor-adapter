package common

const (
	BTM  = "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
	ETH  = "a0889e1080999e59ed552865a1d3ee677202796222141ccc3552041708aad76c"
	USDT = "4483893ef7d9aac788c0e9e49a12398c1d35a2172adb1c8dc551fd169a6f5703"
)

type TokenParam struct {
	Code    string
	Decimal uint8
}

var TokenParams = map[string]TokenParam{
	BTM:  TokenParam{Code: "BTM", Decimal: 8},
	ETH:  TokenParam{Code: "ETH", Decimal: 9},
	USDT: TokenParam{Code: "USDT", Decimal: 6},
}
