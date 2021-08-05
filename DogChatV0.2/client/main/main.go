package main

import (
	"DogChat/client/model"
	"DogChat/client/process"
	"flag"
	"fmt"
	"os"
)

func init() {
	flag.StringVar(&model.ServerIp, "ip", "0.0.0.0", "set server ip(default is 159.75.91.76)")
	flag.IntVar(&model.ServerPort, "port", 7576, "setup up server port(default is 7576)")
	//fmt.Println("main.go init() runned")
	//解析上面定义的标签 没有这个命令行输入就每意义了
	flag.Parse()
	fmt.Printf("ServerIp:%s\tServerPort:%d\n", model.ServerIp, model.ServerPort)
}

func main() {

	var key int

	for {
		fmt.Println("-----------------------欢迎使用DogChat聊天系统----------------------")
		fmt.Println("------------------------1.登录------------------------------------")
		fmt.Println("------------------------2.注册------------------------------------")
		fmt.Println("------------------------3.退出------------------------------------")
		fmt.Println("------------------------4.请输入（1-3）----------------------------")
		fmt.Scanf("%d\n", &key)
		switch key {
		case 1: //登陆
			fmt.Println("你选择了登陆")
			fmt.Println("请输入ID：")
			var userId int
			fmt.Scanf("%d\n", &userId)
			fmt.Println("请输入密码：")
			var userPwd string
			fmt.Scanf("%s\n", &userPwd)
			err := process.ClientLogin(userId, userPwd)
			if err != nil {
				println("main call process.ClientLogin err=", err)
			}
			os.Exit(0)
		case 2: //注册
			fmt.Println("你选择了注册")
			fmt.Println("请输入ID:")
			var userId int
			fmt.Scanf("%d\n", &userId)
			fmt.Println("请输入用户名:")
			var userName string
			fmt.Scanf("%s\n", &userName)
			fmt.Println("请输入密码:")
			var userPwd string
			fmt.Scanf("%s\n", &userPwd)
			err := process.ClientRegister(userId, userName, userPwd)
			if err != nil {
				println("main call process.ClientRegister err=", err)
			}
			os.Exit(0)
		case 3: //退出
			os.Exit(0) //退出程序
		default:
			fmt.Println("输入有误，请输入1到3")

		}

	}

}
