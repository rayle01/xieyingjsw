package BLC

/*
钱包节点，矿工节点

 */

var knowNodesXY = []string{"localhost:3000"}

var nodeAddressXY string //当前节点自己的地址


//记录因该同步，但尚未同步的区块的hash
var blocksArrayXY [][]byte