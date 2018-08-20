package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"encoding/hex"
)

//定义交易的数据
type TransactionXY struct {
	//1.交易ID-->就是交易的Hash
	TxID []byte
	//2.输入
	Vins []*TxInputXY
	//3.输出
	Vouts []*TxOutputXY
}

/*
交易：
1.CoinBase交易：创世区块中
2.转账产生的普通交易：
 */

func NewCoinBaseTransaction(address string) *TransactionXY {
	txInput := &TxInputXY{[]byte{}, -1, "Genesis Data"}
	txOutput := &TxOutputXY{50, address}
	txCoinBaseTransaction := &TransactionXY{[]byte{}, []*TxInputXY{txInput}, []*TxOutputXY{txOutput}}
	//设置交易ID
	txCoinBaseTransaction.SetIDXY()
	return txCoinBaseTransaction
}

//交易ID--->根据tx，生成一个hash
func (tx *TransactionXY) SetIDXY() {
	//1.tx--->[]byte
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	//2.[]byte-->hash
	hash := sha256.Sum256(buf.Bytes())
	//3.为tx设置ID
	tx.TxID = hash[:]
}

//根据转账的信息，创建一个普通的交易
func NewSimpleTransactionXY(from, to string, amount int64,bc *BlockChainXY,txs []*TransactionXY) *TransactionXY {
	//1.定义Input和Output的数组
	var txInputs []*TxInputXY
	var txOuputs [] *TxOutputXY

	//2.创建Input
	/*
	创世区块中交易ID：c16d3ad93450cd532dcd7ef53d8f396e46b2e59aa853ad44c284314c7b9db1b4
	 */

	//获取本次转账要使用output
	total,spentableUTXO := bc.FindSpentableUTXOsXY(from,amount,txs) //map[txID]-->[]int{index}

	for txID,indexArray:=range spentableUTXO{
		txIDBytes,_:=hex.DecodeString(txID)
		for _,index:=range indexArray{
			txInput := &TxInputXY{txIDBytes, index, from}
			txInputs = append(txInputs, txInput)
		}
	}


	//idBytes, _ := hex.DecodeString("c16d3ad93450cd532dcd7ef53d8f396e46b2e59aa853ad44c284314c7b9db1b4")
	//idBytes, _ := hex.DecodeString("143d7db0d5cce24645edb2ba0b503fe15969ade0c721edfd3578cd731c563a16")
	//txInput := &TxInput{idBytes, 1, from}
	//txInputs = append(txInputs, txInput)

	//3.创建Output

	//转账
	txOutput := &TxOutputXY{amount, to}
	txOuputs = append(txOuputs, txOutput)

	//找零
	txOutput2 := &TxOutputXY{total - amount, from}
	txOuputs = append(txOuputs, txOutput2)

	//创建交易
	tx := &TransactionXY{[]byte{}, txInputs, txOuputs}

	//设置交易的ID
	tx.SetIDXY()
	return tx

}

//判断tx是否时CoinBase交易
func (tx *TransactionXY) IsCoinBaseTransactionXY() bool {

	return len(tx.Vins[0].TxID) == 0 && tx.Vins[0].Vout == -1
}
