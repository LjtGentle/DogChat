package model

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net"
	"strconv"

	"github.com/garyburd/redigo/redis"
)

var (
	MyUserDao *UserDao
)

type UserDao struct {
	Pool *redis.Pool
}

func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	userDao = &UserDao{
		Pool: pool,
	}
	return

}

//根据ID查找用户的信息    查操作
func (userDao *UserDao) GetUserById(id int) (user *User, err error) {
	//id为key
	conn := userDao.Pool.Get()
	defer conn.Close()
	res, err := redis.String(conn.Do("HGet", "users", id))
	if err != nil {
		if err == redis.ErrNil { //表示在users哈希中,没有对应的id
			err = ERROR_USER_NOTEXISTS
		}
		return
	}
	//res 是一个json格式的字串
	user = &User{}
	err = json.Unmarshal([]byte(res), user)
	if err != nil {
		fmt.Println("getUserById json.Unmarsh err=", err)
		return
	}

	return
}

//创建一个用户   增操作
func (userDao *UserDao) makeUser(id int, user *User) (err error) {
	//创建一个用户之前得先查找数据库中是否存在这样id的用户
	_, err = userDao.GetUserById(id)
	//根据id查找用户，没有错，就是出错了，改用户已经存在
	if err == nil {
		err = ERROR_USER_EXISTS
		return
	}
	conn := userDao.Pool.Get()
	defer conn.Close()
	date, err := json.Marshal(user)
	if err != nil {
		fmt.Println("makeUser json.Marshal err=", err)
		return
	}
	_, err = conn.Do("Hset", "users", id, string(date))
	if err != nil {
		fmt.Println("makeUser conn.Do err=", err)
		return
	}
	fmt.Println("成功在redis创建了一个用户")
	return
}

//登录操作
func (userDao *UserDao) LoginDAO(id int, pwd string) (err error) {
	user, err := userDao.GetUserById(id)
	if err != nil {
		fmt.Println("LoginADO call getUserById err=", err)
		return
	}
	if user.UserPwd != pwd {
		err = ERROR_USER_PWD
		return
	}
	return
}

//注册操作
func (userDao *UserDao) RegisteredDAO(id int, name, pwd string) (err error) {
	user := &User{
		UserId:   id,
		UserName: name,
		UserPwd:  pwd,
	}
	err = userDao.makeUser(id, user)
	if err != nil {
		fmt.Println("RegisteredDAO call makeUser err=", err)
		return
	}
	fmt.Println("注册成功.....")
	return
}

//////////////////////////--------下面都是不可行的方案

//用户登陆的时候 将用户加入在线用户列表   带用户名的
func (userDao *UserDao) UserOnline2(userId int, userName string, conn net.Conn) (err error) {
	redisconn := userDao.Pool.Get()
	defer redisconn.Close()
	onlineuser := OnlineUser{
		UserId:   userId,
		UserName: userName,
		Conn:     conn,
	}
	// fmt.Println("============>>>>>>", onlineuser)
	// data, err := asn1.Marshal(onlineuser)
	// if err != nil {
	// fmt.Println("UserOnline2 call asn1.Marshal err = ", err)
	// return
	// }
	//fmt.Println("<<<<<<<<<<<<<<<<<<<", data)
	date, err := json.Marshal(onlineuser)
	if err != nil {
		fmt.Println("UserOnline call json.Marshal(conn) err=", err)
		return
	}

	/////testing
	//反序列化回来看看
	// onlineuser2 := OnlineUser{}
	// err = json.Unmarshal(date, &onlineuser2)
	// if err != nil {
	// fmt.Println("test ----- json.Unmarshal err=", err)
	// return
	// }

	// fmt.Println("=========>", onlineuser2)
	// fmt.Println("=========>", onlineuser2.Conn)
	// fmt.Println("=========>", onlineuser2.Conn.(net.Conn))

	/////
	//////testing2
	// onlineuser3 := OnlineUser{}
	// rest, err := asn1.Unmarshal(data, &onlineuser3)
	// if err != nil {
	// fmt.Printf("<<<<<<<<<<<<<<<<<<<<rest:%#v,onlineuser3:%#v\n", rest, onlineuser3)
	// }
	/////

	////testing3  ----序列化
	//注册
	//gob.Register(map[int]string{})
	//gob.Register(&net.TCPConn{})
	gob.Register(OnlineUser{})
	//构造缓冲区
	buf := bytes.NewBuffer(nil)
	//生成god编码器
	g := gob.NewEncoder(buf)
	//编码
	err = g.Encode(onlineuser)
	if err != nil {
		fmt.Println("------UserOnline2 call g.Encode err=", err)
		return
	}
	fmt.Println(".>>>>>>>>>校验>>>>>>>>", buf.Bytes())
	////---反序列化
	var ou OnlineUser
	//构造阅读器
	r := bytes.NewReader(buf.Bytes())
	//构造gob解码器
	dg := gob.NewDecoder(r)
	//解码
	err = dg.Decode(&ou)
	if err != nil {
		fmt.Println("UserOnline2 call dg.Decode err=", err)
		return
	}
	//校验
	fmt.Printf(".>>>>>>>>>>>>>>ou:%#v,*ou.conn:%#v", ou, ou.Conn)

	////
	_, err = redisconn.Do("Hset", "Online", userId, string(date))
	if err != nil {
		fmt.Println("UserOnline call redisconn err=", err)
		return
	}
	fmt.Println("用户登陆后,成功把用户加入到redis在线表中")
	return
}

