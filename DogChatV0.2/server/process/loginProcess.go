package process

import (
	"DogChat/common/message"
	"DogChat/common/utils"
	"DogChat/server/model"
	"encoding/json"
	"fmt"
	"net"
)

func ProcessLoginMessage(conn net.Conn, psm *message.SendMessage) (err error) {
	//反序列化sm.Data
	defer fmt.Println("end of ProcessLoginMessage()")
	var lm message.LoginMessage
	err = json.Unmarshal([]byte(psm.Data), &lm)
	if err != nil {
		fmt.Println("Sub json.Unmarshal err=", err)
		return
	}
	//创建一个LoginResMessage返回给客户端
	var lrm message.LoginResMessage
	err = model.MyUserDao.LoginDAO(lm.UserId, lm.UserPwd)
	if err != nil {
		fmt.Println("ProcessLoginMessage call model.MyUserDao.LoginDAO err=", err)
		lrm.ResCode = 505
		lrm.ErrMess = "登陆失败"
		fmt.Println("登陆失败")

	} else {
		fmt.Println("登陆成功")
		User := &model.User{}
		User, err = model.MyUserDao.GetUserById(lm.UserId)
		if err != nil {
			fmt.Println("ProcessLoginMessage call model.MyUserDao.GetUserById err=", err)
			return
		}
		fmt.Printf("加入在线用户列表,userId:%#v,conn:%#v,userName:%#v\n", lm.UserId, conn, User.UserName)
		ou := &model.OnlineUser{
			UserId:   lm.UserId,
			Conn:     conn,
			UserName: User.UserName,
		}
		err = ProcessNotifyOnline(ou)
		if err != nil {
			fmt.Println("ProcessLoginMessage call UserOnline err=", err)
			return
		}
		fmt.Println("加入在线列表成功")
		lrm.ResCode = 500
		lrm.ErrMess = "登陆成功"
	}
	//登陆成功了,给客户端返回一个信息
	var tf utils.Transfer
	//将LoginResMessage序列化
	data, err := json.Marshal(&lrm) //传不传指针都行
	if err != nil {
		fmt.Println("LoginResMessage json.Marshal err=", err)
		return
	}
	fmt.Printf("test lrm,lrm=%v\n", lrm)
	// var sm *message.SendMessage
	// fmt.Println("test1------")
	// sm.Type = message.LoginResMessageType
	// fmt.Println("test2------")
	// sm.Data = string(data)
	// fmt.Println("test3------")
	fmt.Println("test1.............")
	sm := &message.SendMessage{
		Type: message.LoginResMessageType,
		Data: string(data),
	}
	fmt.Println("test2............")
	err = tf.PkgWrite(sm, conn)
	fmt.Println("test4------")
	if err != nil {
		fmt.Println("ProcessLoginMessage call PkgWrite err=", err)
	}
	return
}

func ProcessShowOnlineMessage(conn net.Conn, psm *message.SendMessage) (err error) {
	sorm := message.ShowOnlineResMessage{}

	//这里psm.Data 是空的
	// err = json.Unmarshal([]byte(psm.Data), &sorm)
	// if err != nil {
	// println("ProcessShowOnlineMessage call json.Unmarshal err=", err)
	// return
	// }
	onlineuser := model.MyUserMgr.GetAllOnlineUsers()
	online := make(map[int]*message.Online)
	for _, vaule := range onlineuser {
		online[vaule.UserId] = &message.Online{
			UserId:   vaule.UserId,
			UserName: vaule.UserName,
		}
	}
	sorm.Onlineusers = online
	sorm.ErrMess = fmt.Sprintf("%#v\n", err)

	data, err := json.Marshal(sorm)
	if err != nil {
		println("ProcessShowOnlineMessage call json.Marshal err=", err)
	}
	sm := &message.SendMessage{
		Type: message.ShowOnlineResMessageType,
		Data: string(data),
	}
	ts := utils.Transfer{}
	err = ts.PkgWrite(sm, conn)
	if err != nil {
		println("ProcessShowOnlineMessage call ts.PkgWrite err=", err)
	}
	return

}

func ProcessSmsMessage(conn net.Conn, psm *message.SendMessage) (err error) {
	//声明一个结构体
	sms := &message.SmsMessage{}
	err = json.Unmarshal([]byte(psm.Data), sms)
	if err != nil {
		fmt.Println("ProcessSmsMessage call json.Unmarshal err=", err)
		return
	}
	if sms.ReceiverID == 000 {
		//公聊
		//除了自己本身的所有在线用户都给发
		onlines := model.MyUserMgr.GetAllOnlineUsers()

		for _, vaule := range onlines {

			if vaule.UserId == sms.SenderID {
				continue
			}
			tf := utils.Transfer{}
			err = tf.PkgWrite(psm, vaule.Conn)
			if err != nil {
				fmt.Println("ProcessSmsMessage call tf.PkgWrite err=", err)
				return
			}
		}

	} else {
		//私聊
		online := &model.OnlineUser{}
		online, err = model.MyUserMgr.GetOnlineUserById(sms.ReceiverID)
		if err != nil {
			fmt.Println("ProcessSmsMessage call model.MyUserMgr.GetOnlineUserById err=", err)
			return
		}
		tf := utils.Transfer{}
		err = tf.PkgWrite(psm, online.Conn)
		if err != nil {
			fmt.Println("ProcessSmsMessage call tf.PkgWrite err=", err)
			return

		}

	}
	return

}

func ProcessNotifyOnline(ou *model.OnlineUser) (err error) {
	//服务器维护的在线列表
	err = model.MyUserMgr.AddOnlineUser(ou)
	if err != nil {
		fmt.Println("ProcessNotifyOnline call model.MyUserMgr.AddOnlineUser err=", err)
		return
	}
	//应该把这个在线的人广播出去
	no := message.NotifyOnline{
		OnlineName: ou.UserName,
		OnlineId:   ou.UserId,
	}
	data, err := json.Marshal(no)
	if err != nil {
		fmt.Println("AddOnlineUser call json.Marshal err=", err)
		return
	}
	sm := message.SendMessage{
		Type: message.NotifyOnlineType,
		Data: string(data),
	}
	//广播除了自己本身
	onlines := model.MyUserMgr.GetAllOnlineUsers()
	for _, vaule := range onlines {
		if vaule.UserId == ou.UserId {
			continue
		}
		tf := utils.Transfer{}
		err = tf.PkgWrite(&sm, vaule.Conn)
		if err != nil {
			fmt.Println("ProcessSmsMessage call tf.PkgWrite err=", err)
			return
		}
	}
	return

}

func ProcessNotifyOutline(conn net.Conn) (err error) {
	//服务器本身维护的在线列表中删除
	UserId, UserName, err := model.MyUserMgr.DelOnlineUserByConn(conn)
	nout := message.NotifyOutline{
		OutlineName: UserName,
		OutlineId:   UserId,
	}
	data, err := json.Marshal(nout)
	if err != nil {
		fmt.Println("ProcessNotifyOutline call json.Marshal err=", err)
		return
	}
	sm := message.SendMessage{
		Type: message.NotifyOutlineType,
		Data: string(data),
	}
	onlines := model.MyUserMgr.GetAllOnlineUsers()
	for _, vaule := range onlines {
		// if vaule.UserId == ou.UserId {
		// continue
		// }
		tf := utils.Transfer{}
		err = tf.PkgWrite(&sm, vaule.Conn)
		if err != nil {
			fmt.Println("ProcessSmsMessage call tf.PkgWrite err=", err)
			return
		}
	}
	return

}
