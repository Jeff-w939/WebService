# WebService
# GO语言学习

### 二：项目—Ihome 租房网

启动命令：1 、先启动consul 服务发现：consul agent -dev（需要到consul的文件夹进入cmd)

​					2、 再启动redis 和mysql:  redis-cli(需要到redis的文件夹进入cmd)  mysql 直接打开即可 输入密码199791

​				   3   再启动server端的各个微服务的main.go： go run main.go

​					4  最后启动 web 端 的 main.go： gu run main.go

​					5 最后浏览器输入url : 自己ip ： 196.128.102.1:8080/home

#### 2.1：单体式服务和微服务

#### 2.2 ：rpc协议和http协议

#### 2.3：grpc

#### 2.4：protobuf

#### 2.5 :  go-micro

micro new --type 服务名     ——创建一个微服务项目

<img src="图片/image-20230218162925223.png" alt="image-20230218162925223" style="zoom: 67%;" />

#### 2.6 ：服务发现-consul

<img src="图片/image-20230218103528806.png" alt="image-20230218103528806" style="zoom:80%;" />

命令：consul agent -dev : 启动一个本机的服务发现

在网站输入 127.0.0.1:8500 web查看

![image-20230218103048723](图片/image-20230218103048723.png)

命令：consul members : 查看有多少个成员

<img src="图片/image-20230218103112511.png" alt="image-20230218103112511" style="zoom: 67%;" />

consul info : 查看当前consul 的 ip 信息

consul leave : 优雅关闭consul

##### 2.6.1 注册服务

##### 2.6.2 查询服务

##### 2.6.3 健康检查

<img src="图片/image-20230218105411359.png" alt="image-20230218105411359" style="zoom: 67%;" />

可以通过consul  实现一个简单的负载均衡， 轮询使用服务。

#### 2.7 gin 框架

1. 封装回调函数， 给 router.Get() 设置

2. 拷贝 微服务的 “ 密码本” protobuf 到  web 中

3. 修改 protobuf文件的 包名。  test66别名  “test66web/proto/test66”

4. 实现 回调函数：

   1. 初始化客户端。    microClient := NewTeset66Sevice(服务名，client.DefaultClient)

   2. 调用远程服务。    resp, err := microClient.Call(context.TODO, &test66.Request{

      ​						       })

   3. 将 返回的 数据， 显示到 浏览器。 context.Writer.WriteString(resp.Msg);

```go
func main () {
    // 1 初始化路由
    router := gin.Default()
    
    // 2 做路由匹配
    router.GET = ("/", func(context *gin.Context) {
        context.Writer.writeString("hello world!")
    } )
    
    //3 启动运行
    router.Run("8080")
}
```

![image-20230218204025565](图片/image-20230218204025565.png)

web服务对于浏览器来说就是服务端，因为web服务需要通过远程调用远程的微服务来进行反馈给浏览器。

优化：1、 可以在浏览器和web服务中间加入nginx实现反向代理。

​			2、 微服务可以采用同名服务，减少故障发生概率。

​            3、 服务发现处实现负载均衡，通过均匀调用同名微服务，减少压力。

#### 2.8 项目准备

1. 准备项目环境。
   1. 创建项目目录  web、service
   2. 在 web 端 使用 MVC
   3. 创建项目常用目录： conf 配置文件、utils 工具类、比如error错误项。bin可执行文件、test测试目录， view包含 前端页面文件（一般前端写好的文件就放到view文件中）
   4. 导入 异常处理error.go
   5. 导入前端资源 html/ 到 view/ 中
2. 开发项目
   1. 开发 微服务端
   2. 开发 web 服务（客户端）



#### 3.1 session 获取

1.  在 web/main.go 中 ， 跟据 gin 框架 使用static() , 设置访问路径
2.  F12 浏览器中，查看 NetWork 中 Headers 和  Respose。 得到 url
3.  查看 《接口文档.doc》, 获取 url 、错误码、错误处理函数。
4.  在 web/ 下 遵循 MVC 设计模式创建  controller 目录。添加 user.go
5.  根据 《接口文档.doc》实现错误函数。
    1. resp[“errno”]
    2. resp[“errmsg”]
    3. ctx.Json(200, resp)  // 将 错误消息，进行序列化。返回给浏览器。
6.  url寻址时，都是从 “/‘’ 开始， 产生歧义
    - router.Static(“/”)    --- 修改为： router.Static(“/home”)

7.  浏览器测试： IP:8080/home 

#### 3.2 图片验证码ID获取

1. 启动 web页面，点击“注册” 按钮。在 NetWork —— Headers 中 看到 错误信息！

2. 从 URL中，提取 图片验证码ID 。 保存 成 uuid

3. 查看 gin 框架 中文文档。—— “获取路径中的参数”

4. web/main.go 中 添加 路由匹配：

   ```go
   router.GET("/api/v1.0/imagecode/:uuid", controller.GetImageCd)
   ```

5. web/controller/user.go 中 ，实现回调：

   ```go
   // 获取图片信息
   func GetImageCd(ctx *gin.Context) {
   	uuid := ctx.Param("uuid")
   	fmt.Println("uuid = ", uuid)
   }
   ```

#### 3.3 图片验证码获取

