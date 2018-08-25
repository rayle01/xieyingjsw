package BLC

import (
	"fmt"
	"os"
)

func (cli *CLIXY) SendXY(from, to, amount [] string,nodeID string) {
	bc := GetBlockChainObject(nodeID)
	if bc == nil {
		fmt.Println("没有BlockChain，无法转账。。")
		os.Exit(1)
	}
	defer bc.DB.Close()
	bc.MineNewBlockXY(from, to, amount,nodeID)
	//添加更新
	utsoSet :=&UTXOSetXY{bc}
	utsoSet.UpdateXY()
}

