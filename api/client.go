package api

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"

	vaporCommon "github.com/bytom/vapor/common"
	"github.com/bytom/vapor/consensus"
	"github.com/bytom/vapor/consensus/segwit"
	"github.com/bytom/vapor/crypto"
	"github.com/bytom/vapor/crypto/ed25519/chainkd"
	"github.com/bytom/vapor/crypto/ed25519/ecmath"
	"github.com/bytom/vapor/errors"
	vaporTypes "github.com/bytom/vapor/protocol/bc/types"

	"vapor-adapter/common"
	"vapor-adapter/types"
)

type ClientAdapter struct {
	netParams *consensus.Params
}

func NewClientAdapter(chainId string) (*ClientAdapter, error) {
	netParams, ok := consensus.NetParams[chainId]
	if !ok {
		return nil, errors.New(fmt.Sprintf("%s does not exist", chainId))
	}
	return &ClientAdapter{netParams: &netParams}, nil
}

func (c *ClientAdapter) Deserialize(rawTxHex string) (*types.Tx, error) {
	decodeTx := &vaporTypes.Tx{}
	if err := decodeTx.UnmarshalText([]byte(rawTxHex)); err != nil {
		return nil, errors.Wrap(err, "unmarshal decodeTx")
	}

	var decodedInputs []*types.UTXO
	for _, input := range decodeTx.Inputs {
		decodedInput, err := c.decodeTxInput(input)
		if err != nil {
			return nil, errors.Wrap(err, "decodeTxInput")
		}

		if decodedInput != nil {
			decodedInputs = append(decodedInputs, decodedInput)
		}
	}

	var decodedOutputs []*types.UTXO
	for _, output := range decodeTx.Outputs {
		decodedOutput, err := c.decodeTxOutput(output)
		if err != nil {
			return nil, errors.Wrap(err, "decodeTxInput")
		}

		if decodedOutput != nil {
			decodedOutputs = append(decodedOutputs, decodedOutput)
		}

	}
	return &types.Tx{Inputs: decodedInputs, Outputs: decodedOutputs}, nil
}

func (c *ClientAdapter) UnsignedTxHash(rawUnsignedTxHex string) (string, error) {
	decodeTx := &vaporTypes.Tx{}
	if err := decodeTx.UnmarshalText([]byte(rawUnsignedTxHex)); err != nil {
		return "", errors.Wrap(err, "unmarshal decodeTx")
	}
	return decodeTx.ID.String(), nil
}

func (c *ClientAdapter) PubkeyToAddress(pubkey string) (string, error) {
	var xPubs []chainkd.XPub
	xPub, err := pubkeyToXPub(pubkey)
	if err != nil {
		return "", errors.Wrap(err, "pubkeyToXPub")
	}
	xPubs = append(xPubs, *xPub)
	path := pathForAddress(1, 1, false)
	derivedXPubs := chainkd.DeriveXPubs(xPubs, path)
	derivedPKs := chainkd.XPubKeys(derivedXPubs)

	pubHash := crypto.Ripemd160(derivedPKs[0])
	address, err := vaporCommon.NewAddressWitnessPubKeyHash(pubHash, c.netParams)
	if err != nil {
		return "", errors.Wrap(err, "NewAddressWitnessPubKeyHash")
	}
	return address.String(), nil
}

func (c *ClientAdapter) decodeTxInput(input *vaporTypes.TxInput) (*types.UTXO, error) {
	assetId := input.AssetAmount().AssetId.String()
	amount := input.AssetAmount().Amount

	var address string
	var err error

	switch i := input.TypedInput.(type) {
	case *vaporTypes.SpendInput:
		if !segwit.IsP2WScript(i.ControlProgram) {
			address = "smart contract"
			break
		}
		address, err = c.scriptToAddress(i.ControlProgram)
		if err != nil {
			return nil, errors.Wrap(err, "ScriptToAddress")
		}
	case *vaporTypes.VetoInput:
		address, err = c.scriptToAddress(i.ControlProgram)
		if err != nil {
			return nil, errors.Wrap(err, "ScriptToAddress")
		}
	}
	tokenParams, ok := common.TokenParams[assetId]
	if !ok {
		//ignore invalid token
		return nil, nil
	}
	utxo := &types.UTXO{
		Address:         address,
		Value:           amount,
		TokenIdentifier: assetId,
		TokenCode:       tokenParams.Code,
		TokenDecimal:    tokenParams.Decimal,
	}
	return utxo, nil
}

