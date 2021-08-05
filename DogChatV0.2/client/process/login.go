package process

import (
	"DogChat/client/model"
	"DogChat/common/message"
	"DogChat/common/utils"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
)

func ClientLogin(userId int, userPwd string) (err error) {
	//先链接服务器
	var address string
	address = fmt.Sprintf("%s:%d", model.ServerIp, model.ServerPort)
	fmt.Printf(">>>>>>address:%s\n", address)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("net.Dial=", err)
	}
	defer conn.Close()
	fmt.Println("成功链接上服务器...")
	//先凑个LoginMessage
	lm := message.LoginMessage{}

	lm.UserId = userId
	lm.UserPwd = userPwd

	//LoginMessage转json

	data, err := json.Marshal(lm)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	//凑一个SendMessage
	psm := &message.SendMessage{}
	psm.Type = message.LoginMessageType
	psm.Data = string(data)
	tf := utils.Transfer{}
	err = tf.PkgWrite(psm, conn)
	if err != nil {
		fmt.Println("ClientLogin call tf.PkgWrite err=", err)
	}
	fmt.Println("ClientLogin....读取...服务器返回的信息....")

	// 声明一个SendMessage读取服务器发送过来的信息
	psm2 := &message.SendMessage{}
	err = tf.PkgRead(psm2, conn)
	if err != nil {
		fmt.Println("ClientLogin call tf.PkgRead err=", err)
		return
	}
	if psm2.Type != message.LoginResMessageType {
		fmt.Println("服务器返回的消息不是message.LoginResMessageType,在这里不作处理")
		return
	}
	//反序列化psm2.Data
	//声明一个
	lrm := &message.LoginResMessage{}
	err = json.Unmarshal([]byte(psm2.Data), lrm)
	if err != nil {
		fmt.Println("ClientLogin call json.Unmarshal err=", err)
		err = errors.New(lrm.ErrMess)
		return
	}
	if lrm.ResCode == 500 {
		fmt.Println("用户登陆成功,----ClientLogin say")
		//....处理一些业务
		onlineUsers(conn) //向服务器请求当前在线用户
		//这里我们还需要在客户端启动一个协程
		//该协程保持和服务器端的通讯.如果服务器有数据推送给客户端
		//则接收并显示在客户端的终端.
		go serverProcessMes(conn)
		showLoginMenu(conn, userId)
	} else {
		fmt.Println("用用户登陆失败,代码是:", lrm.ResCode)

	}

	return
}

func showLoginMenu(conn net.Conn, id int) {

	var key int
	for {
		fmt.Println("------登陆成功--------")
		fmt.Println("------1.显示在线列表--------")
		fmt.Println("------2.公聊--------")
		fmt.Println("------3.私聊--------")
		fmt.Println("------4.退出系统--------")
		fmt.Println("请选择(1-4):")
		fmt.Scanf("%d\n", &key)
		switch key {
		case 1:
			fmt.Println("你选择了显示在线列表")
			onlineUsers2()
		case 2:
			fmt.Println("你选择了个公聊--------")
			chatMode(conn, id, Public)
		case 3:
			fmt.Println("你选择了私聊")
			chatMode(conn, id, Private)

		case 4:
			fmt.Println("你选择了退出系统")
			os.Exit(0)
		default:
			fmt.Println("请选择1-4:")
		}

	}

}

//展示客户端的在线列表
func onlineUsers2() {
	fmt.Println("----在线列表-----")
	for _, value := range model.Users {
		fmt.Printf("***ID:%d\tName:%s***\n", value.UserId, value.UserName)
	}
}

//展示在线列表----向服务器请求的在线列表
func onlineUsers(conn net.Conn) (err error) {
	//som := message.ShowOnlineMessage{}
	sm := message.SendMessage{
		Type: message.ShowOnlineMessageType,
	}
	tf := utils.Transfer{}
	err = tf.PkgWrite(&sm, conn)
	if err != nil {
		fmt.Println("onlineUsers call PkgWrite err=", err)
		return
	}
	return

}

const (
	Public = iota //
	Private
)

//公聊
func chatMode(conn net.Conn, sendId int, flag int) error {
	fmt.Println("come in ChatMode")
	defer fmt.Println("out of ChatMode")
	var id int
	if flag == Private {
		onlineUsers(conn)
		fmt.Println("请选择一个ID")
		fmt.Scanf("%d\n", &id)
	}
	if flag == Public {
		id = 000
	}

	fmt.Println(">>>>>>>>>请输入你要聊天的内容,输入回车发送,输入q退出<<<<<<<<<<")
	tf := utils.Transfer{}
	var mess string
	for {
		//fmt.Scanf("%s\n", &mess)
		//fmt.Scan(&mess)
		//fmt.Scanln(&mess)
		in := bufio.NewReader(os.Stdin)
		ret, err := in.ReadString('\n')
		retslice := strings.Split(ret, "\n")
		if err != nil {
			fmt.Println("chatMode call in.ReadString err=", err)
			return err
		}
		mess = retslice[0]
		// fmt.Printf(">>>>>>ret:%s\n<<<<<<\n", ret)
		// fmt.Printf("-----------retslice[0]:%s----------\n", retslice[0])

		// fmt.Printf(">>>%s<<<\n", mess)
		if mess == "q" {
			println("mess==q")
			break
		} else {
			psm := &message.SendMessage{}
			psm.Type = message.SmsMessageType
			sms := message.SmsMessage{
				SenderID:   sendId,
				ReceiverID: id,
				Message:    mess,
			}
			data, err := json.Marshal(sms)
			if err != nil {
				fmt.Println("pubilcChat call json.Marshal err=", err)
				return err
			}
			psm.Data = string(data)

			err = tf.PkgWrite(psm, conn)
			if err != nil {
				fmt.Println("pubilcChat call tf.PkgWrite err=", err)
				return err
			}

		}

	}
	return nil

}
