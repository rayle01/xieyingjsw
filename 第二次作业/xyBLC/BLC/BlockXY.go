package BLC

import "time"

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
