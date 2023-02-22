package main

import (
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-plugins/registry/consul"
	"github.com/micro/micro/v3/service/logger"
	"users/handler"
	"users/model"
	pb "users/proto"
)

func main() {
	// 初始化Mysql 连接池
	_, err := model.InitDb()
	if err != nil {
		fmt.Print("初始化Mysql 失败 err:", err)
	}
	// 初始化redis连接池
	model.InitRedis()

	//1. 创建consul对象
	consulreg := consul.NewRegistry()

	//2. 注册微服务
	service := micro.NewService(
		micro.Address("192.168.102.1:12342"),
		micro.Name("go.micro.srv.user"),
		micro.Registry(consulreg),
	)

	// Register handler
	err = pb.RegisterUsersHandler(service.Server(), handler.New())
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	// Run service
	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}
