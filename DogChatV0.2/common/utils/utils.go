package utils

import (
	"DogChat/common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

type Transfer struct {
}

func (this *Transfer) PkgWrite(psm *message.SendMessage, conn net.Conn) (err error) {

	//1.对结构体essage.SendMessage序列化
	data, err := json.Marshal(psm)
	if err != nil {
		fmt.Println("PkgWrite json.Marshal 1 err =", err)
		return
	}
	//2.统计结构体的字节数
	pkgLen := uint32(len(data))
	//3.发结构体字节数
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], pkgLen) //把uint32转换成[]byte
	n, err := conn.Write(buf[:])
	if n != 4 || err != nil {
		fmt.Println("PkgWrite conn.Write 1 err=", err)
	}
	//4.发送内容
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("PkgWrite conn.Write 2 err=", err)
	}
	return
}

func (this *Transfer) PkgRead(psm *message.SendMessage, conn net.Conn) (err error) {
	buf := make([]byte, 8096)
	//取出buf对应的字节数
	_, err = conn.Read(buf[:4])
	if err != nil {
		fmt.Println("PkgRead conn.Read  1 err=", err)
		return
	}
	pkgLen := binary.BigEndian.Uint32(buf[:4])
	//再读取消息本身
	n, err := conn.Read(buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		fmt.Println("数据发生年丢包不玩了")
		fmt.Println("PkgRead conn.Read  2 err=", err)
		return
	}
	//buf是json格式的SendMessage结构体
	//将buf反序列化
	err = json.Unmarshal(buf[:pkgLen], &psm)
	if err != nil {
		fmt.Println(" PkgRead json.Unmarshal err =", err)
	}
	return
}
