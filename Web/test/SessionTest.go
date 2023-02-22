package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"

	//"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// 初始化容器
	store, _ := redis.NewStore(10, "tcp", "127.0.0.1:6379", "", []byte("wangbaba"))
	//使用容器
	router.Use(sessions.Sessions("mysession", store))

	router.GET("/test", func(ctx *gin.Context) {
		//调用session, 设置session 数据
		s := sessions.Default(ctx)

		//设置session
		s.Set("wangbaba", "guo")

		s.Save()
		ctx.Writer.WriteString("测试session...")
	})
	router.Run(":9999")
}
