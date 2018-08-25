package BLC

import (
	"fmt"
	"net"
	"log"
	"io/ioutil"

)

//节点的服务端启动，可以接收其他节点传来的数据。。
func startServer(nodeID string, mineAddress string) {
	//当前节点自己地址："192.168.1.100:3001"
	nodeAddressXY = fmt.Sprintf("localhost:%s", nodeID)

	//监听：地址
	listen, err := net.Listen("tcp", nodeAddressXY)
	if err != nil {
		log.Panic(err)
	}
	defer listen.Close()

	//当前的节点，如果全节点：等待
	//如果不是全节点，钱包节点，矿工节点，给全节点发送一个消息。。

	bc := GetBlockChainObject(nodeID)
	if nodeAddressXY != knowNodesXY[0] {
		//钱包节点，矿工节点。。给全节点发送一个消息。。
		//sendMessage(knowNodes[0], "我时王二狗，我的地址是："+nodeAddress)
		sendVersion(knowNodesXY[0], bc)
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Panic(err)
		}

		fmt.Println("$发送方已经链入$", conn.RemoteAddr())

		go handleConnectionXY(conn,bc)
	}

}

func handleConnectionXY(conn net.Conn,bc *BlockChainXY){
	request, err := ioutil.ReadAll(conn) ////command + 数据
	if err != nil {
		log.Panic(err)
	}

	command:=bytesToCommandXY(request[:COMMAND_LENGTHXY])

	fmt.Printf("$接收到的命令是：%s\n", command)

	switch command {
	case COMMAND_VERSIONXY://version
		//此处是处理接收到版本数据
		handleVersion(request,bc)

	case COMMAND_GETBLOCKSXY:
		//对方向看我的信息
		handleGetBlocks(request,bc)

	case COMMAND_INVXY:
		//对方发送过来的inv信息
		handleInv(request,bc)

	case COMMAND_GETDATAXY:
		//对方发来的getdata的信息：根据hahs索要对象
		handleGetData(request,bc)

	case COMMAND_BLOCKDATAXY:
		//对象 发来的时blockdata：发来真正的区块
		handleBlockData(request,bc)
	default:
		fmt.Println("$命令无法解析")
	}
	defer conn.Close()


}







