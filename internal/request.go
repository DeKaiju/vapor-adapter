package internal

type GetUnconfirmedTxReq struct {
	TxId string `json:"tx_id"`
}

type GetBlockReq struct {
	BlockHeight uint64 `json:"block_height"`
}

type GetTxReq struct {
	TxId string `json:"tx_id"`
}

type CreateAccountReq struct {
	RootXpubs []string `json:"root_xpubs"`
	Quorum    int      `json:"quorum"`
	Alias     string   `json:"alias"`
}

type ListBalanceReq struct {
	AccountId string `json:"account_id"`
}

type ListTxReq struct {
	AccountId   string `json:"account_id"`
	Detail      bool   `json:"detail"`
	From        int    `json:"from"`
	Count       int    `json:"count"`
}

type Actions struct {
	AccountId string `json:"account_id,omitempty"`
	Amount    uint64 `json:"amount"`
	AssetId   string `json:"asset_id"`
	Type      string `json:"type"`
	Address   string `json:"address,omitempty"`
}

type BuildTransactionReq struct {
	BaseTransaction interface{} `json:"base_transaction"`
	Actions         []*Actions  `json:"actions"`
	TTL             int         `json:"ttl"`
	TimeRange       int         `json:"time_range"`
}
