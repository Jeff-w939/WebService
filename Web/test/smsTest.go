// 阿里云短信验证码测试服务
package main

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	dysmsapi "github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

func main() {
	config := sdk.NewConfig()

	// accessKeyid 通过本目录下AcessKey.csv获取  这文件也是阿里云申请的
	credential := credentials.NewAccessKeyCredential("LTAI5t6b1s556h1dQSaSfMys", "oAvHtLJQcaAQ5uRAY9YcPJSJ9C21bu")
	/* use STS Token
	credential := credentials.NewStsTokenCredential("<your-access-key-id>", "<your-access-key-secret>", "<your-sts-token>")
	*/
	client, err := dysmsapi.NewClientWithOptions("cn-hangzhou", config, credential)
	if err != nil {
		panic(err)
	}

	request := dysmsapi.CreateSendSmsRequest()

	request.Scheme = "https"
	request.PhoneNumbers = "18813150705"       // 发送手机号码
	request.SignName = "阿里云短信测试"               // 签名需要去阿里云申请
	request.TemplateCode = "SMS_154950909"     // 模板也需要申请
	request.TemplateParam = `{"code":1314520}` // 短信验证码
	response, err := client.SendSms(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Printf("response is %#v\n", response)
}
