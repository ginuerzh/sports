// sign
package controllers

import (
	"bytes"
	//"encoding/binary"
	"encoding/hex"
	//"github.com/conformal/btcec"
	"github.com/conformal/btcnet"
	"github.com/conformal/btcscript"
	"github.com/conformal/btcutil"
	"github.com/conformal/btcwire"
	"log"
)

const (
	txFee       = 100000
	SIGHASH_ALL = 1
)

func dsha256(data []byte) []byte {
	return btcwire.DoubleSha256(data)
}

func makeScriptPubKey(toAddr string) ([]byte, error) {
	addr, err := btcutil.DecodeAddress(toAddr, &btcnet.MainNetParams)
	if err != nil {
		return nil, err
	}
	log.Println("script addr:", hex.EncodeToString(addr.ScriptAddress()))
	builder := btcscript.NewScriptBuilder()
	builder.AddOp(btcscript.OP_DUP).AddOp(btcscript.OP_HASH160)
	builder.AddData(addr.ScriptAddress())
	builder.AddOp(btcscript.OP_EQUALVERIFY).AddOp(btcscript.OP_CHECKSIG)
	//script := "76" + "a9" + "14" + hex.EncodeToString(addr.ScriptAddress()) + "88" + "ac"

	return builder.Script(), nil
}

func signScript(tx *btcwire.MsgTx, idx int, subscript []byte, privKey string) ([]byte, error) {
	wif, err := btcutil.DecodeWIF(privKey)
	if err != nil {
		return nil, err
	}

	return btcscript.SignatureScript(tx, idx, subscript, SIGHASH_ALL, wif.PrivKey.ToECDSA(), wif.CompressPubKey)
}

func makeTx(outputs []output, amount, value int64, toAddr, changeAddr string) (*btcwire.MsgTx, error) {
	msgTx := btcwire.NewMsgTx()

	for _, op := range outputs {
		hash, err := btcwire.NewShaHashFromStr(op.TxHash)
		if err != nil {
			return nil, err
		}
		b, err := hex.DecodeString(op.Script)
		if err != nil {
			return nil, err
		}
		txIn := btcwire.NewTxIn(btcwire.NewOutPoint(hash, op.TxN), b)
		msgTx.AddTxIn(txIn)
	}

	script, err := makeScriptPubKey(toAddr)
	if err != nil {
		return nil, err
	}
	txOut := btcwire.NewTxOut(value, script)
	msgTx.AddTxOut(txOut)

	if amount > value {
		script, err = makeScriptPubKey(changeAddr)
		if err != nil {
			return nil, err
		}
		txOut := btcwire.NewTxOut(amount-value, script)
		msgTx.AddTxOut(txOut)
	}
	return msgTx, nil
}

func CreateRawTx(outputs []output, amount, value int64, toAddr, changeAddr string) (rawtx string, err error) {
	msgTx, err := makeTx(outputs, amount, value, toAddr, changeAddr)
	if err != nil {
		return
	}

	buffer := &bytes.Buffer{}
	if err = msgTx.BtcEncode(buffer, 1); err != nil {
		return
	}

	finalTx := msgTx.Copy()

	for i, op := range outputs {
		b, _ := hex.DecodeString(op.Script)
		scriptSig, err := signScript(msgTx, int(op.TxN), b, op.PrivKey)
		if err != nil {
			return "", err
		}

		finalTx.TxIn[i].SignatureScript = scriptSig
	}

	buffer = &bytes.Buffer{}
	if err = finalTx.BtcEncode(buffer, 1); err != nil {
		return
	}

	return hex.EncodeToString(buffer.Bytes()), nil
}
