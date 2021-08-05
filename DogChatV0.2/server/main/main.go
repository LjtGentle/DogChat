package main

import (
	"DogChat/server/model"
	"DogChat/server/process"
	"fmt"
	"net"
	"time"
)

func init() {
	initPool(16, 32, 5*time.Minute)
	model.MyUserDao = model.NewUserDao(pool)
}

func main() {
	//监听一个端口
	listener, err := net.Listen("tcp", "0.0.0.0:7576")
	if err != nil {
		fmt.Println("net.Listen err=", err)
	}
	defer listener.Close()
	fmt.Println("tf在7576开启服务器成功，等待客户端来连接.....")
	//一直等待链接，一来链接就丢给MainProcess处理
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("server main listener.Accept() err=", err)
				continue
			}
			go process.MainProcess(conn)
		}
	}()
	//阻塞等待
	select {}
}
