package BLC

import (
	"fmt"
	"net"
	"log"
	"io"
	"bytes"
)

//发送消息
func sendData(to string, data []byte) {
	fmt.Println("当前节点可以发送数据")
	conn, err := net.Dial("tcp", to)
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	//发送数据。。
	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}

}

/*
处理节点之间的发送的数据
 */

//发送version
func sendVersion(toAddr string, bc *BlockChainXY) {
	//1.获取当前区块链的队在高高度
	bestHeight := bc.GetBestHeight()
	//创建version对象
	version := VersionXY{NODE_VERSIONXY, bestHeight, nodeAddressXY}
	//将version序列化
	payload := gobEncodeXY(version)

	//拼接命令+数据
	request := append(commandToBytesXY(COMMAND_VERSIONXY), payload...)

	//发送
	sendData(toAddr, request)
}

func sendGetBlocks(toAddr string) {
	getBlocks := GetBlocksXY{nodeAddressXY}
	//将getBlocks序列化
	payload := gobEncodeXY(getBlocks)

	//拼接命令+数据
	request := append(commandToBytesXY(COMMAND_GETBLOCKSXY), payload...)

	//发送
	sendData(toAddr, request)

}

func sendInv(toAddr string, kind string, data [][]byte) {
	inv := InvXY{nodeAddressXY, kind, data}
	//拼接要发送的数据
	payload := gobEncodeXY(inv)

	request := append(commandToBytesXY(COMMAND_INVXY), payload...)
	sendData(toAddr, request)

}

//发送要请求数据ud命令
func sendGetData(toAddr string, kind string, hash [] byte) {
	getData := GetDataXY{nodeAddressXY, kind, hash}
	payload := gobEncodeXY(getData)
	request := append(commandToBytesXY(COMMAND_GETDATAXY), payload...)
	sendData(toAddr, request)
}

func sendBlock(toAddr string, block *BlockXY) {

	blockData := BlockDataXY{nodeAddressXY, block.SerializeXY()}
	payload := gobEncodeXY(blockData)
	request := append(commandToBytesXY(COMMAND_BLOCKDATAXY), payload...)
	sendData(toAddr, request)
}
