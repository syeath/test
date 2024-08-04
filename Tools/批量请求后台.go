package Tools

import (
	"flag"
	"fmt"
	"github.com/gocarina/gocsv"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"web/common"
	"web/utils"
)

// 分批大小
const batchSize = 1000

func RequestBackend() {
	// 定义命令行参数
	urlFlag := flag.String("f", "", "批量请求的URL文件")
	csvFlag := flag.String("o", "", "保存到CSV文件中")
	countFlag := flag.String("c", "10", "设置线程，默认10")
	proxyFlag := flag.String("p", "", "设置代理，格式为 127.0.0.1:7890")
	headerFlag := flag.String("header", "", "自定义请求头")
	adminFlag := flag.String("admin", "", "自定义后台地址")
	sleepFlag := flag.Int("s", 500, "每个线程请求间隔时间（毫秒），默认500毫秒")

	// 自定义帮助信息
	flag.Usage = func() {
		fmt.Printf("%s%s%s\n", common.Yellow, common.RequestBackendLogo, common.Reset)
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Println("Example usage:")
		fmt.Printf("%s  %s -f urls.txt -o output.csv %s\n", common.Blue, os.Args[0], common.Reset)
		fmt.Printf("%s  %s -f urls.txt -o output.csv -p 127.0.0.1:7890 %s\n", common.Blue, os.Args[0], common.Reset)
		fmt.Printf("%s  %s -f urls.txt -o output.csv -c 20 %s\n", common.Blue, os.Args[0], common.Reset)
		fmt.Printf("%s  %s -f urls.txt -o output.csv -c 20 -header xxx== %s\n", common.Blue, os.Args[0], common.Reset)
		fmt.Printf("%s  %s -f urls.txt -o output.csv -c 20 -header xxx== -admin admin %s\n", common.Blue, os.Args[0], common.Reset)
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

	// 读取 URL 文件
	data, err := os.ReadFile(*urlFlag)
	if err != nil {
		fmt.Println("文件读取失败:", err)
		return
	}

	// 分割文件内容为 URL 列表
	lines := strings.Split(string(data), "\n")
	var urls []string
	for _, line := range lines {
		if line = strings.TrimSpace(line); line != "" {
			urls = append(urls, utils.AddHTTPPrefix(line)...)
		}
	}

	// 创建 CSV 文件
	csvFile, err := os.OpenFile(*csvFlag, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("无法创建 CSV 文件:", err)
		return
	}
	defer csvFile.Close()

	// 创建 CSV 写入器
	writer := gocsv.DefaultCSVWriter(csvFile)
	defer writer.Flush()

	// 写入 CSV 头
	if err := writer.Write([]string{"URL", "Title", "Status"}); err != nil {
		fmt.Println("写入 CSV 头部失败:", err)
		return
	}

	// 处理 URL 批次
	for i := 0; i < len(urls); i += batchSize {
		end := i + batchSize
		if end > len(urls) {
			end = len(urls)
		}
		batch := urls[i:end]

		// 使用 WaitGroup 等待所有 goroutine 完成
		var wg sync.WaitGroup
		// URL 传递 channel
		urlChannel := make(chan string, len(batch))
		// 结果传递 channel
		resultChannel := make(chan []string, 100)

		// 启动多个 worker goroutine
		for j := 0; j < threadCount; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for url := range urlChannel {
					rsp, rspURL, err := utils.SendHttp(*adminFlag, *proxyFlag, *headerFlag, "GET", url)
					if err != nil {
						fmt.Printf("[-] %s%s 连接失败 %s\n", url, common.Red, common.Reset)
						resultChannel <- []string{url, "NetError", "-1"}
						continue
					}
					all, _ := io.ReadAll(rsp.Body)
					// 提取 title 内容
					titleRegex := regexp.MustCompile(`<title>(.*?)</title>`)
					matches := titleRegex.FindStringSubmatch(string(all))
					if len(matches) < 2 {
						fmt.Printf("[-] %s%s 找不到标题 %s\n", rspURL, common.Yellow, common.Reset)
						resultChannel <- []string{url, "No Title", strconv.Itoa(rsp.StatusCode)}
						continue
					}
					title := matches[1]
					fmt.Printf("[+] %s%s %s %s\n", rspURL, common.Green, title, common.Reset)
					resultChannel <- []string{rspURL, title, strconv.Itoa(rsp.StatusCode)}
					rsp.Body.Close()

					// 睡眠一段时间，单位为毫秒
					time.Sleep(time.Duration(*sleepFlag) * time.Millisecond)
				}
			}()
		}

		// 启动一个 goroutine 关闭结果 channel
		go func() {
			wg.Wait()
			close(resultChannel)
		}()

		// 将 URL 批次发送到 URL 通道
		for _, url := range batch {
			urlChannel <- url
		}
		close(urlChannel)

		// 写入结果到 CSV 文件
		for result := range resultChannel {
			if err := writer.Write(result); err != nil {
				fmt.Println("写入 CSV 行失败:", err)
			}
		}
	}
}
