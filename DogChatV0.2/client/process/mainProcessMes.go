package process

import (
	"DogChat/client/model"
	"DogChat/common/message"
	"DogChat/common/utils"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
)

func serverProcessMes(conn net.Conn) {
	tf := utils.Transfer{}

	for {
		sm := message.SendMessage{}
		err := tf.PkgRead(&sm, conn)
		if err == io.EOF {
			fmt.Println("服务器关闭了,客户端也关闭吧")
			os.Exit(0)
		} else if err != nil {
			fmt.Println("serverProcessMes call tf.PkgRead err=", err)
			return
		}

		switch sm.Type {
		case message.SmsMessageType:
			ProcessSmsMessage(&sm)
		case message.ShowOnlineResMessageType:
			ProcessShowOnlineMessage(&sm)
		case message.NotifyOnlineType:
			ProcessNotifyOnlineUser(&sm)
		case message.NotifyOutlineType:
			ProcessNotifyOutlineUser(&sm)

		default:
			fmt.Println("来自于服务器的未知信息")

		}

	}
}

func ProcessSmsMessage(psm *message.SendMessage) (err error) {

	sms := message.SmsMessage{}
	err = json.Unmarshal([]byte(psm.Data), &sms)
	if err != nil {
		fmt.Println("ProcessSmsMessage call json.Unmarshal err=", err)
		return
	}
	str := fmt.Sprintf("%#v:%#v", sms.SenderID, sms.Message)
	fmt.Println(str)
	return

}

func ProcessShowOnlineMessage(psm *message.SendMessage) (err error) {
	//反序列化
	sorm := message.ShowOnlineResMessage{}
	err = json.Unmarshal([]byte(psm.Data), &sorm)
	if err != nil {
		println("onlineUsers call json.Unmarshal err=", err)
		return
	}
	if sorm.ErrMess != fmt.Sprintf("%#v\n", err) {
		fmt.Println("--===--ShowOnlineResMessage ErrMess=", sorm.ErrMess)
	}
	for _, value := range sorm.Onlineusers {
		//展示
		fmt.Printf("----ID:%d\tName:%s\n", value.UserId, value.UserName)
		//加入客户端的在线列表
		user := &model.User{
			UserId:   value.UserId,
			UserName: value.UserName,
		}
		model.Users[user.UserId] = user
	}
	return
}

func ProcessNotifyOnlineUser(psm *message.SendMessage) (err error) {
	//显示在线
	no := message.NotifyOnline{}
	err = json.Unmarshal([]byte(psm.Data), &no)
	if err != nil {
		fmt.Println("ProcessNotifyOnlineUser call json.Unmarshal err=", err)
		return
	}
	fmt.Printf("*******************用户上线啦!ID:%d\tName:%s\n", no.OnlineId, no.OnlineName)
	//加入客户端维护的在线列表中
	user := &model.User{
		UserId:   no.OnlineId,
		UserName: no.OnlineName,
	}
	model.Users[user.UserId] = user
	return
}

func ProcessNotifyOutlineUser(psm *message.SendMessage) (err error) {
	//显示下线
	nout := message.NotifyOutline{}

	err = json.Unmarshal([]byte(psm.Data), &nout)
	if err != nil {
		fmt.Printf("ProcessNotifyOutlineUser call json.Unmarshal err=", err)
		return
	}
	fmt.Printf("****************用户下线了!ID:%d\tName:%s\n", nout.OutlineId, nout.OutlineName)
	//在客户端维护的在线列表中删除
	delete(model.Users, nout.OutlineId)
	return

}
