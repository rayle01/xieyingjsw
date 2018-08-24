package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/elliptic"
	"time"
)

//定义交易的数据
type TransactionXY struct {
	//1.交易ID-->就是交易的Hash
	TxID []byte
	//2.输入
	Vins []*TxInputXY
	//3.输出
	Vouts []*TxOutputXY
}

/*
交易：
1.CoinBase交易：创世区块中
2.转账产生的普通交易：
 */

func NewCoinBaseTransaction(address string) *TransactionXY {
	txInput := &TxInputXY{[]byte{}, -1, nil, nil}
	//txOutput := &TxOutput{10, address}
	txOutput := NewTxOutput(10, address)
	txCoinBaseTransaction := &TransactionXY{[]byte{}, []*TxInputXY{txInput}, []*TxOutputXY{txOutput}}
	//设置交易ID
	txCoinBaseTransaction.SetIDXY()
	return txCoinBaseTransaction
}

//交易ID--->根据tx，生成一个hash
func (tx *TransactionXY) SetIDXY() {
	//1.tx--->[]byte
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	//2.[]byte-->hash

	txBytes:=buf.Bytes()

	allBytes:=bytes.Join([][]byte{txBytes,IntToHexXY(time.Now().Unix())},[]byte{})

	hash := sha256.Sum256(allBytes)
	//3.为tx设置ID
	tx.TxID = hash[:]
}

//根据转账的信息，创建一个普通的交易
func NewSimpleTransactionXY(from, to string, amount int64, utxoSet *UTXOSetXY, txs []*TransactionXY) *TransactionXY {
	//1.定义Input和Output的数组
	var txInputs []*TxInputXY
	var txOuputs [] *TxOutputXY

	//2.创建Input
	/*
	创世区块中交易ID：c16d3ad93450cd532dcd7ef53d8f396e46b2e59aa853ad44c284314c7b9db1b4
	 */

	//获取本次转账要使用output
	//total, spentableUTXO := bc.FindSpentableUTXOsXY(from, amount, txs) //map[txID]-->[]int{index}
	total, spentableUTXO := utxoSet.FindSpentableUTXOsXY(from, amount, txs) //map[txID]-->[]int{index}


	//获取钱包的集合：
	wallets := NewWalletsXY()
	wallet := wallets.WalletMap[from]

	for txID, indexArray := range spentableUTXO {
		txIDBytes, _ := hex.DecodeString(txID)
		for _, index := range indexArray {
			txInput := &TxInputXY{txIDBytes, index, nil, wallet.PublickKey}
			txInputs = append(txInputs, txInput)
		}
	}

	//idBytes, _ := hex.DecodeString("c16d3ad93450cd532dcd7ef53d8f396e46b2e59aa853ad44c284314c7b9db1b4")
	//idBytes, _ := hex.DecodeString("143d7db0d5cce24645edb2ba0b503fe15969ade0c721edfd3578cd731c563a16")
	//txInput := &TxInput{idBytes, 1, from}
	//txInputs = append(txInputs, txInput)

	//3.创建Output

	//转账
	//txOutput := &TxOutput{amount, to}
	txOutput := NewTxOutput(amount, to)
	txOuputs = append(txOuputs, txOutput)

	//找零
	//txOutput2 := &TxOutput{total - amount, from}
	txOutput2 := NewTxOutput(total-amount, from)
	txOuputs = append(txOuputs, txOutput2)

	//创建交易
	tx := &TransactionXY{[]byte{}, txInputs, txOuputs}

	//设置交易的ID
	tx.SetIDXY()


	//设置签名
	//bc.SignTrasanctionXY(tx,wallet.PrivateKey)
	utxoSet.BlockChian.SignTrasanctionXY(tx,wallet.PrivateKey,txs)



	return tx

}

//判断tx是否时CoinBase交易
func (tx *TransactionXY) IsCoinBaseTransactionXY() bool {

	return len(tx.Vins[0].TxID) == 0 && tx.Vins[0].Vout == -1
}

//签名
/*
签名：为了对一笔交易进行签名
	私钥：
	要获取交易的Input，引用的output，所在的之前的交易：
 */
