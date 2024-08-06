package utils

import (
	"context"
	"crypto/tls"
	"fmt"
	"golang.org/x/net/proxy"
	"net"
	"net/http"
	"strings"
	"time"
)

// SendHttp 发起请求
func SendHttp(proxySetting, headerSetting, method, url string) (*http.Response, error) {
	transport := &http.Transport{}
	// 忽略https证书问题
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	if proxySetting != "" {
		// 创建一个不验证证书的 HTTP 客户端
		dialer, err := proxy.SOCKS5("tcp", proxySetting, nil, proxy.Direct)
		if err != nil {
			return nil, fmt.Errorf("代理连接失败")
		}
		dialContextFunc := func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		}
		// 配置代理
		transport.DialContext = dialContextFunc
	}
	client := &http.Client{
		Timeout:   5 * time.Second,
		Transport: transport,
	}

	req, err := CreateClient(method, url, headerSetting)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("后台请求失败")
	}
	return resp, nil
}

// CreateClient 创建一个http
func CreateClient(method, url, headersStr string) (*http.Request, error) {
	// url判断有无http或者https
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(headersStr, "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) == 2 {
			request.Header.Add(parts[0], parts[1])
		}
	}

	return request, nil
}

// AddHTTPPrefix 添加http://
func AddHTTPPrefix(ipOrDomain string) []string {
	var urls []string
	if !strings.HasPrefix(ipOrDomain, "http://") && !strings.HasPrefix(ipOrDomain, "https://") {
		urls = append(urls, "http://"+ipOrDomain)
		urls = append(urls, "https://"+ipOrDomain)
		return urls
	}
	return []string{ipOrDomain}
}
