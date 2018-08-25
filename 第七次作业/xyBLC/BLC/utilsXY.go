package BLC

import (
	"bytes"
	"encoding/binary"
	"log"
	"encoding/json"
	"fmt"
	"encoding/gob"
)

/*
将一个int64的整数转为[]byte
 */
func IntToHexXY(num int64)[]byte{
	buff :=new(bytes.Buffer)
	//将二进制数据写入w
	//func Write(w io.Writer, order ByteOrder, data interface{}) error
	err:=binary.Write(buff,binary.BigEndian,num)
	if err!=nil{
		log.Panic(err)
	}
	//转为[]byte并返回
	return buff.Bytes()
}


//json解析的的函数

func JSONToArrayXY(jsonString string) []string {
	var arr [] string
	err := json.Unmarshal([]byte(jsonString), &arr)
	if err != nil {
		log.Panic(err)
	}
	return arr
}

//字节数组反转
func ReverseBytesXY(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

//将给定的字符串的命令，转为字节数组
func commandToBytesXY(command string) []byte {
	var bytes [COMMAND_LENGTHXY]byte

	for i, c := range command {
		bytes[i] = byte(c)
	}

	return bytes[:]
}

//将给定的字节数组，转为string类型的命令
func bytesToCommandXY(bytes []byte) string {
	var command []byte

	for _, b := range bytes {
		if b != 0x0 {
			command = append(command, b)
		}
	}

	return fmt.Sprintf("%s", command)
}

//将对象进行序列化
func gobEncodeXY(data interface{}) []byte {
	var buff bytes.Buffer

	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(data)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}
