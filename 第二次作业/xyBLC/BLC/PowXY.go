package BLC

import (
	"math/big"
	"crypto/sha256"
	"fmt"
	"bytes"
)

const TargetBitNum  = 20 //目标哈希值0的个数

//构建一个POW结构体
type PowXY struct {
	Block *BlockXY
	Target *big.Int
}

//生成一个新的POW对象
func NewPowXY(block *BlockXY) *PowXY  {

	powxy := &PowXY{}
	powxy.Block = block
	targetxy := big.NewInt(1)
	targetxy.Lsh(targetxy,256-TargetBitNum)
	powxy.Target =targetxy

	return powxy
}

//计算哈希函数
func (pow *PowXY) RunXY() ([]byte, int64) {
	/*
	1.将block的字段属性，拼接成一个数组
	2.定义一个nonce的值：初始值为1，
	3.产生hash--->和目标hash比较，
	 */
	//A: 定义一个nonce随机数
	var nonce int64 = 1
	var hash [32]byte
	for {

		//B： 获取拼接后的字节数组
		dataBytes := pow.prepareDataXY(nonce)

		//C:产生hash
		hash = sha256.Sum256(dataBytes) //[32]byte
		fmt.Printf("\r%d,%x", nonce, hash)

		hashInt := new(big.Int)
		hashInt.SetBytes(hash[:])

		if pow.Target.Cmp(hashInt) == 1 {
			break
		}
		nonce++
	}
	return hash[:], nonce
}

//根据block的字段属性，以及传来的nonce值，拼接成一个字节数组
func (pow *PowXY) prepareDataXY(nonce int64) []byte {
	data := bytes.Join([][]byte{
		pow.Block.PrevBlockHash,
		pow.Block.Data,
		IntToHexXY(pow.Block.TimeStamp),
		IntToHexXY(pow.Block.Height),
		IntToHexXY(nonce)}, []byte{})
	return data
}

func (pow *PowXY) IsValidXY() bool {
	hashInt := new(big.Int)
	hashInt.SetBytes(pow.Block.Hash)
	return pow.Target.Cmp(hashInt) == 1
}