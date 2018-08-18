package main

import (
	"./BLC"
)


func main(){

	blockChainXY:=BLC.CreateBlockChainWithGenesisBlockXY("Genesis Block")

	clixy:=BLC.CLIXY{blockChainXY}
	clixy.RunXY()
}


