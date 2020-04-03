package types

type UTXO struct {
	Address         string `json:"address,omitempty"`
	Value           uint64 `json:"value,omitempty"`
	TokenIdentifier string `json:"token_identifier,omitempty"`
	TokenCode       string `json:"token_code,omitempty"`
	TokenDecimal    uint8  `json:"token_decimal,omitempty"`
}

type Tx struct {
	TxHash  string            `json:"tx_hash,omitempty"`
	Inputs  []*UTXO           `json:"inputs"`
	Outputs []*UTXO           `json:"outputs"`
	TxAt    uint64            `json:"tx_at,omitempty"`
	Extra   map[string]string `json:"extra"`
}

type Balance struct {
	TokenCode       string `json:"token_code"`
	TokenIdentifier string `json:"token_identifier"`
	TokenDecimal    uint8  `json:"token_decimal"`
	Balance         uint64 `json:"balance"`
}
