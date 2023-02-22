package controller

import (
	"Project3_WebService/Web/model"
	getCaptcha "Project3_WebService/Web/proto/getCaptcha"
	userMicro "Project3_WebService/Web/proto/users"
	"Project3_WebService/Web/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/afocus/captcha"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/micro/go-micro"
	"github.com/micro/go-plugins/registry/consul"
	"image/png"
	"net/http"
)

// 1. 获取Session信息
func GetSession(ctx *gin.Context) {
	// 初始化错误返回的map
	//resp := make(map[string]string)
	//resp["errno"] = utils.RECODE_SESSIONERR
	//resp["errmsg"] = utils.RecodeText(utils.RECODE_SESSIONERR)
	//ctx.JSON(http.StatusOK, resp)
	resp := make(map[string]interface{}) // 通过

	// 获取 Session 数据
	s := sessions.Default(ctx) // 初始化 Session 对象
	userName := s.Get("userName")

	// 用户没有登录.---没存在 MySQL中, 也没存在 Session 中
	if userName == nil {
		resp["errno"] = utils.RECODE_SESSIONERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_SESSIONERR)
	} else {
		resp["errno"] = utils.RECODE_OK
		resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)

		var nameData struct {
			Name string `json:"name"`
		}
		nameData.Name = userName.(string) // 类型断言
		resp["data"] = nameData
	}

	ctx.JSON(http.StatusOK, resp)
}

//2. 获取图片-微服务实现
func GetImageCd(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	// 指定consul 服务发现
	consulReg := consul.NewRegistry()
	consulService := micro.NewService(
		micro.Registry(consulReg),
	)

	//初始化客户端
	microClient := getCaptcha.NewGetCaptchaService("go-micro-srv-getCaptcha", consulService.Client())

	resp, err := microClient.Call(context.TODO(), &getCaptcha.Request{Uuid: uuid})
	if err != nil {
		fmt.Println("未找到远程服务。。。")
		return
	}
	var img captcha.Image
	json.Unmarshal(resp.Msg, &img)
	png.Encode(ctx.Writer, img)
	fmt.Println("uuid :", uuid)

}

