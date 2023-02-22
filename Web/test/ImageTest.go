package main

import (
	"fmt"
	"github.com/afocus/captcha"
	"image/color"
	"image/png"
	"net/http"
)

func main() {
	// 1.初始化对象
	cap := captcha.New()

	// 2. 设置字体
	cap.SetFont("comic.ttf")

	//3. 设置验证码大小
	cap.SetSize(128, 64)

	//设置干扰强度
	cap.SetDisturbance(captcha.NORMAL)

	//设置前景色
	cap.SetFrontColor(color.RGBA{128, 65, 0, 128})

	//设置背景色
	cap.SetBkgColor(color.RGBA{255, 0, 0, 128}, color.RGBA{65, 126, 0, 255})

	//生成字体 将图片验证码展示到页面中
	http.HandleFunc("/r", func(w http.ResponseWriter, r *http.Request) { // /r 表示访问路径 127.0.0.1:8085/r
		img, str := cap.Create(4, captcha.ALL)
		png.Encode(w, img)
		println(str)
	})
	err := http.ListenAndServe(":8085", nil) // 设置的访问端口是8085
	if err != nil {
		fmt.Println("ListenAndServe err :", err)
		return
	}
}
