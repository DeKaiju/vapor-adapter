package api

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bytom/vapor/consensus"
	"github.com/bytom/vapor/errors"

	"vapor-adapter/common"
	"vapor-adapter/internal"
	"vapor-adapter/types"
)

type ServerAdapter struct {
	nodeAddr    string
	chainId     string
	netParams   *consensus.Params
	accessToken string
}

func NewServerAdapter(chainId, nodeAddr, accessToken string) (*ServerAdapter, error) {
	netParams, ok := consensus.NetParams[chainId]
	if !ok {
		return nil, errors.New(fmt.Sprintf("%s does not exist", chainId))
	}

	return &ServerAdapter{nodeAddr: nodeAddr, chainId: chainId, netParams: &netParams, accessToken: accessToken}, nil
}

func (s *ServerAdapter) PubkeyToAddress(pubkey string) (string, error) {
	clientAdapter, err := NewClientAdapter(s.chainId)
	if err != nil {
		return "", errors.Wrapf(err, "new client adapter")
	}

	return clientAdapter.PubkeyToAddress(pubkey)
}

func (s *ServerAdapter) GetBlockCount() (uint64, error) {
	url := s.nodeAddr + "/get-block-count"
	resp := &internal.GetBlockCountResp{}
	return resp.BlockCount, s.RequestVapor(url, nil, resp)
}

func (s *ServerAdapter) GetRawMemPool() ([]*types.Tx, error) {
	url := s.nodeAddr + "/list-unconfirmed-transactions"
	resp := &internal.ListUnconfirmedTxResp{}
	if err := s.RequestVapor(url, nil, resp); err != nil {
		return nil, errors.Wrapf(err, "request list unconfirmed transaction")
	}

	var txs []*types.Tx
	for _, txId := range resp.TxIds {
		tx, err := s.getUnconfirmedTx(txId)
		if err != nil {
			return nil, err
		}

		txs = append(txs, tx)
	}
	return txs, nil
}

func (s *ServerAdapter) GetBlockTxs(blockNo uint64) ([]*types.Tx, error) {
	url := s.nodeAddr + "/get-block"
	req := &internal.GetBlockReq{BlockHeight: blockNo}
	resp := &internal.GetBlockResp{}
	if err := s.RequestVapor(url, req, resp); err != nil {
		return nil, errors.Wrapf(err, "request get block")
	}

	var txs []*types.Tx
	for _, tx := range resp.Txs {
		temp := transformTx(tx)
		temp.TxHash = tx.ID
		temp.TxAt = resp.Timestamp
		txs = append(txs, temp)
	}
	return txs, nil
}

func (s *ServerAdapter) GetTransaction(txHash string) (*types.Tx, error) {
	url := s.nodeAddr + "/get-transaction"
	req := &internal.GetTxReq{TxId: txHash}
	resp := &internal.Transaction{}
	if err := s.RequestVapor(url, req, resp); err != nil {
		return nil, errors.Wrapf(err, "request get transaction")
	}

	transaction := transformTx(resp)
	transaction.TxHash = resp.TxId
	transaction.TxAt = resp.BlockTime
	return transaction, nil
}

func (s *ServerAdapter) CreateAccount(rootXPub, accountAlias string) (string, error) {
	url := s.nodeAddr + "/create-account"
	req := &internal.CreateAccountReq{RootXpubs: []string{rootXPub}, Quorum: 1, Alias: accountAlias}
	resp := &internal.CreateAccountResp{}
	if err := s.RequestVapor(url, req, resp); err != nil {
		return "", errors.Wrapf(err, "request create account")
	}

	return resp.AccountId, nil
}

func (s *ServerAdapter) BalancesForAddress(accountId string) ([]*types.Balance, error) {
	url := s.nodeAddr + "/list-balances"
	req := &internal.ListBalanceReq{AccountId: accountId}
	var resp []*internal.Balance
	if err := s.RequestVapor(url, req, &resp); err != nil {
		return nil, errors.Wrapf(err, "request list balances")
	}

	var balances []*types.Balance
	for _, b := range resp {
		tokenParams, ok := common.TokenParams[b.AssetId]
		if !ok {
			continue
		}
		balance := &types.Balance{
			Balance:         b.Amount,
			TokenIdentifier: b.AssetId,
			TokenCode:       tokenParams.Code,
			TokenDecimal:    tokenParams.Decimal,
		}
		balances = append(balances, balance)
	}
	return balances, nil
}

