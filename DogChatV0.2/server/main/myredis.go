package main

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

//定义一个全局的pool
var pool *redis.Pool

func initPool(maxIdle, maxActive int, idleTimeout time.Duration) {
	pool = &redis.Pool{
		MaxIdle:     maxIdle,     //最大的空闲连接数
		MaxActive:   maxActive,   //最大的激活连接数，同时最多有N个连接
		IdleTimeout: idleTimeout, //空闲连接等待时间
		Dial: func() (redis.Conn, error) {
			//拨号redis
			conn, err := redis.Dial("tcp", "106.55.36.146:6566")
			if err != nil {
				fmt.Println("redis.Dial err =", err)
				return nil, err
			}
			//密码验证
			if _, err := conn.Do("AUTH", "ljt"); err != nil {
				conn.Close()
				return nil, err
			}
			return conn, err
		},
	}

}
