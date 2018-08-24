package BLC

import (
	"fmt"
	"os"
)

func (cli *CLIXY) GetBalanceXY(address string) {
	bc := GetBlockChainObject()
	if bc == nil {
		fmt.Println("没有BlockChain，无法查询。。")
		os.Exit(1)
	}
	defer bc.DB.Close()
	//total := bc.GetBalanceXY(address,[]*TransactionXY{})
	utxoSet :=&UTXOSetXY{bc}
	total:=utxoSet.GetBalanceXY(address)


	fmt.Printf("%s,余额是：%d\n", address, total)
}
