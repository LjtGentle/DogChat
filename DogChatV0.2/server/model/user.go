package model

import (
	"errors"
	"fmt"
	"net"
)

//用于序列化后存储到redis的value
type User struct {
	UserId   int    `json:"userId"`
	UserName string `json:"userName"`
	UserPwd  string `json:"userPwd"`
}

type OnlineUser struct {
	UserId   int      `json:"userId"`
	UserName string   `json:"userName"`
	Conn     net.Conn `json:"conn"`
}

var (
	MyUserMgr *UserMgr
)

//服务器维护一个用户在线列表
type UserMgr struct {
	onlineUsers map[int]*OnlineUser
}

//初始化工作,只要引入包就会自动调用
func init() {
	MyUserMgr = &UserMgr{
		onlineUsers: make(map[int]*OnlineUser, 1024),
	}

}

//增加一个在线用户

func (this *UserMgr) AddOnlineUser(ou *OnlineUser) (err error) {
	//一个用户登陆两次或以上
	if _, ok := this.onlineUsers[ou.UserId]; ok {
		fmt.Printf("用户Id:%#v已经在服务器的用户在线列表中了\n", ou.UserId)
		err = errors.New("用户已经在用户列表中存在")
		return

	} else { //加入服务器维护的用户在线别表
		this.onlineUsers[ou.UserId] = ou
	}
	// //应该把这个在线的人广播出去,除了他自己本身
	// no := message.NotifyOnline{
	// OnlineName: ou.UserName,
	// OnlineId:   ou.UserId,
	// }
	// data, err := json.Marshal(no)
	// if err != nil {
	// fmt.Println("AddOnlineUser call json.Marshal err=", err)
	// return
	// }
	// sm := message.SendMessage{
	// Type: message.NotifyOnlineType,
	// Data: string(data),
	// }
	// //广播

	return
}

//删除一个在线用户
func (this *UserMgr) DelOnlineUserById(userId int) {
	//服务器维护的在线列表中删除
	delete(this.onlineUsers, userId)
	//
}
func (this *UserMgr) DelOnlineUserByConn(conn net.Conn) (UserId int,UserName string,err error) {
	flag := false
	for index, value := range this.onlineUsers {
		if value.Conn == conn {
			this.DelOnlineUserById(index)
			flag = true
			UserId = value.UserId
			UserName = value.UserName
			break
		}

	}
	if !flag {
		fmt.Printf("该conn:%#v不在在线用户列表中\n", conn)
		err = errors.New("该conn不在在线用户列表中")
	}
	return
}

//返回所有在线用户
func (this *UserMgr) GetAllOnlineUsers() map[int]*OnlineUser {
	return this.onlineUsers

}

//根据用户id 返回在线用户结构体
func (this *UserMgr) GetOnlineUserById(userId int) (ou *OnlineUser, err error) {
	if value, ok := this.onlineUsers[userId]; ok {
		ou = value
	} else {
		err = errors.New("该用户不在用户列表")
		fmt.Printf("用户id:%#v不在用户在线列表\n", userId)
	}
	return
}
