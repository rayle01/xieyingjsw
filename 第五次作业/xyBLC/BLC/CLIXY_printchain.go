package BLC

import (
	"fmt"
	"os"
)

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