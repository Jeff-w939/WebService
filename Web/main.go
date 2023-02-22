package main

import (
	"Project3_WebService/Web/controller"
	"Project3_WebService/Web/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

//中间件过滤
func LoginFilter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 初始化session 对象
		s := sessions.Default(ctx)
		userName := s.Get("userName")
		if userName == nil {
			ctx.Abort() // 从这里返回，不必继续执行
		} else {
			ctx.Next() // 继续执行
		}
	}
}
func main() {
	// 初始化路由
	router := gin.Default()

	// 初始化redis
	model.InitRedis()

	//初始化mysql
	model.InitDb()

	// 初始化session容器
	store, _ := redis.NewStore(10, "tcp", "127.0.0.1:6379", "", []byte("bj38"))
	// 使用容器
	router.Use(sessions.Sessions("mysession", store))

	router.Static("/home", "view") // 自动 获取view 文件夹下的启动页面index.html

	r1 := router.Group("/api/v1.0")
	{
		r1.GET("/session", controller.GetSession)
		r1.GET("imagecode/:uuid", controller.GetImageCd)
		r1.GET("/smscode/:phone", controller.GetSmsCd)
		r1.POST("/users", controller.PostRet)
		r1.GET("/areas", controller.GetArea)
		r1.POST("/sessions", controller.PostLogin)

		r1.Use(LoginFilter()) // 以后的路由都不在校验session

		r1.DELETE("/session", controller.DeleteSession)
		r1.GET("/user", controller.GetuserInfo)
		r1.PUT("/user/name", controller.PutUserInfo)
		r1.POST("/user/avatar", controller.PostAvatar)
	}

	//3. 运行
	router.Run(":8080")
}
