// sign
package controllers

import (
	"bytes"
	//"encoding/binary"
	"encoding/hex"
	"errors"
	"github.com/conformal/btcjson"
	"github.com/conformal/btcnet"
	"github.com/conformal/btcrpcclient"
	//"github.com/conformal/btcscript"
	"github.com/conformal/btcutil"
	//"github.com/conformal/btcwire"
	"log"
)

func btcRpcClient() *btcrpcclient.Client {
	cfg := &btcrpcclient.ConnConfig{
		Host:         "localhost:8110",
		User:         "btcrpc",
		Pass:         "pbtcrpc",
		DisableTLS:   true,
		HttpPostMode: true,
	}
	client, err := btcrpcclient.New(cfg, nil)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

var client = btcRpcClient()

func CreateRawTx2(outputs []output, amount, value int64, toAddr, changeAddr string) (rawtx string, err error) {
	var inputs []btcjson.TransactionInput
	var rawInputs []btcjson.RawTxInput
	var amounts = make(map[btcutil.Address]btcutil.Amount)
	var privKeys []string

	for _, op := range outputs {
		inputs = append(inputs, btcjson.TransactionInput{Txid: op.TxHash, Vout: op.TxN})
		rawInputs = append(rawInputs, btcjson.RawTxInput{
			Txid:         op.TxHash,
			Vout:         op.TxN,
			ScriptPubKey: op.Script,
		})
		privKeys = append(privKeys, op.PrivKey)
	}

	addr, err := btcutil.DecodeAddress(toAddr, &btcnet.MainNetParams)
	if err != nil {
		return
	}
	amounts[addr] = btcutil.Amount(value)
	if amount > value {
		addr, err = btcutil.DecodeAddress(changeAddr, &btcnet.MainNetParams)
		if err != nil {
			return
		}
		amounts[addr] = btcutil.Amount(amount - value)
	}

	txMsg, err := client.CreateRawTransaction(inputs, amounts)
	if err != nil {
		return
	}

	txMsg, complete, err := client.SignRawTransaction3(txMsg, rawInputs, privKeys)
	if err != nil {
		return
	}
	if !complete {
		return "", errors.New("not complete")
	}

	buffer := &bytes.Buffer{}
	if err = txMsg.BtcEncode(buffer, 1); err != nil {
		return
	}
	return hex.EncodeToString(buffer.Bytes()), nil
}
