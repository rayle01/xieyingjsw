package BLC

import "fmt"

func (cli *CLIXY) CreateWalletXY(nodeID string){
	wallets:=NewWalletsXY(nodeID) //获取钱包集合
	wallets.CreateNewWalletXY(nodeID)//创建钱包对象
	fmt.Println("钱包：",wallets.WalletMap)
}

