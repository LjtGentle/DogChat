package process

import (
	"DogChat/common/message"
	"DogChat/common/utils"
	"fmt"
	"io"
	"net"
)

func MainProcess(conn net.Conn) (err error) {

	defer fmt.Println("end of MainProcess()")

	//一直读取客户端的信息
	for {
		//实例化对象
		tf := utils.Transfer{}
		psm := &message.SendMessage{}
		err = tf.PkgRead(psm, conn)
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端退出了程序")
				fmt.Printf("删除用户在线列表,conn=%#v\n", conn)
				ProcessNotifyOutline(conn)
				fmt.Println("删除在线成功")
				return
			}
			fmt.Println("MainProcess tf err=", err)
			return
		}
		//....分析消息的类型，给小弟干活
		go func() {
			err = HandleMessage(conn, psm)
			if err != nil {
				fmt.Println("MainProcess call HandleMessage err=", err)
				return
			}
		}()

	}

}

func HandleMessage(conn net.Conn, psm *message.SendMessage) (err error) {

	switch psm.Type {
	case message.LoginMessageType:
		fmt.Println("服务器接收到客户端的一个登陆请求")
		err = ProcessLoginMessage(conn, psm)
		fmt.Println("服务器接收到客户端的一个登陆请求,处理完毕")
	case message.RegisterMessageType:
		fmt.Println("服务器接收到客户端的一个注册请求")
		err = ProcessRegister(conn, psm)
		fmt.Println("服务器接收到客户端的一个注册请求,处理完毕")
	case message.ShowOnlineMessageType:
		fmt.Println("服务器接收到客户端的展示用户在线列表的请求")
		err = ProcessShowOnlineMessage(conn, psm)
		fmt.Println("服务器接收到客户端的展示用户在线列表的请求,已处理完毕")

	case message.SmsMessageType:
		fmt.Println("服务器收到客户端的一个聊天请求")
		err = ProcessSmsMessage(conn, psm)
		fmt.Println("服务器收到客户端的一个聊天请求,已处理完毕")

	default:
		fmt.Println("未知的请求,服务器无法进行处理")
	}
	return
}
