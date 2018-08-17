package BLC

type BlockChainXY struct {
	Blocks []*BlockXY
}

func CreateBlockChainWithGenesisBlockXY(data string) *BlockChainXY {

	genesisBlock := CreateGenesisBlockXY(data)
	return &BlockChainXY{[]*BlockXY{genesisBlock}}
}

//添加区块到区块链中
func (bc *BlockChainXY) AddBlockToBlockChainXY(data string, prevBlockHash [] byte, height int64) {
	//1.根据参数的数据，创建Block
	newBlock := NewBlockXY(data, prevBlockHash, height)
	//2.将block加入blockchain
	bc.Blocks = append(bc.Blocks, newBlock)
}