package BLC

import "fmt"

func (cli *CLIXY) CreateWalletXY(){
	wallets:=NewWalletsXY() //获取钱包集合
	wallets.CreateNewWalletXY()//创建钱包对象
	fmt.Println("钱包：",wallets.WalletMap)
}

