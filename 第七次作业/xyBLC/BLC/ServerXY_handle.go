package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
	"fmt"
)

//处理接收到的version数据
func handleVersion(request [] byte, bc *BlockChainXY) {
	/*
	1.从request中获取版本的数据：[]byte
	2.反序列化--->version
	3.操作bc，获取自己的最后block的height
	4.根对方的比较
	 */

	versionBytes := request[COMMAND_LENGTHXY:]
	//反序列化
	var version VersionXY
	reader := bytes.NewReader(versionBytes)
	decoder := gob.NewDecoder(reader)
	err := decoder.Decode(&version)
	if err != nil {
		log.Panic(err)
	}

	//获取自己的block的height
	selfHeight := bc.GetBestHeight()          // 0
	foreignerBestHeight := version.BestHeight // 2

	fmt.Printf("$接收到%s,传来的区块高度:%d\n", version.AddrFrom, foreignerBestHeight)
	if selfHeight > foreignerBestHeight {
		//发送版本对象给 对方
		sendVersion(version.AddrFrom, bc)
	} else if selfHeight < foreignerBestHeight {
		fmt.Println("$高度低于你的高度，同步下数据。。。")
		sendGetBlocks(version.AddrFrom)
	}

}

/*
处理接收到的getblocks
 */

func handleGetBlocks(request []byte, bc *BlockChainXY) {

	command := bytesToCommandXY(request[:COMMAND_LENGTHXY])
	getblocksBytes := request[COMMAND_LENGTHXY:]

	//进行发序列化
	var getblocks GetBlocksXY
	reader := bytes.NewReader(getblocksBytes)
	decoder := gob.NewDecoder(reader)
	err := decoder.Decode(&getblocks)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("$接收到了来自%s的数据，%s\n", getblocks.AddrFrom, command)

	//查询自己的数据表，将区块的hash拼接，发送给对象

	blocksHashes := bc.getBlocksHashes() //[][]{block0,block1,block2}
	//发送信息给对方
	sendInv(getblocks.AddrFrom, BLOCK_TYPEXY, blocksHashes)

}

/*
处理对方发来的inv信息
 */

func handleInv(request []byte, bc *BlockChainXY) {
	command := bytesToCommandXY(request[:COMMAND_LENGTHXY])

	invBytes := request[COMMAND_LENGTHXY:]

	var inv InvXY

	//反序列化
	reader := bytes.NewReader(invBytes)
	decoder := gob.NewDecoder(reader)
	err := decoder.Decode(&inv)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("$接收到了来自：%s，传来的%s\n", inv.AddrFrom, command)

	if inv.Type == BLOCK_TYPEXY {
		//发送获取对应数据的请求
		hash:= inv.Items[0] //第一个区块的hash,创世区块  block0
		sendGetData(inv.AddrFrom,BLOCK_TYPEXY,hash)

		if len(inv.Items)>=1{
			blocksArrayXY=inv.Items[1:] //[][]{block1,block2}
		}

	} else if inv.Type == TX_TYPEXY {

	}
}

/*
处理对方发来的getdata：对方向根据发来的hash值索要对应的对象
 */
func handleGetData(request []byte,bc *BlockChainXY){
	command:=bytesToCommandXY(request[:COMMAND_LENGTHXY])

	//反序列化
	var getData GetDataXY

	getDataBytes:=request[COMMAND_LENGTHXY:]

	reader:=bytes.NewReader(getDataBytes)
	decoder:=gob.NewDecoder(reader)
	err := decoder.Decode(&getData)
	if err != nil{
		log.Panic(err)
	}
	fmt.Printf("$接收到了来自：%s，传来的命令：%s\n",getData.AddrFrom,command)
	//
	if getData.Type == BLOCK_TYPEXY{
		//根据getdata中的hash找区块
		block:= bc.GetBlockByHash(getData.Hash) //block1的hash


		sendBlock(getData.AddrFrom,block)




	}else if getData.Type == TX_TYPEXY{

	}
}

/*
处理接收到block数据
 */
 func handleBlockData(request[] byte,bc *BlockChainXY){
 	command :=bytesToCommandXY(request[:COMMAND_LENGTHXY])

 	blockDataBytes:=request[COMMAND_LENGTHXY:]
 	//反序列化
 	var blockData BlockDataXY

 	reader:=bytes.NewReader(blockDataBytes)
 	decoder:=gob.NewDecoder(reader)
 	err := decoder.Decode(&blockData)
 	if err != nil{
 		log.Panic(err)
	}
	fmt.Printf("$收到来自：%s传来的命令：%s\n",blockData.AddrFrom,command)
	//取出block数据，存入自己的数据库
	blockBytes:=blockData.Block //block1

	block:=DeserializeBlockXY(blockBytes)
	//存入到本地的数据库

	bc.AddBlockXY(block)

	if len(blocksArrayXY) == 0{
		//重置：utxoset表
		utxoSet := UTXOSetXY{bc}
		utxoSet.ResetUTXOSetXY()
	}

	if len(blocksArrayXY) > 0{ //{block2}
		//发送请求，获取
		sendGetData(blockData.AddrFrom,BLOCK_TYPEXY,blocksArrayXY[0])
		blocksArrayXY = blocksArrayXY[1:] //{}
	}
 }
