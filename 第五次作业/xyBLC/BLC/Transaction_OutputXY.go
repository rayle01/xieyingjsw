package BLC

import "bytes"

//定义TxOutput的结构体
type TxOutputXY struct {
	//金额
	Value int64  //金额
	//锁定脚本，也叫输出脚本，公钥，目前先理解为用户名，钥花费这笔前，必须钥先解锁脚本
	//ScriptPubKey string
	PubKeyHash [] byte//公钥哈希
}

//判断TxOutput是否时指定的用户解锁
func (txOutput *TxOutputXY) UnlockWithAddress(address string) bool{
	fullPayload :=Base58DecodeXY([]byte(address))

	pubKeyHash:= fullPayload[1:len(fullPayload)-addressCheckSumLen]

	return bytes.Compare(pubKeyHash, txOutput.PubKeyHash) == 0
}

//根据地址创建一个output对象
func NewTxOutput(value int64,address string) *TxOutputXY{
	txOutput:=&TxOutputXY{value,nil}
	txOutput.LockXY(address)
	return txOutput
}

//锁定
func (txOutput *TxOutputXY) LockXY(address string){
	fullPayload := Base58DecodeXY([]byte(address))
	txOutput.PubKeyHash = fullPayload[1:len(fullPayload)-addressCheckSumLen]
}