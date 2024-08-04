package Tools

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/gocarina/gocsv"
	"io"
	"os"
	"regexp"
	"strconv"
	"sync"
	"web/utils"
)

var (
	red    = "\033[0;31m"
	green  = "\033[0;32m"
	yellow = "\033[0;33m"
	blue   = "\033[0;36m"
	reset  = "\033[0m"
	logo   = `
╦═╗┌─┐┌─┐ ┬ ┬┌─┐┌─┐┌┬┐  ┌┐ ┌─┐┌─┐┬┌─┌─┐┌┐┌┌┬┐
╠╦╝├┤ │─┼┐│ │├┤ └─┐ │   ├┴┐├─┤│  ├┴┐├┤ │││ ││
╩╚═└─┘└─┘└└─┘└─┘└─┘ ┴   └─┘┴ ┴└─┘┴ ┴└─┘┘└┘─┴┘
			批量请求后台 by liangc`
)

func RequestBackend() {
	/*
		批量请求后台
		-f 批量请求的url
		-o 保存到csv中
		-p 设置代理 127.0.0.1:7890 支持socket5代理
		-c 设置线程
		--header 设置请求头
		--admin	自定义请求后台
	*/

	// 定义命令行参数
	urlFlag := flag.String("f", "", "批量请求的URL文件")
	csvFlag := flag.String("o", "", "保存到CSV文件中")
	countFlag := flag.String("c", "10", "设置线程，默认10")
	proxyFlag := flag.String("p", "", "设置代理，格式为 127.0.0.1:7890")
	headerFlag := flag.String("header", "", "自定义请求头")
	adminFlag := flag.String("admin", "", "自定义后台地址")

	// 自定义帮助信息
	flag.Usage = func() {
		fmt.Printf("%s%s%s\n", yellow, logo, reset)
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Println("Example usage:")
		fmt.Printf("%s  %s -f urls.txt -o output.csv %s\n", blue, os.Args[0], reset)
		fmt.Printf("%s  %s -f urls.txt -o output.csv -p 127.0.0.1:7890 %s\n", blue, os.Args[0], reset)
		fmt.Printf("%s  %s -f urls.txt -o output.csv -c 20 %s\n", blue, os.Args[0], reset)
		fmt.Printf("%s  %s -f urls.txt -o output.csv -c 20 -header xxx== %s\n", blue, os.Args[0], reset)
		fmt.Printf("%s  %s -f urls.txt -o output.csv -c 20 -header xxx== -admin admin %s\n", blue, os.Args[0], reset)
	}

	// 解析命令行参数
	flag.Parse()

	// 如果没有提供任何参数或者提供了帮助参数，则显示帮助信息
	if *urlFlag == "" || *csvFlag == "" {
		flag.Usage()
		return
	}
	// 设置线程池大小
	threadCount := 10
	fmt.Sscanf(*countFlag, "%d", &threadCount)
	if threadCount <= 0 {
		threadCount = 10
	}

	file, err := os.Open(*urlFlag)
	if err != nil {
		fmt.Println("文件不存在!!!")
		return
	}
	defer file.Close()

	csvFile, _ := os.OpenFile(*csvFlag, os.O_CREATE|os.O_RDWR, 0666)
	defer csvFile.Close()

	// 创建CSV写入器
	writer := gocsv.DefaultCSVWriter(csvFile)
	defer writer.Flush()

	// 写入CSV头
	if err := writer.Write([]string{"URL", "Title", "Status"}); err != nil {
		fmt.Println("写入CSV头部失败!!!")
		return
	}

	// 使用 WaitGroup 等待所有 goroutine 完成
	var wg sync.WaitGroup
	// URL 传递 channel
	urlChannel := make(chan string)
	// 结果传递 channel
	resultChannel := make(chan []string)

	// 启动多个 worker goroutine
	for i := 0; i < threadCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range urlChannel {
				rsp, rspURL, err := utils.SendHttp(*adminFlag, *proxyFlag, *headerFlag, "GET", url)
				if err != nil {
					fmt.Printf("[-] %s%s 连接失败 %s\n", url, red, reset)
					resultChannel <- []string{url, "NetError", "-1"}
					continue
				}
				all, _ := io.ReadAll(rsp.Body)
				// 提取 title 内容
				titleRegex := regexp.MustCompile(`<title>(.*?)</title>`)
				matches := titleRegex.FindStringSubmatch(string(all))
				if len(matches) < 2 {
					fmt.Printf("[-] %s%s 找不到标题 %s\n", rspURL, yellow, reset)
					resultChannel <- []string{url, "No Title", strconv.Itoa(rsp.StatusCode)}
					continue
				}
				title := matches[1]
				fmt.Printf("[+] %s%s %s %s\n", rspURL, green, title, reset)
				resultChannel <- []string{rspURL, title, strconv.Itoa(rsp.StatusCode)}
			}
		}()
	}

	// 启动一个 goroutine 关闭结果 channel
	go func() {
		wg.Wait()
		close(resultChannel)
	}()

	// 逐行读取 URL 判断有无 http:// 或 https:// 并发送到 channel
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		urls := utils.AddHTTPPrefix(scanner.Text())
		for _, url := range urls {
			urlChannel <- url
		}
	}
	close(urlChannel)

	// 写入结果到 CSV 文件
	for result := range resultChannel {
		if err := writer.Write(result); err != nil {
			fmt.Println("写入CSV行失败!!!")
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("读取文件错误:", err)
	}
}
