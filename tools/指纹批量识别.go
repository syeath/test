package tools

import (
	"fmt"
	"github.com/chainreactors/fingers"
	"net/http"
)

func Finger() {
	// 创建 Engine 实例
	engine, err := fingers.NewEngine()
	if err != nil {
		panic(err)
	}

	// 发起 HTTP 请求
	resp, err := http.Get("http://heepay-profession.com/login/index/login")
	if err != nil {
		fmt.Printf("HTTP 请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close() // 确保关闭响应体

	// 检测 favicon 指纹
	frame, err := engine.DetectResponse(resp)

	// 确保 frame 不是 nil
	if frame == nil {
		fmt.Println("未检测到有效的 favicon 指纹")
		return
	}

	// 打印指纹信息
	fmt.Println(frame.String())
}