//用户登陆的时候 将用户加入在线用户列表   不带用户名的
func (userDao *UserDao) UserOnline(userId int, conn net.Conn) (err error) {
	err = MyUserDao.UserOnline2(userId, "***", conn)
	return
}

//用户退出的时候在用户列表中删除
func (userDao *UserDao) UserOutline(UserId int) (err error) {
	redisconn := userDao.Pool.Get()
	defer redisconn.Close()
	_, err = redisconn.Do("HDEL", "Online", UserId)
	if err != nil {
		fmt.Println("redisconn.Do err=", err)
		return
	}
	fmt.Printf("成功将[%d]用户,从redis在线表中删除\n", UserId)
	return
}

//获取全部在线用户
func (userDao *UserDao) GetAllOnline() (onlineusers map[int]*OnlineUser, err error) {

	redisconn := userDao.Pool.Get()
	defer redisconn.Close()

	onlineusers = make(map[int]*OnlineUser, 100)
	// data, err := redis.Values(redisconn.Do("HGETALL", "Online"))
	// fmt.Println("-------------", data, "---------------")

	// for _, value := range data {
	// onlineUser := &OnlineUser{}
	// fmt.Println("-------------", value.([]byte), "---------------")
	// err = json.Unmarshal(value.([]byte), onlineUser)
	// if err != nil {
	// fmt.Println("GetAllOnline call json.Unmarsh err=", err)
	// return
	// }
	// onlineusers[onlineUser.UserId] = onlineUser
	// }
	// return

	data, err := redis.StringMap(redisconn.Do("HGETALL", "Online"))
	fmt.Println("-------------", data, "---------------")
	for index, value := range data {
		fmt.Println("-------------", value, "---------------")
		onlineUser := &OnlineUser{}
		err = json.Unmarshal([]byte(value), onlineUser)
		if err != nil {
			println("GetAllOnline call json.Unmarsh err=", err)
			return
		}
		i, err := strconv.Atoi(index)
		if err != nil {
			println("GetAllOnline call strconv.Atoi err =", err)
		}
		onlineusers[i] = onlineUser
	}
	return

}

//根据conn删除key-value
func (userDao *UserDao) DelByConn(conn net.Conn) (err error) {
	//遍历所有conn,conn相等的value
	onlineusers, err := MyUserDao.GetAllOnline()
	if err != nil {
		fmt.Println("DelByConn call GetAllOnline err = ", err)
		return
	}
	for _, value := range onlineusers {
		if value.Conn == conn {
			//删除
			err = MyUserDao.UserOutline(value.UserId)
			fmt.Println("DelByConn call UserOutline err = ", err)
			return
		}
	}
	return

}
