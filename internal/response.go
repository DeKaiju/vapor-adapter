package internal

import "encoding/json"

type Response struct {
	Status    string          `json:"status"`
	Data      json.RawMessage `json:"data"`
	ErrDetail string          `json:"error_detail"`
}

type GetBlockCountResp struct {
	BlockCount uint64 `json:"block_count"`
}

type ListUnconfirmedTxResp struct {
	TxIds []string `json:"tx_ids"`
}

type Transaction struct {
	ID        string `json:"id"`
	TxId      string `json:"tx_id"`
	BlockTime uint64 `json:"block_time"`
	Inputs    []struct {
		Address string `json:"address"`
		Amount  uint64 `json:"amount"`
		AssetId string `json:"asset_id"`
	} `json:"inputs"`
	Outputs []struct {
		Address string `json:"address"`
		Amount  uint64 `json:"amount"`
		AssetId string `json:"asset_id"`
	} `json:"outputs"`
}

type GetBlockResp struct {
	Timestamp uint64         `json:"timestamp"`
	Txs       []*Transaction `json:"transactions"`
}

type CreateAccountResp struct {
	Alias     string `json:"alias"`
	AccountId string `json:"id"`
}

type Balance struct {
	Amount  uint64 `json:"amount"`
	AssetId string `json:"asset_id"`
}

type BuildTransactionResp struct {
	RawTransaction      string                `json:"raw_transaction"`
	SigningInstructions []SigningInstructions `json:"signing_instructions"`
}

type SigningInstructions struct {
	Position          int `json:"position"`
	WitnessComponents []struct {
		Keys []struct {
			DerivationPath []string `json:"derivation_path"`
			Xpub           string   `json:"xpub"`
		} `json:"keys,omitempty"`
		Quorum     int         `json:"quorum,omitempty"`
		Signatures interface{} `json:"signatures,omitempty"`
		Type       string      `json:"type"`
		Value      string      `json:"value,omitempty"`
	} `json:"witness_components"`
}
