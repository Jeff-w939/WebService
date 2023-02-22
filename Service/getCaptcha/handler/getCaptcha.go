package handler

import (
	"context"
	"encoding/json"
	"getCaptcha/model"
	"github.com/afocus/captcha"
	"image/color"

	getCaptcha "getCaptcha/proto"
)

type GetCaptcha struct{}

// Return a new handler
func New() *GetCaptcha {
	return &GetCaptcha{}
}

// Call is a single request handler called via client.Call or the generated client code
func (e *GetCaptcha) Call(ctx context.Context, req *getCaptcha.Request, rsp *getCaptcha.Response) error {
	// 1.初始化对象
	cap := captcha.New()

	// 2. 设置字体
	cap.SetFont("./conf/comic.ttf")

	//3. 设置验证码大小
	cap.SetSize(128, 64)

	//设置干扰强度
	cap.SetDisturbance(captcha.NORMAL)

	//设置前景色
	cap.SetFrontColor(color.RGBA{128, 65, 0, 128})

	//设置背景色
	cap.SetBkgColor(color.RGBA{255, 0, 0, 128}, color.RGBA{65, 126, 0, 255})

	//生成字体
	img, str := cap.Create(4, captcha.NUM)

	//存储图片验证码到redis
	err := model.SaveImgcode(str, req.Uuid)
	if err != nil {
		return err
	}

	imgbuf, _ := json.Marshal(img)

	rsp.Msg = imgbuf
	return nil
}
