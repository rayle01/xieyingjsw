package BLC

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"bytes"
)

//1.定义一个钱包结构：Wallet
type WalletXY struct {
	//1.私钥
	PrivateKey ecdsa.PrivateKey
	//2.公钥
	PublickKey []byte //原始公钥
}

//step2：产生一对密钥
func newKeyPairXY() (ecdsa.PrivateKey, []byte) {
	/*
	1.根据椭圆曲线算法，产生随机私钥
	2.根据私钥，产生公钥
	椭圆：ellipse，
	曲线：curve，

	椭圆曲线加密：(ECC：ellipse curve Cryptography)，非对称加密
		加密：
			对称加密和非对称机密啊

		SECP256K1,算法

		x轴(32byte)，y轴(32byte)--->

	 */
	//椭圆加密
	curve := elliptic.P256() //根据椭圆加密算法，得到一个椭圆曲线值
	//得到私钥
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader) //*Private
	if err != nil {
		log.Panic(err)
	}

	//产生公钥
	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	return *privateKey, publicKey
}

//step3：创建钱包对象
func NewWalletXY() *WalletXY {
	privateKey, publicKey := newKeyPairXY()
	return &WalletXY{privateKey, publicKey}
}

const version = byte(0x00)
const addressCheckSumLen = 4

//step4：根据公钥获取对应的地址
func (w *WalletXY) GetAddressXY() []byte {
	/*
	1.原始公钥-->sha256-->160-->公钥哈希
	2.版本号+公钥哈希--->校验码
	3.版本号+公钥哈希+校验码--->Base58编码

	 */

	//step1：得到公钥哈希
	pubKeyHash := PubKeyHashXY(w.PublickKey)

	//step2：添加版本号：
	versionedPayload := append([]byte{version}, pubKeyHash...)

	//step3：根据versioned_payload-->两次sha256,取前4位，得到checkSum
	checkSumBytes := CheckSumXY(versionedPayload)

	//step4：拼接全部数据
	fullPayload := append(versionedPayload, checkSumBytes...)
	//fmt.Println("full_payload:", full_payload, ",len:", len(full_payload))
	//step5：Base58编码
	address := Base58EncodeXY(fullPayload)
	return address

}

/*
原始公钥-->公钥哈希
1.sha256
2.ripemd160
 */
func PubKeyHashXY(publickKey []byte) []byte {
	//1.sha256
	hasher := sha256.New()
	hasher.Write(publickKey)
	hash1 := hasher.Sum(nil)

	//2.ripemd160
	hasher2 := ripemd160.New()
	hasher2.Write(hash1)
	hash2 := hasher2.Sum(nil)

	//3.返回
	return hash2
}

//产生校验码
func CheckSumXY(payload [] byte) []byte {
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:]) //[]byte
	return secondHash[:addressCheckSumLen]
}

//校验地址是否有效：
func IsValidAddressXY(address []byte) bool {
	/*

	 */
	//step1：Base58解码
	//version+pubkeyHash+checksum
	fullPayload := Base58DecodeXY(address) //25

	//step2：获取地址中携带的checkSUm
	checkSumBytes := fullPayload[len(fullPayload)-addressCheckSumLen:]

	versionedPayload := fullPayload[:len(fullPayload)-addressCheckSumLen]

	//step3：versioned_payload，生成一次校验码
	checkSumBytes2 := CheckSumXY(versionedPayload)

	//step4：比较checkSumBytes，checkSumBytes2
	return bytes.Compare(checkSumBytes, checkSumBytes2) == 0

}
