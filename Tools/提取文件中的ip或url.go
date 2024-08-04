package Tools

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"web/common"
	"web/utils"
)

/*
	1. 提取所有的ip和域名
	2. 使用httpx探测
	3.
*/

func ParseDomain() {
	// 定义命令行参数
	urlsFlag := flag.String("f", "", "批量请求的URL文件")
	outputFlag := flag.String("o", "", "保存到txt文件中")

	// 自定义帮助信息
	flag.Usage = func() {
		fmt.Printf("%s%s%s\n", common.Yellow, common.RequestBackendLogo, common.Reset)
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Println("Example usage:")
		fmt.Printf("%s  %s -f urls.txt -o output.txt %s\n", common.Blue, os.Args[0], common.Reset)
	}

	// 解析命令行参数
	flag.Parse()

	// 如果没有提供任何参数或者提供了帮助参数，则显示帮助信息
	if *urlsFlag == "" || *outputFlag == "" {
		flag.Usage()
		return
	}

	// 获取保存文件的绝对路径
	executablePath, _ := os.Executable()
	outputPath := fmt.Sprintf("%s\\%s", filepath.Dir(executablePath), *outputFlag)

	fileCtx, err := utils.ReadText(*urlsFlag)
	if err != nil {
		fmt.Println(err)
		return
	}

	lines := strings.Split(fileCtx, "\n")

	// 定义正则
	urlRegex := regexp.MustCompile(common.URLPattern)
	ipRegex := regexp.MustCompile(common.IPPattern)

	var results []string

	for _, line := range lines {
		// 提取url
		urls := urlRegex.FindAllString(line, -1)
		for _, url := range urls {
			results = append(results, url)
		}
		// 提取 IP 地址
		ipMatches := ipRegex.FindAllString(line, -1)
		for _, ip := range ipMatches {
			results = append(results, ip)
		}
	}

	for _, v := range results {
		utils.WriteText(*outputFlag, v)
	}

	// 执行命令
	args, ok := common.HttpxConfig["url"].([]string)
	if !ok {
		log.Fatalf("无法获取 httpx 参数")
	}

	fmt.Println(executablePath)

	args = append(args, outputPath)
	// 创建 exec.Command 对象
	cmd := exec.Command("httpx", args...)
	// 创建 Stdout 和 Stderr 管道
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("创建 Stdout 管道失败: %v", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		log.Fatalf("创建 Stderr 管道失败: %v", err)
	}

	// 启动命令
	if err := cmd.Start(); err != nil {
		log.Fatalf("启动命令失败: %v", err)
	}

	// 实时读取 Stdout
	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			line := scanner.Text()
			// 过滤输出行，确保只输出以 http:// 或 https:// 开头的行
			if strings.HasPrefix(line, "http://") || strings.HasPrefix(line, "https://") {
				fmt.Println(line)
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "读取 Stdout 时出错: %v\n", err)
		}
	}()

	// 实时读取 Stderr
	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			line := scanner.Text()
			// 过滤输出行，确保只输出以 http:// 或 https:// 开头的行
			if strings.HasPrefix(line, "http://") || strings.HasPrefix(line, "https://") {
				fmt.Println(line)
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "读取 Stderr 时出错: %v\n", err)
		}
	}()

	// 等待命令执行完成
	if err := cmd.Wait(); err != nil {
		log.Fatalf("执行命令时发生错误: %v", err)
	}
}
