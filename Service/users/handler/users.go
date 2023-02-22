package handler

import (
	"context"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"math/rand"
	"time"
	"users/model"
	"users/utils"

	users "users/proto"
)

type Users struct{}

// Return a new handler
func New() *Users {
	return &Users{}
}

// 发送短信验证码
func (e *Users) SendSms(ctx context.Context, req *users.Request, rsp *users.Response) error {
	// 校验 图片验证码
	result := model.CheckImgCode(req.Uuid, req.ImageCode)

	if result {
		// 发送短信
		config := sdk.NewConfig()

		// accessKeyid 通过本目录下AcessKey.csv获取  这文件也是阿里云申请的
		credential := credentials.NewAccessKeyCredential("LTAI5t6b1s556h1dQSaSfMys", "oAvHtLJQcaAQ5uRAY9YcPJSJ9C21bu")
		client, err := dysmsapi.NewClientWithOptions("cn-hangzhou", config, credential)
		//if err != nil {
		//	panic(err)
		//}

		request := dysmsapi.CreateSendSmsRequest()

		request.Scheme = "https"
		request.PhoneNumbers = req.Phone       // 发送手机号码
		request.SignName = "阿里云短信测试"           // 签名需要去阿里云申请
		request.TemplateCode = "SMS_154950909" // 模板也需要申请

		//生成一个随机6位数做验证码
		rand.Seed(time.Now().UnixNano())
		smscode := fmt.Sprintf("%06d", rand.Int31n(1000000))
		request.TemplateParam = `{"code": ` + smscode + `}` // 短信验证码

		response, err := client.SendSms(request)

		if response.IsSuccess() {
			// 发送短信验证码成功
			rsp.Errno = utils.RECODE_OK
			rsp.Errmsg = utils.RecodeText(utils.RECODE_OK)

			// 储存验证码到redis
			err := model.SaveSmsCode(req.Phone, smscode)
			if err != nil {
				fmt.Println("存储到redis 失败：", err)
			}
		} else {
			rsp.Errno = utils.RECODE_DATAERR
			rsp.Errmsg = utils.RecodeText(utils.RECODE_DATAERR)
		}
		if err != nil {
			rsp.Errno = utils.RECODE_DATAERR
			rsp.Errmsg = utils.RecodeText(utils.RECODE_DATAERR)
		}

	}
	return nil
}

// 注册用户到Mysql
func (e *Users) Register(ctx context.Context, req *users.RegReq, rsp *users.Response) error {
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
	} else { // 短信验证码错误
		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DATAERR)
	}
	return nil
}
