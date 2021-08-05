package process

import (
	"DogChat/client/model"
	"DogChat/common/message"
	"DogChat/common/utils"
	"encoding/json"
	_ "errors"
	"fmt"
	"net"
	"os"
)

func ClientRegister(userId int, userName, userPwd string) (err error) {
	//与服务器建立链接
	address := fmt.Sprintf("%s:%d", model.ServerIp, model.ServerPort)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("ClientRegister call net.Dial err=", err)
		return
	}
	// var rrm *message.RegisterMessage
	// rrm.UserId = userId
	// rrm.UserName = userName
	// rrm.UserPwd = userPwd
	rm := message.RegisterMessage{
		UserId:   userId,
		UserName: userName,
		UserPwd:  userPwd,
	}
	//对rrm序列化
	data, err := json.Marshal(rm)
	if err != nil {
		fmt.Println("ClientRegister call json.Marshal err=", err)
		return
	}
	//for test
	//fmt.Println("string(data) = ", string(data))

	// var sm *message.SendMessage
	// sm.Data = string(data)
	// sm.Type = message.RegisterMessageType
	sm := &message.SendMessage{
		Type: message.RegisterMessageType,
		Data: string(data),
	}
	//for test
	fmt.Println("sm = ", sm)
	var tf utils.Transfer
	err = tf.PkgWrite(sm, conn)
	if err != nil {
		fmt.Println("ClientRegister call tf.PkgWrite err=", err)
		return
	}
	//等待服务器的消息
	sm2 := &message.SendMessage{}
	err = tf.PkgRead(sm2, conn)
	if err != nil {
		fmt.Println("ClientRegister call tf.PkgRead err=", err)
		return
	}
	if sm2.Type != message.RegisterResMessageType {
		fmt.Println("从服务器读到的消息类型不是RegisterResMessageType,将不进行处理")
		return
	}
	//声明一个结构体
	var rrm *message.RegisterResMessage = new(message.RegisterResMessage)
	err = json.Unmarshal([]byte(sm2.Data), rrm)
	if err != nil {
		fmt.Println("ClientRegister call json.Unmarshal err=", err)
		return
	}
	if rrm.ResCode == 400 {
		fmt.Println("ClientRegister say 注册成功,你可以登陆了")
		os.Exit(0)
	} else {
		fmt.Println(rrm.ErrMess)
		os.Exit(0) //退出系统
	}
	return

}
