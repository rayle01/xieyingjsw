package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
)

type TxOutputsXY struct {
	UTXOs []*UTXOXY
}


//序列化
func (outs *TxOutputsXY) SerializeXY()[]byte{
	var buff bytes.Buffer

	encoder:=gob.NewEncoder(&buff)

	err :=encoder.Encode(outs)
	if err != nil{
		log.Panic(err)
	}
	return buff.Bytes()
}

//反序列化
func DeserializeTxOutputsXY(data []byte) *TxOutputsXY{
	txOutputs:=TxOutputsXY{}


	reader:=bytes.NewReader(data)
	decoder:=gob.NewDecoder(reader)
	err :=decoder.Decode(&txOutputs)
	if err != nil{
		log.Panic(err)
	}
	return &txOutputs
}