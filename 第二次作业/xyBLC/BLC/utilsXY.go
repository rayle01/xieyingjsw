package BLC

import (
	"bytes"
	"encoding/binary"
	"log"
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