// 3. 发送短信验证码 - 微服务实现
func GetSmsCd(ctx *gin.Context) {
	// 获取短信验证码
	phone := ctx.Param("phone")

	imgCode := ctx.Query("text") // 这个是自己手写的图片验证码

	uuid := ctx.Query("id")

	//指定consul 服务发现 三步
	consulReg := consul.NewRegistry()
	consulService := micro.NewService(
		micro.Registry(consulReg),
	)
	// 调用远程服务 SendSms()
	microClient := userMicro.NewUsersService("go.micro.srv.user", consulService.Client())
	resp, err := microClient.SendSms(context.TODO(), &userMicro.Request{Phone: phone, ImageCode: imgCode, Uuid: uuid})
	if err != nil {
		fmt.Println("调用远程SendSms 失败")
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

//4. 发送注册信息 - 微服务实现
func PostRet(ctx *gin.Context) {

	//request-payload 数据格式需要这样获取
	var regData struct {
		Mobile   string `json:"mobile"`
		PassWord string `json:"password"`
		SmsCode  string `json:"sms_code"`
	}
	ctx.Bind(&regData)
	//指定consul 服务发现 三步
	consulReg := consul.NewRegistry()
	consulService := micro.NewService(
		micro.Registry(consulReg),
	)
	microClient := userMicro.NewUsersService("go.micro.srv.user", consulService.Client())
	//调用远程函数  Register()
	resp, err := microClient.Register(context.TODO(), &userMicro.RegReq{
		Mobile:   regData.Mobile,
		SmsCode:  regData.SmsCode,
		Password: regData.PassWord,
	})

	if err != nil {
		fmt.Println("注册用户找不到远程服务：", err)
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

//5. 获取地域信息
func GetArea(ctx *gin.Context) {
	////提高用户感受的常见方法：先查缓存redis， 缓存没有查MySQL， 写入redis缓存。
	var areas []model.Area

	// 从缓存redis 中, 获取数据, 先得到一个redispool中的连接conn
	conn := model.RedisPool.Get()
	// 当初使用 "字节切片" 存入, 现在使用 切片类型接收
	areaData, _ := redis.Bytes(conn.Do("get", "areaData"))
	// 没有从 Redis 中获取到数据
	if len(areaData) == 0 {

		fmt.Println("从 MySQL 中 获取数据...")
		model.GlobalConn.Find(&areas)
		// 把数据写入到 redis 中. , 存储结构体序列化后的 json 串
		areaBuf, _ := json.Marshal(areas)
		_, err := conn.Do("set", "areaData", areaBuf)
		if err != nil {
			fmt.Println(" 写入redis 失败：", err)
		}
	} else {
		fmt.Println("从 redis 中 获取数据...")
		// redis 中有数据
		json.Unmarshal(areaData, &areas)
	}

	resp := make(map[string]interface{})

	resp["errno"] = "0"
	resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
	resp["data"] = areas

	ctx.JSON(http.StatusOK, resp)
}

//6.处理登录业务
func PostLogin(ctx *gin.Context) {
	fmt.Println("PostLogin....")
	var loginData struct {
		Mobile   string `json:"mobile"`
		PassWord string `json:"password"`
	}
	ctx.Bind(&loginData)
	resp := make(map[string]interface{})

	//获取 数据库数据, 查询是否和数据的数据匹配
	userName, err := model.Login(loginData.Mobile, loginData.PassWord)
	fmt.Println("开始登录。。。")
	if err == nil {
		// 登录成功!
		fmt.Println("登录成功。。。")
		resp["errno"] = utils.RECODE_OK
		resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)

		// 将 登录状态, 保存到Session中
		s := sessions.Default(ctx)  // 初始化session
		s.Set("userName", userName) // 将用户名设置到session中.
		s.Save()

	} else {
		// 登录失败!
		fmt.Println("登录失败。。。")
		resp["errno"] = utils.RECODE_LOGINERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_LOGINERR)
	}

	ctx.JSON(http.StatusOK, resp)
}

//7.退出登录
func DeleteSession(ctx *gin.Context) {
	resp := make(map[string]interface{})
	s := sessions.Default(ctx)
	s.Delete("userName")
	err := s.Save()
	if err != nil {
		resp["errno"] = utils.RECODE_IOERR // 没有合适错误,使用 IO 错误!
		resp["errmsg"] = utils.RecodeText(utils.RECODE_IOERR)

	} else {
		resp["errno"] = utils.RECODE_OK
		resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
	}
	ctx.JSON(http.StatusOK, resp)
}

//8.获取用户信息
func GetuserInfo(ctx *gin.Context) {
	// 获取session, 得到当前用户信息
	s := sessions.Default(ctx)
	userName := s.Get("userName")
	resp := make(map[string]interface{})
	defer ctx.JSON(http.StatusOK, resp)

	if userName == nil {
		resp["errno"] = utils.RECODE_SESSIONERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_SESSIONERR)
		return
	}
	user, err := model.GetInfo(userName.(string))
	if err != nil {
		resp["errno"] = utils.RECODE_DBERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_DBERR)
		return
	}
	resp["errno"] = utils.RECODE_OK
	resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)

	temp := make(map[string]interface{})
	temp["user_id"] = user.ID
	temp["name"] = user.Name
	temp["mobile"] = user.Mobile
	temp["real_name"] = user.Real_name
	temp["Id_card"] = user.Id_card
	temp["Avatar_url"] = user.Avatar_url

	resp["data"] = temp
}

//9. 修改用户名
func PutUserInfo(ctx *gin.Context) {
	//先从session 中获取旧用户名
	s := sessions.Default(ctx)
	userName := s.Get("userName")

	// 获取新用户名
	var nameData struct {
		Name string `json:"name"`
	}
	ctx.Bind(&nameData)

	// 更新用户名
	resp := make(map[string]interface{})
	defer ctx.JSON(http.StatusOK, resp)

	// 更新数据库mysql中的 name
	err := model.UpdateUserName(nameData.Name, userName.(string))
	if err != nil {
		resp["errno"] = utils.RECODE_DBERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_DBERR)
		return
	}

	// 更新session中的数据
	s.Set("userName", nameData.Name)
	err = s.Save() // 必须保存
	if err != nil {
		resp["errno"] = utils.RECODE_SESSIONERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_SESSIONERR)
		return
	}
	resp["errno"] = utils.RECODE_OK
	resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
	resp["data"] = nameData
}

//10. 上传用户头像
func PostAvatar(ctx *gin.Context) {
	// 获取图片文件， 静态文件对象
	file, _ := ctx.FormFile("avatar")

	// 上传文件到项目中
	err := ctx.SaveUploadedFile(file, "test/"+file.Filename)
	fmt.Println(err)
}
