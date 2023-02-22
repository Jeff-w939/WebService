// 微服务： 获取图片验证码
package main

import (
	"fmt"
	"getCaptcha/handler"
	pb "getCaptcha/proto"
	"github.com/micro/go-micro"
	"github.com/micro/go-plugins/registry/consul"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	//srv := service.New(            代码自动生成的是这一行 是错的 得改成下面的
	//	service.Name("getcaptcha"),
	//)
	// 初始化consul对象
	consulReg := consul.NewRegistry()

	// new 一个微服务
	srv := micro.NewService(
		micro.Address("192.168.102.1:12345"),
		micro.Name("go-micro-srv-getCaptcha"),
		micro.Registry(consulReg),
	)
	// Register handler
	err := pb.RegisterGetCaptchaHandler(srv.Server(), handler.New())
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
