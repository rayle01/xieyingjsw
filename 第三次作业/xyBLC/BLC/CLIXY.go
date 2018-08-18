package BLC

import (
	"flag"
	"os"
	"log"
	"fmt"
)

type CLIXY struct {
	BlockChain *BlockChainXY
}


//修改添加区块函数，让data为用户输入的数据

func (cli *CLIXY) RunXY() {

	isValidArgsXY()

	adddata := "hello world"
	//1.创建flagset命令对象
	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	CreateBlockChainCmd:=flag.NewFlagSet("createblockchain",flag.ExitOnError)

	//2.设置命令后的参数对象
	flagAddBlockData:=addBlockCmd.String("data",adddata,"区块的数据")
	flagCreateBlockChainData:=CreateBlockChainCmd.String("data","GenesisBlock","创世区块的信息")

	//3.解析
	switch os.Args[1] {
	case "addblock":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := CreateBlockChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		printUsageXY()
		os.Exit(1)

	}
	//4.根据终端输入的命令执行对应的功能
	if addBlockCmd.Parsed() {
		//fmt.Println("添加区块。。。",*flagAddBlockData)
		if *flagAddBlockData == ""{
			printUsageXY()
			os.Exit(1)
		}
		adddata = os.Args[2]
		//添加区块
		cli.AddBlockToBlockChainXY(adddata)

	}

	if printChainCmd.Parsed() {
		//fmt.Println("打印区块。。。")
		//cli.BlockChain.PrintChains()
		cli.PrintChainsXY()
	}

	//添加创世区块的创建
	if CreateBlockChainCmd.Parsed(){
		if *flagCreateBlockChainData ==""{
			printUsageXY()
			os.Exit(1)
		}
		cli.CreateBlockChainXY(*flagCreateBlockChainData)
	}

}

//判断终端输入的参数的长度
func isValidArgsXY() {
	if len(os.Args) < 2 {
		printUsageXY()
		os.Exit(1)
	}
}

//添加说明
func printUsageXY() {
	fmt.Println("Usage:")
	fmt.Println("\tcreateblockchain -data DATA -- 创建创世区块")
	fmt.Println("\taddblock -data DATA -- 添加区块")
	fmt.Println("\tprintchain -- 打印区块")
}


func (cli *CLIXY) PrintChainsXY(){
	cli.BlockChain.PrintChainsXY()
}


func (cli *CLIXY) AddBlockToBlockChainXY(data string){
	cli.BlockChain.AddBlockToBlockChainXY(data)
}

func(cli *CLIXY) CreateBlockChainXY(data string){
	fmt.Println("创世区块。。。")

}