- 去 github 中搜索 “captcha” 关键词。 过滤 Go语言。 ——  [afocus/*captcha*](https://github.com/afocus/captcha)


- 使用 go get github.com/afocus/captcha 下载源码包。

- 参照 github 中的示例代码，测试生成 图片验证码：

  ```go
  package main
  
  import (
  	"github.com/afocus/captcha"   // 按住 Ctrl ，鼠标左键点击 captcha 看到 examples， 从中可以提取到 “comic.ttf”
  	"image/color"
  	"image/png"
  	"net/http"
  )
  
  func main()  {
  	// 初始化对象
  	cap := captcha.New()
  
  	// 设置字体
  	cap.SetFont("comic.ttf")
  
  	// 设置验证码大小
  	cap.SetSize(128, 64)
  
  	// 设置干扰强度
  	cap.SetDisturbance(captcha.MEDIUM)
  
  	// 设置前景色
  	cap.SetFrontColor(color.RGBA{0,0,0, 255})
  
  	// 设置背景色
  	cap.SetBkgColor(color.RGBA{100,0,255, 255}, color.RGBA{255,0,127, 255}, color.RGBA{255,255,10, 255})
  
  	// 生成字体 -- 将图片验证码, 展示到页面中.
  	http.HandleFunc("/r", func(w http.ResponseWriter, r *http.Request) {
  		img, str := cap.Create(4, captcha.NUM)
  		png.Encode(w, img)
  
  		println(str)
  	})
  
  	// 或者 自定固定的数据,来做图片内容.
  	http.HandleFunc("/c", func(w http.ResponseWriter, r *http.Request) {
  		str := "itcast"
  		img := cap.CreateCustom(str)
  		png.Encode(w, img)
  	})
  
  	// 启动服务
  	http.ListenAndServe(":8086", nil)
  }
  ```

#### 3.4 图片验证码集成到web端

1. 将 ImgTest.go 中的测试代码(最后两个回调函数不用）， 粘贴到 web/controller/user.go 中的 GetImageCd() 函数
2. 导入需要的包。 将 “comic.ttf” 文件 存放到 web/conf/ 中。 对应修改访问代码。
3. 浏览器中，“注册” 页面测试！

#### 3.5 web 端对接微服务实现

1. 拷贝密码本。 将 service 下的 proto/  拷贝 web/   下

2. 在 GetImageCd() 中 导入包，起别名：

   `getCaptcha "bj38web/web/proto/getCaptcha" `

3. 指定consul 服务发现：

   ```go
   // 指定 consul 服务发现
   consulReg := consul.NewRegistry()
   
   consulService := micro.NewService(
       micro.Registry(consulReg),
   )
   ```

4. 初始化客户端

   ```go
   microClient := getCaptcha.NewGetCaptchaService("getCaptcha", consulService.Client())
   
   ```

5. 调用远程函数

   ```go
   resp, err := microClient.Call(context.TODO(), &getCaptcha.Request{})
   if err != nil {
       fmt.Println("未找到远程服务...")
       return
   }
   ```

6. 将得到的数据,反序列化,得到图片数据

   ```go
   var img captcha.Image
   json.Unmarshal(resp.Img, &img)
   ```

7. 将图片写出到 浏览器.

   ```go
   png.Encode(ctx.Writer, img)
   ```

8. 测试：

   1. 启动 consul  ，  consul agent -dev
   2. 启动 service/  下的  main.go
   3. 启动 web/ 下的  main.go
   4. 浏览器中 192.168.IP: port/home    点击注册 查看图片验证码！

####  4.1 redis 操作回顾

1. 修改 配置文件。 /etc/redis/redis.conf . 
   - bind 地址。修改成当前主机地址。 —— 192.168.6.108 
2. port：
   - 6379
3. 开启 redis：
   - sudo  redis-server  /etc/redis/redis.conf
   - 验证 ： ps xua | grep redis  —— iP 和 port
4. 连接 redis ：
   - redis-cli -h 192.168.6.108 -p 6379
5. 查看所有：
   - keys *
6. 删除所有：
   - flushall
7. 添加一条：
   - set key  value  ——   set  hello  world
8. 获取一条：
   - get key 

#### 4.2 go 语言操作redis

- 从 redis.cn —— 客户端 —— go语言 —— 选择 redigo —— https://godoc.org/github.com/gomodule/redigo/redis#pkg-examples  查看 API

- 主要分为 3 类：

  1.  连接数据库。   
      - API文档中，所有以 Dial 开头。
  2.  操作数据库。   
      - Do() 函数【推荐】;  Send()函数, 需要配合Flush()、Receive() 3 个函数使用。
  3.  回复助手。   
      - 相当于 “类型断言”。根据使用的具体数据类型，选择调用。

- 添加测试案例：

  ```go
  package main
  
  import (
  	"github.com/gomodule/redigo/redis"
  	"fmt"
  )
  
  func main()  {
  	// 1. 链接数据库
  	conn, err := redis.Dial("tcp", "192.168.6.108:6379")
  	if err != nil {
  		fmt.Println("redis Dial err:", err)
  		return
  	}
  	defer conn.Close()
  
  	// 2. 操作数据库
  	reply, err := conn.Do("set", "itcast", "itheima")
  
  	// 3. 回复助手类函数. ---- 确定成具体的数据类型
  	r, e := redis.String(reply, err)
  
  	fmt.Println(r, e)
  }
  ```

  

#### 4.3 redis 存储 图片验证码的uuid 和 码值

思路： web 端产生一个图片的uuid 传给微服务端，微服务端产生一个码值，并同时把uuid作为key值，码值作为value值存入redis, 为的是后面验证输入的图片验证码是否与产生的验证码一致。

- 微服务端

  1. 修改 service/proto 中 getCaptcha.proto 的 Request 消息体，填加 uuid 成员。

     ```go
     message Request {
     	string uuid = 1;
     }
     ```

  2. 使用 make 命令，重新生成 getCaptcha.proto 对应的文件。

  3. 遵循 MVC 代码组织架构，在 service/getCaptcha/ 中 创建 model 目录

  4. 创建 modelFunc.go 文件 封装并实现 SaveImgCode() 函数：

     ```go
     // 存储图片id 到redis 数据库
     func SaveImgCode(code, uuid string) error {
     	// 1. 链接数据库
     	conn, err := redis.Dial("tcp", "192.168.6.108:6379")
     	if err != nil {
     		fmt.Println("redis Dial err:", err)
     		return err
     	}
     	defer conn.Close()
     
     	// 2. 写数据库  --- 有效时间 5 分钟
     	_, err = conn.Do("setex", uuid, 60*5, code)
     
     	return err  // 不需要回复助手!
     }
     ```

  5. 在 getCaptcha.go 文件的 Call() 方法中， cap.Create() 后， 调用 SaveImgCode() 传参。

     

- web端

  1. 修改密码本！因为 微服务端 修改了 proto/  , 添加消息体成员。 需要重新拷贝 proto/ 到web

  2. 修改 web/controller/user.go 中 Call() 方法传参。给 Request{} 初始化。

     ```go
     resp, err := microClient.Call(context.TODO(), &getCaptcha.Request{Uuid:uuid})
     if err != nil {
         fmt.Println("未找到远程服务...")
         return
     }
     ```

  3. 测试验证：

     1. 确认 consul 已经启动 
     2. 启动 service/    main.go 
     3. 启动 web/   main.go
     4. 浏览器：192.168.6.108:8080/home  --- 点击“注册”。看到 图片验证码。
     5. 查看 redis 数据库， keys * 能看到 图片验证吗，对应 uuid。 校验！！

#### 4.4 短信验证码准备工作

##### 开发者平台

- 中国移动 —— 短信验证码 业务 —— 下放到各个大公司（资质）。

- 常用平台：
  1. 聚合数据：
  2. 京东万象：
  3. 腾讯云：
  4. 阿里云（推荐）： 推荐，服务器。—— 生态好！ API接口丰富。 对应开发友好度高！

- 资料搜索。
  - 免费：—— github 
  - 收费：—— 各种开发者平台。刁钻、生僻领域。
  - 最后Google。



##### 注册阿里云账号

1. 通过实名认证
2. 开通短信验证码功能 —— 充值 2 元左右。
3. 申请 AccessKey
4. 申请签名。国内消息 —— 签名管理
5. 申请模板。国内消息 —— 模板管理
6. 测试使用  OpenAPI Explorer， 成功发送一条短信验证码！

#### 4.5 短信验证码



##### 测试短信验证码

1. 申请 阿里云账号、开通短信验证码功能、申请签名、申请模板、申请 AccessKey

2. 打开 OpenAPI  Explorer。

3. 选择 左侧 SendSms

4. 在中间位置依次填：华东1（杭州）、手机号、签名的名称、模板Code、{”code“:验证码}

5. 在右侧自动生成代码。 拷贝至，测试.go 程序中

6. 将 dysmsapi.NewClientWithAccessKey(） 函数的 ：\<accessKeyId> 和 \<accessSecret> 替换为我们申请到的 AccessKey 对应值。

7. 虚拟机安装的 SDK 版本，比从 OpenAPI  Explorer 工具拿到的代码版本低。需要添加一行代码：

   `request.Domain = "dysmsapi.aliyuncs.com" `

8. 运行 测试.go 程序。 —— 成功：在手机上收到 短信验证码。



##### 将短信验证码集成到项目

1. 修改 router 分组。

   ```go
   --- 在 web/main.go 中
   // 添加路由分组
   r1 := router.Group("/api/v1.0")
   {
       r1.GET("/session", controller.GetSession)
       r1.GET("/imagecode/:uuid", controller.GetImageCd)
       r1.GET("/smscode/:phone", controller.GetSmscd)
   }
   ```

2. 提取Get请求中的数据

   ```go
   --- 在 web/controller/user.go 中
   GET 请求 URL 格式： http://IP:port/资源路径?key=value&key=value&key=value...
   
   func GetSmscd(ctx *gin.Context)  {
   	// 获取短信验证码
   	phone := ctx.Param("phone")
   	// 拆分 GET 请求中 的 URL === 格式: 资源路径?k=v&k=v&k=v
   	imgCode := ctx.Query("text")
   	uuid := ctx.Query("id")
   
   	fmt.Println("---out---:", phone, imgCode, uuid)
   }
   ```

3. 封装实现 校验图片 验证码

   ```go
   --- 依据 MVC 代码架构。 创建 model/modelFunc.go
   
   // 校验图片验证码
   func CheckImgCode(uuid, imgCode string) bool {
   	// 链接 redis
   	conn, err := redis.Dial("tcp", "192.168.6.108:6379")
   	if err != nil {
   		fmt.Println("redis.Dial err:", err)
   		return false
   	}
   	defer conn.Close()
   
   	// 查询 redis 数据
   	code, err := redis.String(conn.Do("get", uuid))
   	if err != nil {
   		fmt.Println("查询错误 err:", err)
   		return false
   	}
   	
   	// 返回校验结果
   	return code == imgCode
   }
   ```

4. 根据校验结果，发送短信验证码

   ```go
   result := model.CheckImgCode(uuid, imgCode)
   if result {  // 校验成功
       // 发送短信验证码
       response, _ := client.SendSms(request)
       if response.IsSuccess() {
           // 发送短信验证码 成功
       } else {
           // 发送端验证码 失败.
       }
   } else {
       // 校验失败
   }
   ```

5. 发送短信验证码实现

   ```go
   client, _ := dysmsapi.NewClientWithAccessKey("cn-hangzhou", "LTAI4FgbQXjf117SX7E75Rmn", "6icOghQlhjevrTM5PxfiB8nDTxB9z6")
   
   request := dysmsapi.CreateSendSmsRequest()
   request.Scheme = "https"
   
   request.Domain = "dysmsapi.aliyuncs.com"  //域名  ---参考讲义补充!
   request.PhoneNumbers = phone
   request.SignName = "爱家租房网"
   request.TemplateCode = "SMS_183242785"
   
   // 生成一个随机 6 位数, 做验证码
   rand.Seed(time.Now().UnixNano())		// 播种随机数种子.
   // 生成6位随机数.
   smsCode := fmt.Sprintf("%06d", rand.Int31n(1000000))
   
   request.TemplateParam = `{"code":"` + smsCode + `"}`
   
   response, _ := client.SendSms(request)
   ```

   

6. 根据发送结果，给前端反馈消息

   ```go
   // 校验图片验证码 是否正确
   result := model.CheckImgCode(uuid, imgCode)
   if result {
       // 发送短信
       .....
       response, _ := client.SendSms(request)
       if response.IsSuccess() {
           // 发送短信验证码 成功
           resp["errno"] = utils.RECODE_OK
           resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
       } else {
           // 发送端验证码 失败.
           resp["errno"] = utils.RECODE_SMSERR
           resp["errmsg"] = utils.RecodeText(utils.RECODE_SMSERR)
       }
   } else {
       // 校验失败, 发送错误信息
       resp["errno"] = utils.RECODE_DATAERR
       resp["errmsg"] = utils.RecodeText(utils.RECODE_DATAERR)
   }
   ```

7. 测试：

   1. 启动 consul、启动微服务、启动 web 服务
   2. 浏览器 输入 手机号、图片验证码、短信验证码
   3. 成功：
      1. 手机会收到短信验证码。 
      2. 浏览器的 Network -- Response 中 会有 “成功”数据。
      3. redis 数据库中有  图片验证码 对应数据。



##### 短信验证码存入 Redis

- 使用 Redis 连接池！

  ```go
  // redis.Pool ——  Ctrl-B 查看 redis库， 连接池属性
  type Pool struct {
  	Dial func() (Conn, error)	// 连接数据库使用
  	。。。。
  	MaxIdle int		// 最大空闲数 == 初始化连接数
  	MaxActive int	// 最大存活数 > MaxIdle
  	IdleTimeout time.Duration	// 空闲超时时间。
  	。。。。
  	MaxConnLifetime time.Duration	// 最大生命周期。
  	。。。。
  }
  ```

- 连接池代码实现：

  ```go
  -- 在 model/modelFunc.go 中 创建redis连接池的 初始化函数。
  
  // 创建全局redis 连接池 句柄
  var RedisPool redis.Pool
  
  // 创建函数, 初始化Redis连接池
  func InitRedis()  {
  	RedisPool = redis.Pool{
  		MaxIdle:20,
  		MaxActive:50,
  		MaxConnLifetime:60 * 5,
  		IdleTimeout:60,
  		Dial: func() (redis.Conn, error) {
  			return redis.Dial("tcp", "192.168.6.108:6379")
  		},
  	}
  }
  
  ---在 web/main.go 中，使用该 InitRedis() , 在项目启动时，自动初始化连接池！
  ```

- 修改了 CheckImgCode() , 使用连接池。

- 实现 SaveSmsCode() 函数，将数据存入 redis

  ```go
  // 链接 Redis --- 从链接池中获取一条链接
  conn := RedisPool.Get()
  defer conn.Close()
  
  // 存储短信验证码到 redis 中
  _, err := conn.Do("setex", phone+"_code", 60 * 3, code)
  ```

  

- 在 GetSmscd() 函数中，在 短信验证码发送成功之后：

  ```go
  response, _ := client.SendSms(request)
  if response.IsSuccess() {
      // 发送短信验证码 成功
      resp["errno"] = utils.RECODE_OK
      resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
  
      // 将 电话号:短信验证码 ,存入到 Redis 数据库
      err := model.SaveSmsCode(phone, smsCode)
      if err != nil {
          fmt.Println("存储短信验证码到redis失败:", err)
          resp["errno"] = utils.RECODE_DBERR
          resp["errmsg"] = utils.RecodeText(utils.RECODE_DBERR)
      }
  
  } else {
      // 发送端验证码 失败.
      resp["errno"] = utils.RECODE_SMSERR
      resp["errmsg"] = utils.RecodeText(utils.RECODE_SMSERR)
  }
  ```

- 测试：

  - 保证 consul、service的 main.go  web的main.go 、redis服务  启动。
  - 浏览器 —— 注册 —— 获取短信验证码
  - 成功：
    - 手机收到 验证码
    - Network —— Respose —— “成功”
    - 打开 redis 数据库。 多 “手机号_code” 为 key 值的 一条数据。
      - get “手机号_code”  取值 == 手机收到 验证码



##### 短信验证码 分离成微服务

1. 创建 微服务。 将 “登录”、“短信验证”、“注册” 使用 user 微服务 实现。

   `micro new --type srv bj38web/service/user`

   ```shell
   Creating service go.micro.srv.user in /home/itcast/workspace/go/src/bj38web/service/user
   .
   ├── main.go
   ├── plugin.go
   ├── handler
   │   └── user.go
   ├── subscriber
   │   └── user.go
   ├── proto/user
   │   └── user.proto
   ├── Dockerfile
   ├── Makefile
   ├── README.md
   └── go.mod
   
   download protobuf for micro:
   
   brew install protobuf
   go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
   go get -u github.com/micro/protoc-gen-micro
   
   compile the proto file user.proto:
   
   cd /home/itcast/workspace/go/src/bj38web/service/user
   protoc --proto_path=.:$GOPATH/src --go_out=. --micro_out=. proto/user/user.proto
   
   ```

   

2. 修改密码本 —— proto文件。

   ```protobuf
   // 修改 Call --- SendSms。  删除 Stream、PingPong 函数。
   // 删除 除 Request、Response 之外的其他 message 消息体。
   // 根据传入、传出修改 Request、Response 
   syntax = "proto3";
   
   package go.micro.srv.user;
   
   service User {
   	rpc SendSms(Request) returns (Response) {}
   }
   message Request {
   	string phone = 1;
   	string imgCode = 2;
   	string uuid = 3;
   }
   message Response {
   	string errno = 1;
   	string errmsg = 2;
   }
   ```

3. 编译 proto文件，生成 2 个新文件 xxx.micro.go 和 xxx.pb.go ，用于 grpc 远程调用！

   ```shell
   make proto
   ```

4. 修改 service/user/main.go

   ```go
   import (
   	"github.com/micro/go-micro/util/log"
   	"github.com/micro/go-micro"
   	"bj38web/service/user/handler"
   	user "bj38web/service/user/proto/user"
   )
   
   func main() {
   	// New Service
   	service := micro.NewService(
   		micro.Name("go.micro.srv.user"),
   		micro.Version("latest"),
   	)
   
   	// Register Handler
   	user.RegisterUserHandler(service.Server(), new(handler.User))
   
   	// Run service
   	if err := service.Run(); err != nil {
   		log.Fatal(err)
   	}
   }
   ```

   

5. 修改 service/user/handler/user.go

   ```go
   import (
   	"context"
   	user "bj38web/service/user/proto/user"
   )
   
   type User struct{}
   
   // Call is a single request handler called via client.Call or the generated client code
   func (e *User) SendSms(ctx context.Context, req *user.Request, rsp *user.Response) error {
   
   	return nil
   }
   ```

   

6. 修改服务发现 mdns ——> consul:

   ```go
   --- 在 /service/user/main.go 中
   
   // 初始化 Consul
   consulReg := consul.NewRegistry()
   
   // New Service  -- 指定 consul
   service := micro.NewService(
       micro.Address("192.168.6.108:12342"),
       micro.Name("go.micro.srv.user"),
       micro.Registry(consulReg),
       micro.Version("latest"),
   )
   ```

7. 移植 web/controller/user.go 中 “发送短信验证码” 代码，到 service/user/handler/user.go 中，实现微服务版的 短信验证码功能。

   ```go
   --- 修改对应的 import 包。
   --- 从web/下 拷贝 model/ 和 utils/ 包  到  service/user/ 中。
   --- 修改代码，使用 SendSms() 的 req， rsp 进行传参！
   // Call is a single request handler called via client.Call or the generated client code
   func (e *User) SendSms(ctx context.Context, req *user.Request, rsp *user.Response) error {
   
   	// 校验图片验证码 是否正确
   	result := model.CheckImgCode(req.Uuid, req.ImgCode)
   	if result {
   		// 发送短信
   		client, _ := dysmsapi.NewClientWithAccessKey("cn-hangzhou", "LTAI4FgbQXjf117SX7E75Rmn", "6icOghQlhjevrTM5PxfiB8nDTxB9z6")
   
   		request := dysmsapi.CreateSendSmsRequest()
   		request.Scheme = "https"
   
   		request.Domain = "dysmsapi.aliyuncs.com" //域名  ---参考讲义补充!
   		request.PhoneNumbers = req.Phone
   		request.SignName = "爱家租房网"
   		request.TemplateCode = "SMS_183242785"
   
   		// 生成一个随机 6 位数, 做验证码
   		rand.Seed(time.Now().UnixNano()) // 播种随机数种子.
   		// 生成6位随机数.
   		smsCode := fmt.Sprintf("%06d", rand.Int31n(1000000))
   
   		request.TemplateParam = `{"code":"` + smsCode + `"}`
   
   		response, _ := client.SendSms(request)
   		if response.IsSuccess() {
   			// 发送短信验证码 成功
   			rsp.Errno = utils.RECODE_OK
   			rsp.Errmsg = utils.RecodeText(utils.RECODE_OK)
   
   			// 将 电话号:短信验证码 ,存入到 Redis 数据库
   			err := model.SaveSmsCode(req.Phone, smsCode)
   			if err != nil {
   				fmt.Println("存储短信验证码到redis失败:", err)
   				rsp.Errno = utils.RECODE_DBERR
   				rsp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
   			}
   
   		} else {
   			// 发送端验证码 失败.
   			rsp.Errno = utils.RECODE_SMSERR
   			rsp.Errmsg = utils.RecodeText(utils.RECODE_SMSERR)
   		}
   
   	} else {
   		// 校验失败, 发送错误信息
   		rsp.Errno = utils.RECODE_DATAERR
   		rsp.Errmsg = utils.RecodeText(utils.RECODE_DATAERR)
   	}
   	return nil
   }
   
   ```

8. 在微服务项目中，初始化 Redis 连接池！service/user/main.go 中 调用 InitRedis() 



##### 实现短信验证码 客户端调用

1.  拷贝密码本。

- 将 service/user/proto/user/  拷贝至  web/proto/   —— 与 getCaptcha 同级。

2. 先导包，起别名。

   - `userMicro "bj38web/web/proto/user"   // 给包起别名`

3. 指定Consul服务发现。 mdns --> consul

   ```go
   consulReg := consul.NewRegistry()
   consulService := micro.NewService(
       micro.Registry(consulReg),
   )
   ```

4. 初始化客户端

   ```go
   microClient := userMicro.NewUserService("go.micro.srv.user", consulService.Client())
   参1： service/user/main.go 中 指定的 服务名。
   ```

5. 调用远程函数, 封装调用结果程json 发送给浏览器

   ```go
   resp, err := microClient.SendSms(context.TODO(), &userMicro.Request{Phone:phone, ImgCode:imgCode,Uuid:uuid})
   if err != nil {
       fmt.Println("调用远程函数 SendSms 失败:", err)
       return
   }
   
   // 发送校验结果 给 浏览器
   ctx.JSON(http.StatusOK, resp)
   ```

6. 测试：

   1. 启动 consul  。  consul agent -dev
   2. 启动 service/getCaptcha/ 下的 main.go
   3. 启动 service/user/ 下的 main.go
   4. 启动 web/  下的 main.go
   5. 浏览器 ：192.168.6.108:8080/home  -- > 注册流程
      - 成功：
        1. 手机获取到短信验证码
        2. NetWork -- Respose 收到 errno = 0  errmsg = 成功！
        3. 打开 redis 数据库。 验证 图片验证码、短信验证码    码值。 



##### 获取数据 绑定数据

- 前端传递数据种类：
  1. form表单：数据为 form data
  2. ajax(阿贾克斯)： 数据为 json 格式。 体现成 —— Request Payload
- 默认 postForm() 方法 只能获取 form 表单传递的数据。
- 针对 Request Payload 数据形式，需要 使用 “数据绑定“ 来获取传递的数据。
  - `ctx.Bind()` 将 数据绑定到对象中。

##### 获取数据

浏览器 ：注册流程， —— Network 中 Headers 中 获取

-  http请求方法：Post
-  Request Payload：  mobile、password、sms_code 作为 key 

```go
-- web/main.go 中 添加

r1.POST("/users", controller.PostRet)

-- web/controller/user.go 中添加、实现 
// 发送注册信息
func PostRet(ctx *gin.Context) {
	/*	mobile := ctx.PostForm("mobile")
		pwd := ctx.PostForm("password")
		sms_code := ctx.PostForm("sms_code")
		fmt.Println("m = ", mobile, "pwd = ", pwd, "sms_code = ",sms_code)
	*/
	// 获取数据
	var regData struct {
		Mobile   string `json:"mobile"`
		PassWord string `json:"password"`
		SmsCode  string `json:"sms_code"`
	}
	ctx.Bind(&regData)
	fmt.Println("获取到的数据为:", regData)
}
```

#### 5 微服务实现用户注册

##### 微服务端

1. 修改密码本 —— proto 文件

   ```go
   syntax = "proto3";
   
   package go.micro.srv.user;
   
   service User {
   	rpc SendSms(Request) returns (Response) {};
   	rpc Register(RegReq) returns (Response) {};   // 注册用户
   }
   
   message RegReq {
       string mobile = 1;
       string password = 2;
       string sms_code = 3;
   }
   
   message Request {
   	string phone = 1;
   	string imgCode = 2;
   	string uuid = 3;
   }
   
   message Response {
   	string errno = 1;
   	string errmsg = 2;
   }
   ```

2. make 编译生成 xxx.micro.go 文件。

3. 修改 service/user/main.go    --- 没有需要修改的

4. 修改 handler/user.go 

   ```go
   // 添加 方法
   func (e *User) Register(ctx context.Context, req *user.RegReq, rsp *user.Response) error {
   	return nil
   }
   ```

5. 需要操作，MySQL数据库， 拷贝 web/model/model.go 到  微服务项目中。

6. 在 service/user/model/modelFunc.go 中 添加 校验短信验证码函数实现

   ```go
   // 校验短信验证码
   func CheckSmsCode(phone, code string) error {
   	// 链接redis
   	conn := RedisPool.Get()
   
   	// 从 redis 中, 根据 key 获取 Value --- 短信验证码  码值
   	smsCode, err := redis.String(conn.Do("get", phone+"_code"))
   	if err != nil {
   		fmt.Println("redis get phone_code err:", err)
   		return err
   	}
   	// 验证码匹配  失败
   	if smsCode != code {
   		return errors.New("验证码匹配失败!")
   	}
   	// 匹配成功!
   	return nil
   }
   ```

7. service/user/model/modelFunc.go 中, 添加函数RegisterUser， 实现 用户注册信息，写入MySQL数据库

   ```go
   // 注册用户信息,写 MySQL 数据库.
   func RegisterUser(mobile, pwd string) error {
   	var user User
   	user.Name = mobile		// 默认使用手机号作为用户名
   
   	// 使用 md5 对 pwd 加密
   	m5 := md5.New()			// 初始md5对象
   	m5.Write([]byte(pwd))			// 将 pwd 写入缓冲区
   	pwd_hash := hex.EncodeToString(m5.Sum(nil))	// 不使用额外的秘钥
   
   	user.Password_hash = pwd_hash
   
   	// 插入数据到MySQL
   	return GlobalConn.Create(&user).Error
   }
   ```

8. 完成  Register  函数 实现

   ```go
   func (e *User) Register(ctx context.Context, req *user.RegReq, rsp *user.Response) error {
   	// 先校验短信验证码,是否正确. redis 中存储短信验证码.
   	err := model.CheckSmsCode(req.Mobile, req.SmsCode)
   	if err == nil {
   
   		// 如果校验正确. 注册用户. 将数据写入到 MySQL数据库.
   		err = model.RegisterUser(req.Mobile, req.Password)
   		if err != nil {
   			rsp.Errno = utils.RECODE_DBERR
   			rsp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
   		} else {
   			rsp.Errno = utils.RECODE_OK
   			rsp.Errmsg = utils.RecodeText(utils.RECODE_OK)
   		}
   	} else {  // 短信验证码错误
   		rsp.Errno = utils.RECODE_DATAERR
   		rsp.Errmsg = utils.RecodeText(utils.RECODE_DATAERR
   	}
   	return nil
   }
   ```


##### web 端

1. 拷贝密码本 ——proto

2. 创建 web/utils/utils.go 文件， 封装 函数实现 初始 consul 客户端代码

   ```go
   // 初始化micro
   func InitMicro() micro.Service {
   	// 初始化客户端
   	consulReg := consul.NewRegistry()
   	return micro.NewService(
   		micro.Registry(consulReg),
   	)
   }
   ```

3. 实现 web/controller/user.go 中的 PostRet 函数

   ```go
   // 发送注册信息
   func PostRet(ctx *gin.Context) {
   	// 获取数据
   	var regData struct {
   		Mobile   string `json:"mobile"`
   		PassWord string `json:"password"`
   		SmsCode  string `json:"sms_code"`
   	}
   	ctx.Bind(&regData)
   
   	// 初始化consul
   	microService := utils.InitMicro()
       
       // 初始化客户端
   	microClient := userMicro.NewUserService("go.micro.srv.user", microService.Client())
   
   	// 调用远程函数
   	resp, err := microClient.Register(context.TODO(), &userMicro.RegReq{
   		Mobile:regData.Mobile,
   		SmsCode:regData.SmsCode,
   		Password:regData.PassWord,
   	})
   	if err != nil {
   		fmt.Println("注册用户, 找不到远程服务!", err)
   		return
   	}
   	// 写给浏览器
   	ctx.JSON(http.StatusOK, resp)
   }
   ```

4. 测试：

   - consul 启动
   - getCaptcha 服务启动  --- 12341
   - user 服务启动  --- 12342
   - web 启动  --- 8080
   - 浏览器测试，注册流程。
     - 成功：
       - 界面跳转。
       - 查询 MySQL数据库， 多一条用户信息。

#### 6 获取地域信息

##### 导入 SQL脚本

1. 将 home.sql 保存至 Linux 系统。建议放 家目录。
2. 登录 MySQL数据库。选择数据库： use  search_house;
3. 执行 source  /home/itcast/home.sql  ——  运行 脚本文件。向表插入数据。



##### web端实现

1. 在 web/main.go 中，添加路由， 设置回调。

   ```go
   r1.GET("/areas", controller.GetArea)
   ```

2. 在 web/controller/use.go 中， 添加 GetArea() 函数。

   ```go
   func GetArea(ctx *gin.Context)  {
   }
   ```

3. 从数据库获取数据，提高用户感受的常见方法：先查缓存， 缓存没有查MySQL， 写入redis缓存。

   ```go
   // 测试实现：
   // 获取地域信息
   func GetArea(ctx *gin.Context)  {
   	// 先从MySQL中获取数据.
   	var areas []model.Area
   
   	model.GlobalConn.Find(&areas)
   
   	// 再把数据写入到 redis 中.
   	conn := model.RedisPool.Get()		// 获取链接
   	conn.Do("set", "areaData", areas)
   
   	resp := make(map[string]interface{})
   
   	resp["errno"] = "0"
   	resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
   	resp["data"] = areas
   
   	ctx.JSON(http.StatusOK, resp)
   }
   ```

4. 测试：登录 redis ，指定 --raw 参数，显示中文。

   ```shell
   itcast@ubuntu:~$ redis-cli -h 192.168.6.108 -p 6379 --raw
   192.168.6.108:6379> keys *
   areaData
   hello
   itcast
   192.168.6.108:6379> get areaData
   [{1 东城区 []} {2 西城区 []} {3 朝阳区 []} {4 海淀区 []} {5 昌平区 []} {6 丰台区 []} {7 房山区 []} {8 通州区 []} {9 顺义区 []} {10 大兴区 []} {11 怀柔区 []} {12 平谷区 []} {13 密云区 []} {14 延庆区 []} {15 石景山区 []}]
   192.168.6.108:6379> 
   
   ```

5. 思考：按如上方法存储数据到 Redis 中 `conn.Do("set", "areaData", areas)`， 将来 使用 Do 获取数据时！不好获取！没有对应的 回复助手函数来完成 “类型断言”。 —— 重新选择 存储 redis 的方法: 将 数据转换成 josn 字节流存储。 

![1582438662615](图片/1582438662615.png)

6. 重新实现获取地域信息， 没数据，读MySQL，写redis；有数据，直接读 redis

   - 强调：写入 Redis 中的数据 —— 序列化后的字节流数据。

   ```go
   // 获取地域信息
   func GetArea(ctx *gin.Context) {
   	// 先从MySQL中获取数据.
   	var areas []model.Area
   
   	// 从缓存redis 中, 获取数据
   	conn := model.RedisPool.Get()
   	// 当初使用 "字节切片" 存入, 现在使用 切片类型接收
   	areaData, _ := redis.Bytes(conn.Do("get", "areaData"))
   	// 没有从 Redis 中获取到数据
   	if len(areaData) == 0 {
   
   		fmt.Println("从 MySQL 中 获取数据...")
   		model.GlobalConn.Find(&areas)
   		// 把数据写入到 redis 中. , 存储结构体序列化后的 json 串
   		areaBuf, _ := json.Marshal(areas)
   		conn.Do("set", "areaData", areaBuf)
   
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
   ```

7. 测试：

   1. GoLand 中 借助输出，看到 展示的数据来源。 



#### 7 Cookie 和 Session  实现

##### Cookie 和 Session简介

- http协议，有 3 个版本：
  - http/1.0 版：无状态，短连接。
  - http/1.1 版：可以记录状态。—— 默认支持。
  - http/2.0 版：可以支持长连接。 协议头：Connection: keep-alive 。

##### Cookie

- 最早的 http/1.0 版，提供 Cookie 机制， 但是没有 Session。
- Cookie 作用：一定时间内， 存储用户的连接信息。如：用户名、登录时间 ... 不敏感信息。
- Cookie 出身：http自带机制。Session不是！
- Cookie 存储：Cookie 存储在 客户端 (浏览器) 中。—— 浏览器可以存储数据。少
  - 存储形式：key - value
  - 可以在浏览器中查看。
  - Cookie 不安全。直接将数据存储在浏览器上。

##### Session

- ”会话“：在一次会话交流中，产生的数据。不是http、浏览器自带。
- Session 作用：一定时间内， 存储用户的连接信息。
- Session 存储：在服务器中。一般为 临时 Session。—— 会话结束 (浏览器关闭) ， Session被干掉！

##### 对比 Cookie 和 Session

1.  Cookie 存储在 浏览器， 在哪生成呢？
2.  Session 存储在 服务器，在哪生成呢？
3.  什么时候生成Cookie ， 什么时候生成 Session？

![1582441308868](图片/1582441308868.png)

##### Cookie操作

##### 设置Cookie

```go
func (c *Context) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) 
name: 名称。 相当于 key
value：内容。
maxAge：最大生命周期。
	 = 0 : 表示没指定该属性。
	 < 0 ：表示删除。 ---- // 删除Cookie 的操作， 可以使用 该属性实现。
	 > 0 ：指定生命周期。 单位：s
path：路径。—— 通常传""
domain：域名。 IP地址。
secure：设置是否安全保护。true：不能在 地址栏前，点击查看。 可以使用 F12 查看。
					   false：能在 地址栏前，点击查看。
httpOnly：是否只针对http协议。
```

测试案例：

```go
package main

import "github.com/gin-gonic/gin"

func main()  {
	router := gin.Default()

	router.GET("/test", func(context *gin.Context) {
		// 设置 Cookie
		//context.SetCookie("mytest", "chuanzhi", 60*60, "", "", true, true)
        //context.SetCookie("mytest", "chuanzhi", 60*60, "", "", false, true)
		context.SetCookie("mytest", "chuanzhi", 0, "", "", false, true)
		context.Writer.WriteString("测试 Cookie ...")
	})

	router.Run(":9999")
}
```



##### 获取Cookie

```go
// 获取Cookie
cookieVal, _ := context.Cookie("mytest")

fmt.Println("获取的Cookie 为:", cookieVal)
```



##### Session 操作

- gin 框架， 默认不支持Session功能。要想在 gin 中使用 Session，需要添加插件！

- gin 框架中的 “插件”  —— 中间件 —— gin MiddleWare

- 去 github 搜索，gin Session 可以得到：https://github.com/gin-contrib/sessions

- 安装 Session 插件。

- ```
  $ go get github.com/gin-contrib/sessions
  ```



##### 设置session

- 容器的初始化：

  ```go
  func NewStore(size int, network, address, password string, keyPairs ...[]byte) (Store, error)
  size:容器大小。
  network：协议
  address：IP：port
  password：使用redis做容器使用的密码。 没有特殊设定，传 “”
  []byte(“secret”)： 加密密钥！
  ```

- 使用容器：

  ```go
  func (engine *Engine) Use(middleware ...HandlerFunc) IRoutes {}
  router.Use(sessions.Sessions("mysession", store))
  ```

测试案例：

```go
package main

import (
	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-contrib/sessions"
)

func main()  {
	router := gin.Default()

	// 初始化容器.
	store, _ := redis.NewStore(10, "tcp", "192.168.6.108:6379", "", []byte("bj38"))

	// 使用容器
	router.Use(sessions.Sessions("mysession", store))

	router.GET("/test", func(context *gin.Context) {
		// 调用session, 设置session数据
		s := sessions.Default(context)
		// 设置session
		s.Set("itcast", "itheima")
		// 修改session时, 需要Save函数配合.否则不生效
		s.Save()

		context.Writer.WriteString("测试 Session ...")
	})

	router.Run(":9999")
}
```



##### 获取session

```go
// 建议：不要修改 session属性，使用默认属性。
v := s.Get("itcast")
fmt.Println("获取 Session:", v.(string))
```



#### 8  实现用户登录

1. 浏览器 访问 ： 192.168.6.108:8080/home   点击登录按钮。  跳“登录页面”。 输手机号、输密码，登录

2. 在 Name的 General  和  Request Payload 中获取到 路由 和 方法以及数据信息。

3. web/main.go 添加：

   ```go
   r1.POST("/sessions", controller.PostLogin)    // 注意 “s”
   ```

4. web/controller/user.go 增加函数

   ```go
   // 处理登录业务
   func PostLogin(ctx *gin.Context) {
   }
   ```

5. 实现 PostLogin 函数

   1. 获取数据。 因为数据来自  Request Payload , 所以：需要通过“Bind”来获取输入数据

      ```go
      var loginData struct {
          Mobile   string `json:"mobile"`
          PassWord string `json:"password"`
      }
      ctx.Bind(&loginData)
      ```

   2. web/model/modelFunc.go  创建函数, 处理登录业务，根据手机号/密码 获取用户名

      ```go
      // 处理登录业务,根据手机号/密码 获取用户名
      func Login(mobile, pwd string) (string, error) {
      
      	var user User
      
      	// 对参数 pwd 做md5 hash
      	m5 := md5.New()
      	m5.Write([]byte(pwd))
      	pwd_hash := hex.EncodeToString(m5.Sum(nil))
      
      	err := GlobalConn.Where("mobile = ?", mobile).Select("name").
      		Where("password_hash = ?", pwd_hash).Find(&user).Error
      
      	return user.Name, err
      }
      ```

6. 获取数据库数据，查询是否和输入数据匹配	

   ```go
   userName, err  := model.Login(loginData.Mobile, loginData.PassWord)
   resp := make(map[string]interface{})
   if err == nil {
       // 登录成功
       resp["errno"] = utils.RECODE_OK
       resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
   
       // 将 登录状态保存到 session 中 
   
   } else {
       // 登录失败.
       resp["errno"] = utils.RECODE_LOGINERR
       resp["errmsg"] = utils.RecodeText(utils.RECODE_LOGINERR)
   }
   ```

   

7. 在 web/main.go 中 ， 初始化容器， 使用容器

   ```go
   // 初始化容器
   store, _ := redis.NewStore(10, "tcp", "192.168.6.108:6379", "", []byte("bj38"))	
   
   // 使用容器
   router.Use(sessions.Sessions("mysession", store))
   ```

   

8. 完整实现 ，处理登录业务  PostLogin() 函数

   ```go
   // 处理登录业务
   func PostLogin(ctx *gin.Context) {
   	// 获取前端数据
   	var loginData struct {
   		Mobile   string `json:"mobile"`
   		PassWord string `json:"password"`
   	}
   	ctx.Bind(&loginData)
   
   	resp := make(map[string]interface{})
   
   	//获取 数据库数据, 查询是否和数据的数据匹配
   	userName, err := model.Login(loginData.Mobile, loginData.PassWord)
   	if err == nil {
   		// 登录成功!
   		resp["errno"] = utils.RECODE_OK
   		resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
   
   		// 将 登录状态, 保存到Session中
   		s := sessions.Default(ctx)	  // 初始化session
   		s.Set("userName", userName)   // 将用户名设置到session中.
   		s.Save()
   
   	} else {
   		// 登录失败!
   		resp["errno"] = utils.RECODE_LOGINERR
   		resp["errmsg"] = utils.RecodeText(utils.RECODE_LOGINERR)
   	}
   
   	ctx.JSON(http.StatusOK, resp)
   }
   
   ```

9. 测试：

   1. go  run   web/main.go 即可！ 其他的不用启动！
   2. 浏览器， 192.168.6.108:8080/home  ——> 登录 ——> 输入用户名、密码 ——> 登录！
   3. 看不到变化，是因为：我们写的第一个 Session 相关函数 GetSession()，里面 直接发送的假数据，并没有真正获取 Session。现在我们有真正Session了。
   4. 请大家尝试修改实现 GetSession() ！最终能在浏览器中看到 登录变化。

#### 9 用户基本信息查看

##### 用户登录

##### 修改 GetSession 方法

- 之前实现的 web/controller/user.go 中的 GetSession() 方法，是一个伪实现，没有真正的获取 Session。 

- 从 容器中， 真正的获取 Session，展示给浏览器。

  ```go
  func GetSession(ctx *gin.Context)  {
  	resp := make(map[string]interface{})
  
  	// 获取 Session 数据
  	s := sessions.Default(ctx)		// 初始化 Session 对象
  	userName := s.Get("userName")
  
  	// 用户没有登录.---没存在 MySQL中, 也没存在 Session 中
  	if userName == nil {
  		resp["errno"] = utils.RECODE_SESSIONERR
  		resp["errmsg"] = utils.RecodeText(utils.RECODE_SESSIONERR)
  	} else {
  		resp["errno"] = utils.RECODE_OK
  		resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
  
  		var nameData struct{
  			Name string `json:"name"`
  		}
  		nameData.Name = userName.(string)		// 类型断言
  		resp["data"] = nameData
  	}
  
  	ctx.JSON(http.StatusOK, resp)
  }
  ```

  

##### 退出登录

- 将 Session 数据删除！保存修改！

  ```go
  // 退出登录
  func DeleteSession(ctx *gin.Context)  {
  	resp := make(map[string]interface{})
  
  	// 初始化 Session 对象
  	s := sessions.Default(ctx)
  	// 删除 Session 数据
  	s.Delete("userName")		// 没有返回值
  	// 必须使用 Save 保存
  	err := s.Save()				// 有返回值
  
  	if err != nil {
  		resp["errno"] = utils.RECODE_IOERR	// 没有合适错误,使用 IO 错误!
  		resp["errmsg"] = utils.RecodeText(utils.RECODE_IOERR)
  
  	} else {
  		resp["errno"] = utils.RECODE_OK
  		resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
  	}
  	ctx.JSON(http.StatusOK, resp)
  }
  ```

  

##### 获取用户基本信息

1. 登录成功后，点击 “用户名(手机号)” 进入 用户页面。 尚未显示用户 “用户名”、“手机号” 用户信息。

2. web/main.go 中 添加 路由回调。—— 参考 《接口文档.doc》

   ```go
   r1.GET("/user", controller.GetUserInfo)
   ```

3. web/controller/user.go 中添加 GetUserInfo 函数， 并实现

   ```go
   // 获取用户基本信息
   func GetUserInfo(ctx *gin.Context)  {
   }
   ```

4. 实现：

   1. 准备 存储输出结果(成功、失败) 的 数据结构：

      ```go
      resp := make(map[string]interface{})
      defer ctx.JSON(http.StatusOK, resp)
      ```

   2. 获取 Session, 得到 当前 用户信息

      ```go
      s := sessions.Default(ctx)  // Session 初始化
      userName := s.Get("userName")	// 根据key 获取Session
      
      if userName == nil {		// 用户没登录, 但进入该页面, 恶意进入.
          resp["errno"] = utils.RECODE_SESSIONERR
          resp["errmsg"] = utils.RecodeText(utils.RECODE_SESSIONERR)
          return		// 如果出错, 报错, 退出
      }
      ```

   3. 访问 MySQL数据库， 按照 userName ， 提取用户信息

      ```go
      --- 在 web/model/modleFunc.go 中 创建 GetUserInfo() 函数，并实现。
      // 获取用户信息
      func GetUserInfo(userName string) (User, error) {
      	// 实现SQL: select * from user where name = userName;
      	var user User
      	err := GlobalConn.Where("name = ?", userName).First(&user).Error
      	return user, err
      }
      
      上述函数，也可以写成。---- go语法写。
      func GetUserInfo(userName string) (user User, err error) {
      	// 实现SQL: select * from user where name = userName;
      	err = GlobalConn.Where("name = ?", userName).First(&user).Error
      	return 
      }
      ```

   4. 调用 GetUserInfo 获取用户信息。

      ```go
      // 根据用户名, 获取 用户信息  ---- 查 MySQL 数据库  user 表.
      user, err := model.GetUserInfo(userName.(string))  // 类型断言，传参
      if err != nil {
          resp["errno"] = utils.RECODE_DBERR
          resp["errmsg"] = utils.RecodeText(utils.RECODE_DBERR)
          return		// 如果出错, 报错, 退出
      }
      ```

   5. 参照 《接口文档.doc》 中 发送成功消息字段，实现

      ```go
      resp["errno"] = utils.RECODE_OK
      resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
      
      temp := make(map[string]interface{})
      temp["user_id"] = user.ID
      temp["name"] = user.Name
      temp["mobile"] = user.Mobile
      temp["real_name"] = user.Real_name
      temp["id_card"] = user.Id_card
      temp["avatar_url"] = user.Avatar_url
      
      resp["data"] = temp
      
      // 前面已经添加 defer ctx.JSON() 
      ```

5. 测试：

   1. 重启  web/main.go 
   2. 刷新 浏览器，看到：用户名：xxxx  手机号 ： xxxxx   。 



##### 更新用户名

1. web/main.go 添加 路由回调

   ```go
   r1.PUT("/user/name", controller.PutUserInfo)
   ```

2. web/controller/user.go 定义，实现函数 PutUserInfo

   ```go
   func PutUserInfo(ctx *gin.Context) {
   	// 获取当前用户名
   	// 获取新用户名
   	// 更新用户名
   }
   ```

3. 获取 当前用户名

   ```go
   s := sessions.Default(ctx)			// 初始化Session 对象
   userName := s.Get("userName")
   ```

   

4. 获取 新用户名

   ```go
   // 获取新用户名		---- 处理 Request Payload 类型数据. Bind()
   var nameData struct {
       Name string `json:"name"`
   }
   ctx.Bind(&nameData)
   ```

   

5. 更新 用户名

   1. 更新 MySQL 数据库

      ```go 
      --- 在 web/model/modelFunc.go 实现
      // 更新用户名
      func UpdateUserName(newName, oldName string) error {
      	// update user set name = 'itcast' where name = 旧用户名
      	return GlobalConn.Model(new(User)).Where("name = ?", oldName).Update("name", newName).Error
      }
      ```

   2. 调用 model 中的 UpdateUserName() 函数

      ```go
      // 更新用户名
      resp := make(map[string]interface{})
      defer ctx.JSON(http.StatusOK, resp)
      
      // 更新数据库中的 name
      err := model.UpdateUserName(nameData.Name, userName.(string))
      if err != nil {
          resp["errno"] = utils.RECODE_DBERR
          resp["errmsg"] = utils.RecodeText(utils.RECODE_DBERR)
          return
      }
      ```

   3. 更新 Session 

      ```go
      // 更新 Session 数据
      s.Set("userName", nameData.Name)
      err = s.Save()		// 必须保存
      if err != nil {
          resp["errno"] = utils.RECODE_SESSIONERR
          resp["errmsg"] = utils.RecodeText(utils.RECODE_SESSIONERR)
          return
      }
      resp["errno"] = utils.RECODE_OK
      resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
      resp["data"] = nameData
      ```

6. 测试：

   - 修改以后的用户名，保存在 Session中， 后退页面也能看到！



#### 10 中间件 (middleWare)

- 中间件，对以后的路由全部生效。
  - 设置好中间件以后，所有的路由都会使用这个中间件。
  - 设置以前的路由，不生效。



##### 什么是 “中间件”：

- 早期：
  - 用于 系统 和 应用之间。
  - 中间件： 内核 —— 中间件 ——  用户应用

- 现在：
  - 用于 两个模块之间的 功能 软件(模块)
  - 中间件：—— 承上启下。  前后台开发： 路由 ——> 中间件 (起过滤作用) ——> 控制器
  - 特性：对 “中间件”指定位置 ， 以下的路由起作用！以上的，作用不到。

##### 中间件类型

- gin 框架规定：中间件类型为：gin.HandlerFunc 类型。

- gin.HandlerFunc 类型。就是 ：

  ```go
  func (c *gin.Context) { 
      
  }
  ```

  ```go
  // 示例：
  func Logger() gin.HandlerFunc {
      return func (c *gin.Context) {   
      }
  }
  r.Use(Logger())		// 传 “中间件” 做参数。
  ```

##### 中间件测试

- 中间件使用的 3 个知识：

##### Next()

- 表示，跳过当前中间件剩余内容， 去执行下一个中间件。 当所有操作执行完之后，以出栈的执行顺序返回，执行剩余代码。

- ```go
    
  // 创建中间件
  func Test1(ctx *gin.Context)  {
  	fmt.Println("1111")
  	ctx.Next()
  	fmt.Println("4444")
  }
  // 创建 另外一种格式的中间件.
  func Test2() gin.HandlerFunc {
  	return func(context *gin.Context) {
  		fmt.Println("3333")
  		context.Next()
  		fmt.Println("5555")
  	}
  }
  func main()  {
  	router := gin.Default()
  
  	// 使用中间件
  	router.Use(Test1)
  	router.Use(Test2())
  
  	router.GET("/test", func(context *gin.Context) {
  		fmt.Println("2222")
  		context.Writer.WriteString("hello world!")
  	})
  
  	router.Run(":9999")
  }
  ```

##### return 

- 终止执行当前中间件剩余内容，执行下一个中间件。 当所有的函数执行结束后，以出栈的顺序执行返回，但，不执行return后的代码！

  ```go
  // 创建中间件
  func Test1(ctx *gin.Context)  {
  	fmt.Println("1111")
  	
  	ctx.Next()
  
  	fmt.Println("4444")
  }
  // 创建 另外一种格式的中间件.
  func Test2() gin.HandlerFunc {
  	return func(context *gin.Context) {
  		fmt.Println("3333")
  
  		return
  		context.Next()
  
  		fmt.Println("5555")
  	}
  }
  func main()  {
  	router := gin.Default()
  
  	// 使用中间件
  	router.Use(Test1)
  	router.Use(Test2())
  
  	router.GET("/test", func(context *gin.Context) {
  		fmt.Println("2222")
  		context.Writer.WriteString("hello world!")
  	})
  
  	router.Run(":9999")
  }
  ```

  

##### Abort()

- 只执行当前中间件， 操作完成后，以出栈的顺序，依次返回上一级中间件。

  ```go
  // 创建中间件
  func Test1(ctx *gin.Context)  {
  	fmt.Println("1111")
  
  	ctx.Next()
  
  	fmt.Println("4444")
  }
  // 创建 另外一种格式的中间件.
  func Test2() gin.HandlerFunc {
  	return func(context *gin.Context) {
  		fmt.Println("3333")
  
  		context.Abort()
  
  		fmt.Println("5555")
  	}
  }
  func main()  {
  	router := gin.Default()
  
  	// 使用中间件
  	router.Use(Test1)
  	router.Use(Test2())
  
  	router.GET("/test", func(context *gin.Context) {
  		fmt.Println("2222")
  		context.Writer.WriteString("hello world!")
  	})
  
  	router.Run(":9999")
  }
  
  ```

  

<img src="图片/1582529575294.png" alt="1582529575294" style="zoom:80%;" />



<img src="图片/1582529769290.png" alt="1582529769290" style="zoom: 80%;" />

##### 中间件测试业务时间：

```go
// 创建中间件
func Test1(ctx *gin.Context)  {
	fmt.Println("1111")

	t := time.Now()

	ctx.Next()

	fmt.Println(time.Now().Sub(t))
}

// 创建 另外一种格式的中间件.
func Test2() gin.HandlerFunc {
	return func(context *gin.Context) {
		fmt.Println("3333")
        context.Abort()		// 将 Abort() 替换成 Next()， 反复测试，获取时间差！
		fmt.Println("5555")
	}
}
func main()  {
	router := gin.Default()

	// 使用中间件
	router.Use(Test1)
	router.Use(Test2())

	router.GET("/test", func(context *gin.Context) {
		fmt.Println("2222")
		context.Writer.WriteString("hello world!")
	})

	router.Run(":9999")
}
```



##### 小结

- 2种 书写格式：（见 前面笔记）
- 3个 操作函数/关键字： Next()、return、Abort()
- 作用域：作用域 以下 的路由。（ 对以上的 路由 无效！）



##### 项目中使用中间件

1. 在 web/main.go 中 创建 中间件。

   ```go
   func LoginFilter() gin.HandlerFunc {
   	return func(ctx *gin.Context) {
   		// 初始化 Session 对象
   		s := sessions.Default(ctx)
   		userName := s.Get("userName")
   
   		if userName == nil {
   			ctx.Abort()			// 从这里返回, 不必继续执行了
   		} else {
   			ctx.Next()			// 继续向下
   		}
   	}
   }
   ```

2. 在 所有需要进行 Session 校验操作之前， 添加、使用这个中间件。

   ```go
   // 添加路由分组
   r1 := router.Group("/api/v1.0")
   {
       r1.GET("/session", controller.GetSession)
       r1.GET("/imagecode/:uuid", controller.GetImageCd)
       r1.GET("/smscode/:phone", controller.GetSmscd)
       r1.POST("/users", controller.PostRet)
       r1.GET("/areas", controller.GetArea)
       r1.POST("/sessions", controller.PostLogin)
   
   r1.Use(LoginFilter())  //以后的路由,都不需要再校验 Session 了. 直接获取数据即可!
   
       r1.DELETE("/session", controller.DeleteSession)
       r1.GET("/user", controller.GetUserInfo)
       r1.PUT("/user/name", controller.PutUserInfo)
   }
   ```



#### 11 用户头像

1. 登录成功后，点击 “修改” 进入 “个人信息”页面。点击 “选择文件” ， 在系统选择一张图片作为头像。F12查看， 会在 Header 中， 看到 “avatar” 错误。

2. web/main.go 中 添加 路由回调。—— 参考 《接口文档.doc》

   ```go
   r1.GET("/user/avatar", controller.PostAvatar)
   ```

3. web/controller/user.go 中添加 GetUserInfo 函数， 并实现

   ```go
   // 上传头像
   func PostAvatar(ctx *gin.Context)  {
   }
   ```

4. 实现 PostAvatar 函数。 ———— 测试 gin 框架上传文件函数。

   ```go
   func PostAvatar(ctx *gin.Context) {
   	// 获取图片文件, 静态文件对象
   	file, _ := ctx.FormFile("avatar")
   	// 上传文件到项目中
   	err := ctx.SaveUploadedFile(file, "test/"+file.Filename)
   	fmt.Println(err)
   }
   ```

   可以将 头像，上传到 test/ 目录中。 



- 课后作业：
  - 看 《fastdfs.pdf》课件！ 配置 fastDFS、 Nginx 相应的环境。

#### 12  FastDFS 和 Nginx

##### FastDFS 和 Nginx

##### FastDFS

##### 三端：

- 客户端：client
- 监听端（监听服务器）：tracker
- 存储端（存储服务器）：storage

![1582680996644](图片/1582680996644.png)

##### 使用步骤：

1. 监听服务器定时查看存储服务器的状态。
2. client 访问监听服务器， 获取到可用的 存储服务器地址。
3. 客户端根据返回的地址，访问存储服务器。
4. 存储服务器存储文件，并返回凭证。如：“组名/M00/00/00/xxxxxx” 

##### 与以往的区别

- 图片需要上传！但是不需要下载。 直接按 “凭证” 展示到 浏览器即可。
- fastDFS 不提供 “展示图片” 功能。 —— Nginx。



##### fastDFS 安装和配置

- 安装 ： 参照 《fastdfs.pdf》实施

- 配置：

  - 修改sudo vim /etc/fdfs/storage.conf。	—— tracker_server=你自己的IP:22122
  - 修改sudo vim /etc/fdfs/tracker.conf        —— 不需要。
  - 修改sudo vim /etc/fdfs/client.conf          —— tracker_server=你自己的IP:22122

- 启动 fastDFS：

  1. 启动 存储服务器 storage：   sudo fdfs_storaged /etc/fdfs/storage.conf

  2. 启动 监听服务器 tracker： sudo fdfs_trackerd /etc/fdfs/tracker.conf

  3. ps aux | grep fdfs

     ![1582682178020](图片/1582682178020.png)



##### Nginx

- 安装： fastDFS 的 Nginx 插件。
  - 参照 《fastdfs.pdf》 3.2.5.7安装fastdfs-nginx-module

- 拷贝：
  - 将 解压缩的 fastdfs-master目录中 mod_fastdfs.conf 拷贝， 放到 /etc/fdfs/ 目录下。
  - 将 解压缩的 fastdfs-master目录中 http.conf  拷贝， 放到 /etc/fdfs/ 目录下。 --- 不需要修改
  - 将 解压缩的 fastdfs-master目录中 mime.types 拷贝， 放到 /etc/fdfs/ 目录下。 --- 不需要修改

- 修改 mod_fastdfs.conf
  - sudo  vim mod_fastdfs.conf 文件， 修改： tracker_server=你自己的IP地址:22122
- 修改 Nginx：
  - sudo vim /usr/local/nginx/conf/nginx.conf 
  - 参考 《fastdfs.pdf》 3.2.5.7 安装fastdfs-nginx-module 第 9 条。修改：

![1582683497645](图片/1582683497645.png)

- 启动nginx
  - 启动命令：`sudo /usr/local/nginx/sbin/nginx`
  - 提示：类似于：`ngx_http_fastdfs_set pid=78481`
  - 查看：
    - ps aux | grep nginx

##### Go语言 使用 FastDFS和Nginx

- 《fastdfs.pdf》中， 3.2.6 小结中的 方法， 已经过时！不要使用！

- 去 github 搜索 “fastdfs”  --- https://github.com/tedcy/fdfs_client

- 使用开源包，必须确认在  /etc/fdfs/client.conf 中 添加了 配置

  ```shell
  maxConns=10		# 设置最大连接数。
  ```

- 查看 demo：client_test.go ，确认，我们使用的方法：

  ```go
  client, err := NewClientWithConfig("fdfs.conf")
  	参数： /etc/fdfs/client.conf
  
  client.UploadByBuffer([]byte("hello world"), "go"); 
  	参1： []byte 的图片数据。
  	参2： 去除 “.” 文件后缀名。
  ```

  

##### 测试

```go
package main

import (
	"github.com/tedcy/fdfs_client"
	"fmt"
)

func main()  {
	// 初始化客户端 --- 配置文件
	clt, err := fdfs_client.NewClientWithConfig("/etc/fdfs/client.conf")
	if err != nil {
		fmt.Println("初始化客户端错误, err:", err)
		return
	}

	// 上传文件 -- 尝试文件名上传! 传入到 storage
	resp, err := clt.UploadByFilename("头像1.jpg")

	fmt.Println(resp, err)
}

-- 运行成功后 去 ~/fastdfs/storage/data/00/00/ 中查看！   ls | grep  xxx
```



##### 用户头像上传

1. 获取图片文件, 静态文件对象

   ```go
   file, _ := ctx.FormFile("avatar")
   ```

2. 上传头像到fastdfs 中, 按字节流。

   ```go
   clt, _ := fdfs_client.NewClientWithConfig("/etc/fdfs/client.conf")
   
   // 打开文件,读取文件内容
   f, _ := file.Open()			// 只读打开.
   
   buf := make([]byte, file.Size)	// 按文件实际大小,创建切片.
   
   f.Read(buf)		// 读取文件内容, 保存至buf缓冲区.
   
   // go语言根据文件名获取文件后缀
   fileExt := path.Ext(file.Filename)		// 传文件名, 获取文件后缀---- 带有"."
   
   // 按字节流上传图片内容
   remoteId, _ := clt.UploadByBuffer(buf, fileExt[1:])
   
   ```

3. 获取session , 得到当前用户

   ```go
   userName := sessions.Default(ctx).Get("userName")
   ```

4. web/model/modelFunc.go 中添加函数， 更新用户头像到数据中。 将 图片 凭证写入 avatar_url

   ```go
   func UpdateAvatar(userName, avatar string) error {
   	// update user set avatar_url = avatar, where name = username
   	return GlobalConn.Model(new(User)).Where("name = ?", userName).
   		Update("avatar_url", avatar).Error
   }
   ```

   

5. 根据用户名, 更新用户头像  --- MySQL数据库

   ```go
   model.UpdateAvatar(userName.(string), remoteId)
   ```

6. 这里 不在 remoteId 前，拼接 Nginx 使用的 IP:port。 这样写，会写死到 数据库。在 web/controller/user.go 中 修改 GetUserInfo() 函数， 在获取用户信息时， 添加  Nginx 使用的 IP:port。

   ```go
   func GetUserInfo(ctx *gin.Context) {
       ......
       temp["name"] = user.Name
   	temp["mobile"] = user.Mobile
   	temp["real_name"] = user.Real_name
   	temp["id_card"] = user.Id_card
       // 修改 avatar_url 写入的值。
   	temp["avatar_url"] = "http://192.168.6.108:8888/" + user.Avatar_url
   
   	resp["data"] = temp
   }
   ```

7. 成功，写出返回值给前端

   ```go
   resp := make(map[string]interface{})
   
   resp["errno"] = utils.RECODE_OK
   resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
   
   temp := make(map[string]interface{})
   temp["avatar_url"] = "http://192.168.6.108:8888/" + remoteId
   resp["data"] = temp
   
   ctx.JSON(http.StatusOK, resp)
   ```







