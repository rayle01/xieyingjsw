package BLC

import (
	"flag"
	"os"
	"log"
	"fmt"
)

type CLIXY struct {
	//	BlockChain *BlockChainXY
}


func (cli *CLIXY) RunXY() {

	/*
	Usage:
		addblock -data DATA
		printchain


	./bc printchain
		-->执行打印的功能

	 ./bc send -from '["wangergou"]' -to '["lixiaohua"]' -amount '["4"]'
	./bc send -from '["wangergou","rose"]' -to '["lixiaohua","jace"]' -amount '["4","5"]'


	 */
	isValidArgsXY()

	//1.创建flagset命令对象

	CreateBlockChainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)

	//2.设置命令后的参数对象
	flagCreateBlockChainData := CreateBlockChainCmd.String("address", "GenesisBlock", "创世区块的信息")

	flagSendFromData := sendCmd.String("from", "", "转账源地址")
	flagSendToData := sendCmd.String("to", "", "转账目标地址")
	flagSendAmountData := sendCmd.String("amount", "", "转账金额")

	flagGetBalanceData := getBalanceCmd.String("address", "", "要查询余额的账户")

	//3.解析
	switch os.Args[1] {
	case "send":
		err := sendCmd.Parse(os.Args[2:])
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
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		printUsageXY()
		os.Exit(1)

	}
	//4.根据终端输入的命令执行对应的功能
	if sendCmd.Parsed() {
		//fmt.Println("添加区块。。。",*flagAddBlockData)
		if *flagSendFromData == "" || *flagSendToData == "" || *flagSendAmountData == "" {
			fmt.Println("转账信息有误。。")
			printUsageXY()
			os.Exit(1)
		}
		//添加区块
		//cli.AddBlockToBlockChain([]*Transaction{})
		//from:=*flagSendFromData
		//to:=*flagSendToData
		//amount:=*flagSendAmountData
		from := JSONToArrayXY(*flagSendFromData)     //[]string
		to := JSONToArrayXY(*flagSendToData)         //[]string
		amount := JSONToArrayXY(*flagSendAmountData) //[]string
		//fmt.Println(from)
		//fmt.Println(to)
		//fmt.Println(amount)
		cli.SendXY(from, to, amount)
	}

	if printChainCmd.Parsed() {
		//fmt.Println("打印区块。。。")
		//cli.BlockChain.PrintChains()
		cli.PrintChainsXY()
	}

	//添加创世区块的创建
	if CreateBlockChainCmd.Parsed() {
		if *flagCreateBlockChainData == "" {
			printUsageXY()
			os.Exit(1)
		}
		cli.CreateBlockChainXY(*flagCreateBlockChainData)
	}

	if getBalanceCmd.Parsed() {
		if *flagGetBalanceData == "" {
			fmt.Println("查询地址有误。。")
			printUsageXY()
			os.Exit(1)
		}
		cli.GetBalanceXY(*flagGetBalanceData)
	}

}

//判断终端输入的参数的长度
func isValidArgsXY() {
	if len(os.Args) < 2 {
		printUsageXY()
		os.Exit(1)
	}
}

//添加程序运行的说明
func printUsageXY() {
	fmt.Println("Usage:")
	fmt.Println("\tcreateblockchain -address DATA -- 创建创世区块")
	fmt.Println("\tsend -from From -to To -amount Amount -- 转账交易")
	fmt.Println("\tprintchain -- 打印区块")
	fmt.Println("\tgetbalance -address Data -- 查询余额")
}

func (cli *CLIXY) PrintChainsXY() {
	//cli.BlockChain.PrintChains()
	bc := GetBlockChainObject() //bc{Tip,DB}
	if bc == nil {
		fmt.Println("没有BlockChain，无法打印任何数据。。")
		os.Exit(1)
	}
	defer bc.DB.Close()
	bc.PrintChainsXY()
}

func (cli *CLIXY) AddBlockToBlockChainXY(txs []*TransactionXY) {
	//cli.BlockChain.AddBlockToBlockChain(data)
	bc := GetBlockChainObject()
	if bc == nil {
		fmt.Println("没有BlockChain，无法添加新的区块。。")
		os.Exit(1)
	}
	defer bc.DB.Close()
	bc.AddBlockToBlockChainXY(txs)
}

func (cli *CLIXY) CreateBlockChainXY(address string) {
	//fmt.Println("创世区块。。。")
	CreateBlockChainWithGenesisBlockXY(address)

}

func (cli *CLIXY) SendXY(from, to, amount []string) {
	bc := GetBlockChainObject()
	if bc == nil {
		fmt.Println("没有BlockChain，无法转账。。")
		os.Exit(1)
	}
	defer bc.DB.Close()
	bc.MineNewBlockXY(from, to, amount)
}

func (cli *CLIXY) GetBalanceXY(address string) {
	bc := GetBlockChainObject()
	if bc == nil {
		fmt.Println("没有BlockChain，无法查询。。")
		os.Exit(1)
	}
	defer bc.DB.Close()
	total := bc.GetBalanceXY(address,[]*TransactionXY{})
	fmt.Printf("%s,余额是：%d\n", address, total)
}