func (tx *TransactionXY) Sign(privateKey ecdsa.PrivateKey, prevTxsmap map[string]*TransactionXY) {
	//1.判断当前tx是否时coinbase交易
	if tx.IsCoinBaseTransactionXY() {
		return
	}

	//2.获取input对应的output所在的tx，如果不存在，无法进行签名
	for _, input := range tx.Vins {
		if prevTxsmap[hex.EncodeToString(input.TxID)] == nil {
			log.Panic("当前的Input，没有找到对应的output所在的Transaction，无法签名。。")
		}
	}

	//即将进行签名:私钥，要签名的数据
	txCopy := tx.TrimmedCopyXY()

	for index, input := range txCopy.Vins {
		// input--->5566

		prevTx := prevTxsmap[hex.EncodeToString(input.TxID)]

		txCopy.Vins[index].Signature = nil                                 //仅仅是一个双重保险，保证签名一定为空
		txCopy.Vins[index].PublicKey = prevTx.Vouts[input.Vout].PubKeyHash //设置input中的publickey为对应的output的公钥哈希


		txCopy.TxID = txCopy.NewTxIDXY()//产生要签名的数据：

		//为了方便下一个input，将数据再置为空
		txCopy.Vins[index].PublicKey = nil


		//获取要交易的数据

		/*
		第一个参数
		第二个参数：私钥
		第三个参数：要签名的数据


		func Sign(rand io.Reader, priv *PrivateKey, hash []byte) (r, s *big.Int, err error)
		r + s--->sign
		input.Signatrue = sign
	 */
		r,s,err:=ecdsa.Sign(rand.Reader, &privateKey, txCopy.TxID )
		if err != nil{
			log.Panic(err)
		}

		sign:=append(r.Bytes(),s.Bytes()...)
		tx.Vins[index].Signature = sign
	}

}

//获取要签名tx的副本
/*
要签名tx中，并不是所有的数据都要作为签名数据，生成签名
txCopy = tx{签名所需要的部分数据}
TxID

Inputs
	txid,vout,sign,publickey

Outputs
	value,pubkeyhash


交易的副本中包含的数据
	包含了原来tx中的输入和输出。
		输入中：sign，publickey
 */

func (tx *TransactionXY) TrimmedCopyXY() *TransactionXY {
	var inputs [] *TxInputXY
	var outputs [] *TxOutputXY
	for _, in := range tx.Vins {
		inputs = append(inputs, &TxInputXY{in.TxID, in.Vout, nil, nil})
	}

	for _, out := range tx.Vouts {
		outputs = append(outputs, &TxOutputXY{out.Value, out.PubKeyHash})
	}

	txCopy := &TransactionXY{tx.TxID, inputs, outputs}
	return txCopy

}

func (tx *TransactionXY) SerializeXY() [] byte {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	return buf.Bytes()
}

func (tx *TransactionXY) NewTxIDXY() []byte {
	txCopy := tx
	txCopy.TxID = []byte{}
	hash := sha256.Sum256(txCopy.SerializeXY())
	return hash[:]
}



//验证交易
/*
验证的原理：
公钥 + 要签名的数据 验证 签名：rs
 */
func (tx *TransactionXY) VerifityXY(prevTxs map[string]*TransactionXY)bool{
	//1.如果时coinbase交易，不需要验证
	if tx.IsCoinBaseTransactionXY(){
		return true
	}

	//prevTxs
	for _,input:=range prevTxs{
		if prevTxs[hex.EncodeToString(input.TxID)] == nil{
			log.Panic("当前的input没有找到对应的Transaction，无法验证。。")
		}
	}

	//验证
	txCopy:= tx.TrimmedCopyXY()

	curev:= elliptic.P256() //曲线

	for index,input:=range tx.Vins{
		//原理：再次获取 要签名的数据  + 公钥哈希 + 签名
		/*
		验证签名的有效性：
		第一个参数：公钥
		第二个参数：签名的数据
		第三、四个参数：签名：r，s
		func Verify(pub *PublicKey, hash []byte, r, s *big.Int) bool
		 */
		//ecdsa.Verify()

		//获取要签名的数据
		prevTx:=prevTxs[hex.EncodeToString(input.TxID)]

		txCopy.Vins[index].Signature = nil
		txCopy.Vins[index].PublicKey = prevTx.Vouts[input.Vout].PubKeyHash
		txCopy.TxID = txCopy.NewTxIDXY() //要签名的数据

		txCopy.Vins[index].PublicKey = nil

		//获取公钥
		/*
		type PublicKey struct {
			elliptic.Curve
			X, Y *big.Int
		}
		 */

		x:=big.Int{}
		y:=big.Int{}
		keyLen:=len(input.PublicKey)
		x.SetBytes(input.PublicKey[:keyLen/2])
		y.SetBytes(input.PublicKey[keyLen/2:])



		rawPublicKey:=ecdsa.PublicKey{curev,&x,&y}


		//获取签名：

		r :=big.Int{}
		s :=big.Int{}

		signLen:=len(input.Signature)
		r.SetBytes(input.Signature[:signLen/2])
		s.SetBytes(input.Signature[signLen/2:])

		if ecdsa.Verify(&rawPublicKey,txCopy.TxID,&r,&s) == false{
			return false
		}

	}
	return true
}