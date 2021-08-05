package message

import (
	_ "encoding/json"
)

const (
	LoginMessageType         = "LoginMessage"
	LoginResMessageType      = "LoginResMessage"
	RegisterMessageType      = "RegisterMessage"
	RegisterResMessageType   = "RegisterResMessage"
	ShowOnlineMessageType    = "ShowOnlineMessage"
	ShowOnlineResMessageType = "ShowOnlineResMessage"
	SmsMessageType           = "SmsMessage"
	NotifyOnlineType         = "NotifyOnline"
	NotifyOutlineType        = "NotifyOutline"
)

type SendMessage struct {
	//PkgLen uint32 `json:"pkgLen"`
	Type string `json:"type"`
	Data string `json:"data"`
}

type LoginMessage struct {
	//UserName string
	UserId  int    `json:"userId"`
	UserPwd string `json:"userPwd"`
}

type LoginResMessage struct {
	//错误代码
	ResCode uint32 `json:"errCode"`
	ErrMess string `json:"errMess"`
}

type RegisterMessage struct {
	UserId   int    `json:"userId"`
	UserName string `json:"userName"`
	UserPwd  string `json:"userPwd"`
}

type RegisterResMessage struct {
	ResCode uint32 `json:"errCode"`
	ErrMess string `json:"errMess"`
}

type ShowOnlineMessage struct {
}

type Online struct {
	UserId   int    `json:"userId"`
	UserName string `json:"userName"`
}

type ShowOnlineResMessage struct {
	ErrMess     string          `json:"errMess"`
	Onlineusers map[int]*Online `json:"onlineusers"`
}

type User struct {
	UserId   int    `json:"userId"`
	UserName string `json:"userName"`
	UserPwd  string `json:"userPwd"`
}

type SmsMessage struct {
	SenderID   int    `json:"sender"`
	ReceiverID int    `json:"receiver"` //如果id为000即广播
	Message    string `json:"message"`
}

type NotifyOnline struct {
	OnlineName string `json:"onlineName"`
	OnlineId   int    `json:"onlineId"`
}

type NotifyOutline struct {
	OutlineName string `json:"outlineName"`
	OutlineId   int    `json:"outlineId"`
}
