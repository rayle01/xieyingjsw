package BLC

import (
	"fmt"
	"os"
)

func (cli *CLIXY) CreateBlockChainXY(address string) {
	//fmt.Println("创世区块。。。")
	CreateBlockChainWithGenesisBlockXY(address)

	//Reset
	bc :=GetBlockChainObject()
	if bc == nil{
		fmt.Println("没有数据库。。")
		os.Exit(1)
	}
	defer bc.DB.Close()
	utxoSet:=&UTXOSetXY{bc}
	utxoSet.ResetUTXOSetXY()

}