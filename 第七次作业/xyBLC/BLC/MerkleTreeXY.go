package BLC

import (
	"crypto/sha256"
	"math"
)

//第一步：创建结构体对象，表示节点和树

type MerkleNodeXY struct {
	LeftNode  *MerkleNodeXY
	RightNode *MerkleNodeXY
	DataHash  []byte
}

type MerkleTreeXY struct {
	RootNode *MerkleNodeXY
}


//给一个左右节点，生成一个新的节点
func NewMerkleNodeXY(leftNode, rightNode *MerkleNodeXY, txHash []byte) *MerkleNodeXY {
	//1.创建当前的节点
	mNode := &MerkleNodeXY{}

	//2.赋值
	if leftNode == nil && rightNode == nil {
		//mNode就是个叶子节点
		hash := sha256.Sum256(txHash)
		mNode.DataHash = hash[:]
	} else {
		//mNOde是非叶子节点
		prevHash := append(leftNode.DataHash, rightNode.DataHash...)
		hash := sha256.Sum256(prevHash)
		mNode.DataHash = hash[:]
	}
	mNode.LeftNode = leftNode
	mNode.RightNode = rightNode
	return mNode
}

//生成merkleTree
/*
func NewMerkleTreeXY(txHashData [][]byte) *MerkleTreeXY{


	//1.创建一个数组，用于存储node节点XY
	var nodes []*MerkleNodeXY

	//2.判断交易量的奇偶性
	if len(txHashData) %2 !=0{
		//奇数，复制最后一个
		txHashData = append(txHashData,txHashData[len(txHashData)-1])
	}
	//3.创建一排的叶子节点
	for _,datum :=range txHashData{
		node :=NewMerkleNodeXY(nil,nil,datum)
		nodes = append(nodes,node)
	}

	//4.生成树其他的节点
	for i:=0;i<len(txHashData)/2;i++{ // 2
		var newLevel[]*MerkleNodeXY

		for j:=0;j<len(nodes);j+=2{//j=0  tx12 tx33
			node:=NewMerkleNodeXY(nodes[j],nodes[j+1],nil)
			newLevel = append(newLevel,node)

		}

		//判断newLevel的长度的奇偶性
		if len(newLevel) % 2 != 0{
			newLevel = append(newLevel,newLevel[len(newLevel)-1])
		}

		nodes = newLevel // 3
	}

	mTree:=&MerkleTreeXY{nodes[0]}

	return mTree

}
*/
/*
//生成merkleTree 其他方法 非对称二叉树 每两个节点hash得到根节点hash

func NewMerkleTreeXY(txHashData [][]byte) *MerkleTreeXY{


	//1.创建一个数组，用于存储node所有子节点
	var nodes []*MerkleNodeXY
	var rootnode *MerkleNodeXY

	//3.创建所有子节点
	for _,datum :=range txHashData{
		node :=NewMerkleNodeXY(nil,nil,datum)
		nodes = append(nodes,node)
	}
	rootnode =nodes[0]
	for i:=0;i<len(txHashData)-1;i++{
		node:=NewMerkleNodeXY(rootnode,nodes[i+1],nil)
		rootnode =node

	}
	mTree:=&MerkleTreeXY{rootnode}
	return mTree

}
*/
func NewMerkleTreeXY(txHashData [][]byte) *MerkleTreeXY{


	//创建一个数组，用于存储node节点
	var nodes []*MerkleNodeXY

	//判断交易量的奇偶性
	if len(txHashData)%2 != 0 {
		//奇数，复制最后一个
		txHashData = append(txHashData, txHashData[len(txHashData)-1])
	}
	//创建一排的叶子节点
	for _, datum := range txHashData {
		node := NewMerkleNodeXY(nil, nil, datum)
		nodes = append(nodes, node)
	}

	count := GetCircleCountXY(len(nodes))

	for i := 0; i < count; i++ {
		var newLevel []*MerkleNodeXY

		for j := 0; j < len(nodes); j += 2 { //j=0  tx12 tx33
			node := NewMerkleNodeXY(nodes[j], nodes[j+1], nil)
			newLevel = append(newLevel, node)

		}

		//判断newLevel的长度的奇偶性
		if len(newLevel)%2 != 0 {
			newLevel = append(newLevel, newLevel[len(newLevel)-1])
		}

		nodes = newLevel

	}

	mTree := &MerkleTreeXY{nodes[0]}

	return mTree

}

func GetCircleCountXY(len int) int {
	count := 0
	for {
		if int(math.Pow(2, float64(count))) >= len {
			return count
		}
		count++
	}
}