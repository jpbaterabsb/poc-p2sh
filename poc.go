package main

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

func main() {
	// config := rpcclient.ConnConfig{
	// 	Host:         "localhost:5555",
	// 	User:         "test",
	// 	Pass:         "test",
	// 	Params:       "regtest",
	// 	DisableTLS:   true,
	// 	HTTPPostMode: true,
	// }

	// p2sh, err := BuildMultiSigP2SHAddr()

	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(p2sh)

	raxTransaction, err := SpendMultiSig()

	if err != nil {
		panic(err)
	}

	fmt.Println(raxTransaction)
}

func BuildMultiSigP2SHAddr() (string, error) {

	// you can use your wif
	wifStr1 := "cVB8APhPmiMx6vmUJ2XFgg1X9PvCqS6dNbh5nu3RzXqRUbrn1LPT"
	wif1, err := btcutil.DecodeWIF(wifStr1)
	if err != nil {
		return "", err
	}
	// public key extracted from wif.PrivKey
	pk1 := wif1.PrivKey.PubKey().SerializeCompressed()

	wifStr2 := "cVgxEkRBtnfvd41ssd4PCsiemahAHidFrLWYoDBMNojUeME8dojZ"
	wif2, err := btcutil.DecodeWIF(wifStr2)
	if err != nil {
		return "", err
	}
	pk2 := wif2.PrivKey.PubKey().SerializeCompressed()

	wifStr3 := "cPXZBMz5pKytwCyUNAdq94R9VafU8L2QmAW8uw3gKrzjuCWCd3TM"
	wif3, err := btcutil.DecodeWIF(wifStr3)
	if err != nil {
		return "", nil
	}
	pk3 := wif3.PrivKey.PubKey().SerializeCompressed()

	// create redeem script for 2 of 3 multi-sig
	builder := txscript.NewScriptBuilder()
	// add the minimum number of needed signatures
	builder.AddOp(txscript.OP_2)
	// add the 3 public key
	builder.AddData(pk1).AddData(pk2).AddData(pk3)
	// add the total number of public keys in the multi-sig screipt
	builder.AddOp(txscript.OP_3)
	// add the check-multi-sig op-code
	builder.AddOp(txscript.OP_CHECKMULTISIG)
	// redeem script is the script program in the format of []byte
	redeemScript, err := builder.Script()
	if err != nil {
		return "", err
	}

	// calculate the hash160 of the redeem script
	redeemHash := btcutil.Hash160(redeemScript)

	// if using Bitcoin main net then pass &chaincfg.MainNetParams as second argument
	addr, err := btcutil.NewAddressScriptHashFromHash(redeemHash, &chaincfg.RegressionNetParams)
	if err != nil {
		return "", err
	}

	return addr.EncodeAddress(), nil
}

func SpendMultiSig() (string, error) {
	// you can use your wif
	wifStr1 := "cVB8APhPmiMx6vmUJ2XFgg1X9PvCqS6dNbh5nu3RzXqRUbrn1LPT"
	wif1, err := btcutil.DecodeWIF(wifStr1)
	if err != nil {
		return "", err
	}
	// public key extracted from wif.PrivKey
	pk1 := wif1.PrivKey.PubKey().SerializeCompressed()

	wifStr2 := "cVgxEkRBtnfvd41ssd4PCsiemahAHidFrLWYoDBMNojUeME8dojZ"
	wif2, err := btcutil.DecodeWIF(wifStr2)
	if err != nil {
		return "", err
	}
	pk2 := wif2.PrivKey.PubKey().SerializeCompressed()

	wifStr3 := "cPXZBMz5pKytwCyUNAdq94R9VafU8L2QmAW8uw3gKrzjuCWCd3TM"
	wif3, err := btcutil.DecodeWIF(wifStr3)
	if err != nil {
		return "", nil
	}
	pk3 := wif3.PrivKey.PubKey().SerializeCompressed()

	// create redeem script for 2 of 3 multi-sig
	builder := txscript.NewScriptBuilder()
	// add the minimum number of needed signatures
	builder.AddOp(txscript.OP_2)
	// add the 3 public key
	builder.AddData(pk1).AddData(pk2).AddData(pk3)
	// add the total number of public keys in the multi-sig screipt
	builder.AddOp(txscript.OP_3)
	// add the check-multi-sig op-code
	builder.AddOp(txscript.OP_CHECKMULTISIG)
	// redeem script is the script program in the format of []byte
	redeemScript, err := builder.Script()
	if err != nil {
		return "", err
	}

	redeemTx := wire.NewMsgTx(wire.TxVersion)

	// you should provide your UTXO hash
	utxoHash, err := chainhash.NewHashFromStr("1b6c1fd0b15e9b2c2c886fd11121d488cf64f6ad031722c2128696522d12fae5")

	if err != nil {
		return "", nil
	}

	// and add the index of the UTXO
	outPoint := wire.NewOutPoint(utxoHash, 0)

	txIn := wire.NewTxIn(outPoint, nil, nil)

	redeemTx.AddTxIn(txIn)

	// adding the output to tx
	decodedAddr, err := btcutil.DecodeAddress("mhRnDJY6iKpbiadqeCMs1NQAKsQwKVn5xU", &chaincfg.RegressionNetParams)
	if err != nil {
		return "", err
	}
	destinationAddrByte, err := txscript.PayToAddrScript(decodedAddr)
	if err != nil {
		return "", err
	}

	// adding the destination address and the amount to the transaction
	redeemTxOut := wire.NewTxOut(43000, destinationAddrByte)
	redeemTx.AddTxOut(redeemTxOut)

	// signing the tx

	sig1, err := txscript.RawTxInSignature(redeemTx, 0, redeemScript, txscript.SigHashAll, wif1.PrivKey)
	if err != nil {
		return "", err
	}

	//sig2, err := txscript.RawTxInSignature(redeemTx, 0, redeemScript, txscript.SigHashAll, wif2.PrivKey)
	//if err != nil {
	//	return "", err
	//}

	sig3, err := txscript.RawTxInSignature(redeemTx, 0, redeemScript, txscript.SigHashAll, wif3.PrivKey)
	if err != nil {
		fmt.Println("got error in constructing sig3")
		return "", err
	}

	signature := txscript.NewScriptBuilder()
	signature.AddOp(txscript.OP_FALSE).AddData(sig1)
	signature.AddData(sig3).AddData(redeemScript)
	signatureScript, err := signature.Script()
	if err != nil {
		// Handle the error.
		return "", err
	}

	redeemTx.TxIn[0].SignatureScript = signatureScript

	var signedTx bytes.Buffer
	redeemTx.Serialize(&signedTx)

	hexSignedTx := hex.EncodeToString(signedTx.Bytes())

	return hexSignedTx, nil
}
