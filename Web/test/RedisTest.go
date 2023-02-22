package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

func main() {
	//1 连接redis
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("Dial err :", err)
		return
	}
	defer conn.Close()

	//2. 操作数据库
	reply, err := conn.Do("set", "itwang", "itbaba")

	//3.回复助手类函数 ---- 确定成具体的数据类型
	rep, e := redis.String(reply, err)
	fmt.Print(rep, e)
}
