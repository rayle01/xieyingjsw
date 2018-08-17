package main

import (
	"./BLC"
)


func main(){

	blockChainXY:=BLC.CreateBlockChainWithGenesisBlockXY("Genesis Block")

	powxy := BLC.NewPowXY(blockChainXY.Blocks[0])

	powxy.RunXY()

}

