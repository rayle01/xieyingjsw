package BLC

import (
	"github.com/boltdb/bolt"
	"log"
	"encoding/hex"
	"fmt"
	"bytes"
)

type UTXOSetXY struct {
	BlockChian *BlockChainXY
}

const utxosettableXY = "utxoset"

//提供一个重置的功能：获取blockchain中所有的未花费utxo

func (utxoset *UTXOSetXY) ResetUTXOSetXY() {
	err := utxoset.BlockChian.DB.Update(func(tx *bolt.Tx) error {
		//如果utxoset表存在，删除
		b := tx.Bucket([]byte(utxosettableXY))
		if b != nil {
			err := tx.DeleteBucket([]byte(utxosettableXY))
			if err != nil {
				log.Panic("重置时，删除表失败。。")
			}
		}

		//创建utxoset
		b, err := tx.CreateBucket([]byte(utxosettableXY))
		if err != nil {
			log.Panic("重置时，创建表失败。。")
		}
		if b != nil {
			//将map数据--->表
			unUTXOMap := utxoset.BlockChian.FindUnspentUTXOMapXY()

			for txIDStr, outs := range unUTXOMap {
				txID, _ := hex.DecodeString(txIDStr)
				b.Put(txID, outs.SerializeXY())
			}
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

}

//查询余额
func (utxoSet *UTXOSetXY) GetBalanceXY(address string) int64 {
	utxos := utxoSet.FindUnspentUTXOsByAddressXY(address)

	var total int64

	for _, utxo := range utxos {
		total += utxo.Output.Value
	}
	return total
}

//根据指定的地址，找出所有的utxo
func (utxoSet *UTXOSetXY) FindUnspentUTXOsByAddressXY(address string) []*UTXOXY {
	var utxos []*UTXOXY
	err := utxoSet.BlockChian.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxosettableXY))
		if b != nil {
			c := b.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				txOutputs := DeserializeTxOutputsXY(v)
				for _, utxo := range txOutputs.UTXOs { //txid, index,output
					if utxo.Output.UnlockWithAddress(address) {
						utxos = append(utxos, utxo)
					}
				}
			}
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return utxos
}


//添加一个方法，用于查询要转账的utxo
func (utxoSet *UTXOSetXY) FindSpentableUTXOsXY(from string, amount int64, txs []*TransactionXY) (int64, map[string][]int) {
	var total int64
	//用于存储转账所使用utxo
	spentableUTXOMap := make(map[string][]int)
	//1.查询未打包utxo：txs
	unPackageSpentableUTXOs := utxoSet.FindUnpackeSpentableUTXOXY(from, txs)

	for _, utxo := range unPackageSpentableUTXOs {
		total += utxo.Output.Value
		txIDStr := hex.EncodeToString(utxo.TxID)
		spentableUTXOMap[txIDStr] = append(spentableUTXOMap[txIDStr], utxo.Index)
		if total >= amount {
			return total, spentableUTXOMap
		}
	}

	//2.查询utxotable，查询utxo
	err := utxoSet.BlockChian.DB.View(func(tx *bolt.Tx) error {
		//查询utxotable中，未花费的utxo
		b := tx.Bucket([]byte(utxosettableXY))
		if b != nil {
			//查询
			c := b.Cursor()
		dbLoop:
			for k, v := c.First(); k != nil; k, v = c.Next() {
				txOutputs := DeserializeTxOutputsXY(v)
				for _, utxo := range txOutputs.UTXOs {
					if utxo.Output.UnlockWithAddress(from) {
						total += utxo.Output.Value
						txIDStr := hex.EncodeToString(utxo.TxID)
						spentableUTXOMap[txIDStr] = append(spentableUTXOMap[txIDStr], utxo.Index)
						if total >= amount {
							break dbLoop
							//return nil
						}
					}
				}

			}

		}

		return nil

	})
	if err != nil {
		log.Panic(err)
	}

	return total, spentableUTXOMap
}

//查询未打包的交易中，可以使用的utxo
func (utxoSet *UTXOSetXY) FindUnpackeSpentableUTXOXY(from string, txs []*TransactionXY) []*UTXOXY {
	//存储可以使用的未花费utxo
	var unUTXOs []*UTXOXY

	//存储已经花费的input
	spentedMap := make(map[string][]int)

	for i := len(txs) - 1; i >= 0; i-- {
		//func caculate(tx *Transaction, address string, spentTxOutputMap map[string][]int, unSpentUTXOs []*UTXO) []*UTXO {
		unUTXOs = caculateXY(txs[i], from, spentedMap, unUTXOs)
	}

	return unUTXOs
}

/*
每次转账后，更新UTXOSet：
 */
func (utxoSet *UTXOSetXY) UpdateXY() {

	//获取最后一个区块,遍历该区块中的所有tx
	newBlock := utxoSet.BlockChian.IteratorXY().Next()
	//获取所有的input
	inputs := [] *TxInputXY{}
	//遍历交易，获取所有的input
	for _, tx := range newBlock.Txs {
		if !tx.IsCoinBaseTransactionXY() {
			for _, in := range tx.Vins {
				inputs = append(inputs, in)
			}
		}
	}

	fmt.Println(len(inputs)) //5

	//存储该区块中的，tx中的未花费
	outsMap := make(map[string]*TxOutputsXY)

	//3.获取所有的output
	for _, tx := range newBlock.Txs {
		utxos := []*UTXOXY{}
		for index, output := range tx.Vouts {
			isSpent := false
			//遍历inputs的数组，比较是否有intput和该output对应
			for _, input := range inputs {
				if bytes.Compare(tx.TxID, input.TxID) == 0 && index == input.Vout {
					if bytes.Compare(output.PubKeyHash, PubKeyHashXY(input.PublicKey)) == 0 {
						isSpent = true
					}
				}
			}
			if isSpent == false {
				//output未花
				utxo := &UTXOXY{tx.TxID, index, output}
				utxos = append(utxos, utxo)
			}
		}

		//utxos,
		if len(utxos) > 0 {
			txIDStr := hex.EncodeToString(tx.TxID)
			outsMap[txIDStr] = &TxOutputsXY{utxos}
		}

	}

	//删除花费了数据,添加未花费
	err := utxoSet.BlockChian.DB.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(utxosettableXY))
		if b != nil {
			//遍历inputs，删除
			for _, input := range inputs {
				txOutputsBytes := b.Get(input.TxID)
				if len(txOutputsBytes) == 0 {
					continue
				}

				txOutputs := DeserializeTxOutputsXY(txOutputsBytes)

				isNeedDelete := false

				utxos := []*UTXOXY{}

				for _, utxo := range txOutputs.UTXOs {
					if bytes.Compare(utxo.Output.PubKeyHash, PubKeyHashXY(input.PublicKey)) == 0 && input.Vout == utxo.Index {
						isNeedDelete = true
					} else {
						utxos = append(utxos, utxo)
					}
				}

				if isNeedDelete == true {
					b.Delete(input.TxID)
					if len(utxos) > 0 {
						txOutputs := &TxOutputsXY{utxos}
						b.Put(input.TxID, txOutputs.SerializeXY())
					}
				}
			}

			//遍历map，添加
			for txIDStr, txOutputs := range outsMap {
				txID, _ := hex.DecodeString(txIDStr)
				b.Put(txID, txOutputs.SerializeXY())

			}

		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}
