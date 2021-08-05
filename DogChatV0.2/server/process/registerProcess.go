package process

import (
	"DogChat/common/message"
	"DogChat/common/utils"
	"DogChat/server/model"
	"encoding/json"
	"fmt"
	"net"
)

func ProcessRegister(conn net.Conn, psm *message.SendMessage) (err error) {
	var rm message.RegisterMessage
	err = json.Unmarshal([]byte(psm.Data), &rm)
	if err != nil {
		fmt.Println("ProessRegister json.Unmarshal err=", err)
		return
	}
	var rrm message.RegisterResMessage
	err = model.MyUserDao.RegisteredDAO(rm.UserId, rm.UserName, rm.UserPwd)
	if err != nil {
		fmt.Println("ProcessRegister call model.MyUserDAO err=", err)
		rrm.ResCode = 404
		rrm.ErrMess = "注册失败"

	} else {
		rrm.ResCode = 400
		rrm.ErrMess = "注册成功"
		fmt.Println("注册成功")
	} ////////////////////////////////////////////////////////////////////??////////////////
	var sm *message.SendMessage = new(message.SendMessage)
	sm.Type = message.RegisterResMessageType
	data, err := json.Marshal(&rrm)
	if err != nil {
		fmt.Println("ProcessRegister call json.Marshal err=", err)
		return
	}

	sm.Data = string(data)

	var tf utils.Transfer
	err = tf.PkgWrite(sm, conn)
	if err != nil {
		fmt.Println("ProcessRegister call tf.PkgWrite err=", err)
		return
	}

	return

}
