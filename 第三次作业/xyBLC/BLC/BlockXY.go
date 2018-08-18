package BLC

import (
	"time"
	"bytes"
	"encoding/gob"
	"log"
)

type BlockXY struct {
	//区块高度
	Height int64
	//随机数
	Nonce int64
	//时间戳
	TimeStamp int64
	//上一个区块得哈希值
	PrevBlockHash []byte
	//本区块哈希值
	Hash []byte
	//交易数据
	Data []byte
}

//创建一个新的区块
func NewBlockXY(data string,prevBlockHash []byte,height int64) *BlockXY{

	blockxy:= &BlockXY{Height:height,Nonce:0,TimeStamp:time.Now().Unix(),PrevBlockHash:prevBlockHash,Hash:nil,Data:[]byte(data)}

	pow:=NewPowXY(blockxy)
	hash,nonce:=pow.RunXY()
	blockxy.Hash = hash
	blockxy.Nonce = nonce

	return blockxy
}

func CreateGenesisBlockXY(data string) *BlockXY{
	return NewBlockXY(data,make([]byte,32,32),0)
}

//定义block的方法，用于序列化该block对象，获取[]byte
func (block *BlockXY) SerializeXY()[]byte{
	//1.创建一个buff
	var buf bytes.Buffer

	//2.创建一个编码器
	encoder:=gob.NewEncoder(&buf)

	//3.编码
	err:=encoder.Encode(block)
	if err != nil{
		log.Panic(err)
	}

	return buf.Bytes()
}

//定义一个函数，用于将[]byte，转为block对象，反序列化
func DeserializeBlockXY(blockBytes [] byte) *BlockXY{
	var block BlockXY
	//1.先创建一个reader
	reader:=bytes.NewReader(blockBytes)
	//2.创建解码器
	decoder:=gob.NewDecoder(reader)
	//3.解码
	err:=decoder.Decode(&block)
	if err != nil{
		log.Panic(err)
	}
	return &block
}