func (s *ServerAdapter) TxsForAddress(accountId string, start, limit int) ([]*types.Tx, error) {
	url := s.nodeAddr + "/list-transactions"
	req := &internal.ListTxReq{AccountId: accountId, Detail: true, From: start, Count: limit}
	var resp []*internal.Transaction
	if err := s.RequestVapor(url, req, &resp); err != nil {
		return nil, errors.Wrapf(err, "request list balances")
	}

	var txs []*types.Tx
	for _, tx := range resp {
		temp := transformTx(tx)
		temp.TxHash = tx.TxId
		temp.TxAt = tx.BlockTime
		txs = append(txs, temp)
	}
	return txs, nil
}

func (s *ServerAdapter) BuildTransaction(accountId, toAddress, tokenIdentifier string, amount uint64) (*internal.BuildTransactionResp, error) {
	url := s.nodeAddr + "/build-transaction"
	var actions []*internal.Actions
	spendAction := &internal.Actions{
		AccountId: accountId,
		Amount:    amount,
		AssetId:   tokenIdentifier,
		Type:      "spend_account",
	}
	actions = append(actions, spendAction)

	controlAction := &internal.Actions{
		Amount:  amount,
		AssetId: tokenIdentifier,
		Type:    "control_address",
		Address: toAddress,
	}
	actions = append(actions, controlAction)

	req := &internal.BuildTransactionReq{Actions: actions}
	resp := &internal.BuildTransactionResp{}
	if err := s.RequestVapor(url, req, resp); err != nil {
		return nil, errors.Wrapf(err, "request build transaction")
	}

	return resp, nil
}

func (s *ServerAdapter) RequestVapor(url string, req interface{}, resp interface{}) error {
	header := make(map[string]string)
	header, err := setAccessToken(header, s.accessToken)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return err
	}

	result := &internal.Response{}
	if err := common.Post(url, header, payload, result); err != nil {
		return err
	}

	if result.Status != "success" {
		return errors.New(result.ErrDetail)
	}

	return json.Unmarshal(result.Data, resp)
}

func (s *ServerAdapter) getUnconfirmedTx(txId string) (*types.Tx, error) {
	url := s.nodeAddr + "/get-unconfirmed-transaction"
	req := &internal.GetUnconfirmedTxReq{TxId: txId}
	resp := &internal.Transaction{}
	if err := s.RequestVapor(url, req, resp); err != nil {
		return nil, errors.Wrapf(err, "request get unconfirmed transaction")
	}

	transaction := transformTx(resp)
	transaction.TxHash = resp.ID
	return transaction, nil
}

func transformTx(transaction *internal.Transaction) *types.Tx {
	inputs := transformInput(transaction)
	outputs := transformOutput(transaction)
	return &types.Tx{Inputs: inputs, Outputs: outputs}
}

func transformInput(transaction *internal.Transaction) []*types.UTXO {
	var inputs []*types.UTXO
	for _, input := range transaction.Inputs {
		tokenParams, ok := common.TokenParams[input.AssetId]
		if !ok {
			//ignore invalid token
			continue
		}
		utxo := &types.UTXO{
			Address:         input.Address,
			Value:           input.Amount,
			TokenIdentifier: input.AssetId,
			TokenCode:       tokenParams.Code,
			TokenDecimal:    tokenParams.Decimal,
		}
		inputs = append(inputs, utxo)
	}
	return inputs
}

func transformOutput(transaction *internal.Transaction) []*types.UTXO {
	var outputs []*types.UTXO
	for _, output := range transaction.Outputs {
		tokenParams, ok := common.TokenParams[output.AssetId]
		if !ok {
			//ignore invalid token
			continue
		}
		utxo := &types.UTXO{
			Address:         output.Address,
			Value:           output.Amount,
			TokenIdentifier: output.AssetId,
			TokenCode:       tokenParams.Code,
			TokenDecimal:    tokenParams.Decimal,
		}
		outputs = append(outputs, utxo)
	}
	return outputs
}

func setAccessToken(header map[string]string, accessToken string) (map[string]string, error) {
	if accessToken != "" {
		splits := strings.Split(accessToken, ":")
		if len(splits) != 2 {
			return nil, common.ErrInvalidAccessToken
		}
		header["Authorization"] = "Basic " + common.BasicAuth(splits[0], splits[1])
	}
	return header, nil
}
