package BLC

import (
	"fmt"
	"os"
)

func (cli *CLIXY) SendXY(from, to, amount []string) {
	bc := GetBlockChainObject()
	if bc == nil {
		fmt.Println("没有BlockChain，无法转账。。")
		os.Exit(1)
	}
	defer bc.DB.Close()
	bc.MineNewBlockXY(from, to, amount)
}

