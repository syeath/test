package common

const (
	Red        = "\033[0;31m"
	Green      = "\033[0;32m"
	Yellow     = "\033[0;33m"
	Blue       = "\033[0;36m"
	Reset      = "\033[0m"
	URLPattern = `https?:\/\/[^\s,"]+`
	IPPattern  = `((?:\d{1,3}\.){3}\d{1,3})(?::\d{1,5})?`

	RequestBackendLogo = `
╦═╗┌─┐┌─┐ ┬ ┬┌─┐┌─┐┌┬┐  ┌┐ ┌─┐┌─┐┬┌─┌─┐┌┐┌┌┬┐
╠╦╝├┤ │─┼┐│ │├┤ └─┐ │   ├┴┐├─┤│  ├┴┐├┤ │││ ││
╩╚═└─┘└─┘└└─┘└─┘└─┘ ┴   └─┘┴ ┴└─┘┴ ┴└─┘┘└┘─┴┘
			批量请求后台 by liangc`
)

// 设置要执行的命令和参数

var (
	HttpxConfig = map[string]interface{}{
		// 探测url以及指纹和ip
		"url": []string{
			"-title",
			"-status-code",
			"-tech-detect",
			"-l",
		},
		// 爆破子域名
		"domain": []string{
			"-title",
			"-status-code",
			"-tech-detect",
			"-l",
		},
		// c段探测
		"c": []string{},
		// 指定端口探测
		"port": []string{},
		// 爆破指定目录
		"path": []string{},
	}
)
