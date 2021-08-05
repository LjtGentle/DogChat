package model

//客户端自己维护一个用户在线列表
var (
	Users map[int]*User
)

//引入改包就会自己调用该函数完成初始化
func init() {
	Users = make(map[int]*User, 1024)
}

type User struct {
	UserId   int    `json:"userId"`
	UserName string `json:"userName"`
}