func (c *ClientAdapter) decodeTxOutput(output *vaporTypes.TxOutput) (*types.UTXO, error) {
	assetId := output.AssetAmount().AssetId.String()
	amount := output.AssetAmount().Amount

	var address string
	var err error

	switch o := output.TypedOutput.(type) {
	case *vaporTypes.CrossChainOutput:
		if assetId == common.ETH || assetId == common.USDT {
			address = "0x" + hex.EncodeToString(output.ControlProgram())
		}
	case *vaporTypes.IntraChainOutput:
		address, err = c.scriptToAddress(o.ControlProgram)
		if err != nil {
			return nil, errors.Wrap(err, "ScriptToAddress")
		}
	case *vaporTypes.VoteOutput:
		address, err = c.scriptToAddress(o.ControlProgram)
		if err != nil {
			return nil, errors.Wrap(err, "ScriptToAddress")
		}
	}
	tokenParams, ok := common.TokenParams[assetId]
	if !ok {
		//ignore invalid token
		return nil, nil
	}
	utxo := &types.UTXO{
		Address:         address,
		Value:           amount,
		TokenIdentifier: assetId,
		TokenCode:       tokenParams.Code,
		TokenDecimal:    tokenParams.Decimal,
	}
	return utxo, nil

}

func (c *ClientAdapter) scriptToAddress(script []byte) (string, error) {
	isP2WPKH, isP2WSH := segwit.IsP2WPKHScript(script), segwit.IsP2WSHScript(script)
	if !isP2WPKH && !isP2WSH {
		return "smart contract", nil
	}

	segwitHash, err := segwit.GetHashFromStandardProg(script)
	if err != nil {
		return "", errors.Wrap(err, "GetHashFromStandardProg")
	}

	var address vaporCommon.Address
	switch {
	case isP2WPKH:
		if address, err = vaporCommon.NewAddressWitnessPubKeyHash(segwitHash, c.netParams); err != nil {
			return "", err
		}
	case isP2WSH:
		if address, err = vaporCommon.NewAddressWitnessScriptHash(segwitHash, c.netParams); err != nil {
			return "", err
		}
	}
	return address.EncodeAddress(), nil
}

func pathForAddress(accountIdx, addressIndex uint64, change bool) [][]byte {
	path := [][]byte{
		{0x2C, 0x00, 0x00, 0x00},
		{0x99, 0x00, 0x00, 0x00},
		{0x00, 0x00, 0x00, 0x00},
		{0x00, 0x00, 0x00, 0x00},
		{0x00, 0x00, 0x00, 0x00},
	}
	binary.LittleEndian.PutUint32(path[2], uint32(accountIdx))
	binary.LittleEndian.PutUint32(path[4], uint32(addressIndex))
	if change {
		binary.LittleEndian.PutUint32(path[3], uint32(1))
	}
	return path
}

func pubkeyToXPub(str string) (*chainkd.XPub, error) {
	validStringLen := 128
	if len(str) != validStringLen {
		return nil, common.ErrBadLenXPubStr
	}

	var xPub chainkd.XPub
	if _, err := hex.Decode(xPub[:], []byte(str)); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("decode pubkey: %s", str))
	}

	pubkey := [32]byte{}
	copy(pubkey[:], xPub[:32])
	P := ecmath.Point{}
	if _, ok := P.Decode(pubkey); !ok {
		return nil, common.ErrInvalidXPub
	}

	return &xPub, nil
}
