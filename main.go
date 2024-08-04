package main

import "web/Tools"

func main() {
	/*
		request		批量请求后台
		parse		从文件中提取域名或IP
		place		批量获取归属地
	*/

	// 批量请求后台
	//Tools.RequestBackend()
	// 文件提取IP或域名
	Tools.ParseDomain()
}
