package tools

import (
	"fmt"
	"github.com/chainreactors/fingers"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"web/common"
	"web/utils"
)

// Finger 进行指纹识别，支持多线程
func Finger(inputFileName, proxy, header string, ua bool, numThreads int) {
	// 获取文件的绝对路径
	executablePath, _ := os.Executable()
	outputPath := filepath.Join(filepath.Dir(executablePath), inputFileName)

	fileCtx, err := utils.ReadText(outputPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	lines := strings.Split(strings.TrimSpace(fileCtx), "\n")

	// 创建 Engine 实例
	engine, err := fingers.NewEngine()
	if err != nil {
		panic(err)
	}
	if ua {
		header = header + ",User-Agent: " + common.RandomUserAgent()
	}

	// 创建 WaitGroup 和信号通道
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, numThreads)

	for _, line := range lines {
		wg.Add(1)
		semaphore <- struct{}{} // 添加信号，确保并发控制

		go func(url string) {
			defer wg.Done()
			defer func() { <-semaphore }() // 释放信号

			resp, err := utils.SendHttp(proxy, header, "GET", url)
			if err != nil {
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 && resp.StatusCode != 300 && resp.StatusCode != 301 && resp.StatusCode != 302 && resp.StatusCode != 400 && resp.StatusCode != 404 {
				return
			}

			frame, err := engine.DetectResponse(resp)
			// 确保 frame 不是 nil
			if frame == nil {
				fmt.Println("未检测到有效的 favicon 指纹")
				return
			}

			// 收集指纹信息
			fingerprints := collectFingerprints(frame.String())
			uniqueFingerprints := unique(fingerprints)

			formattedFingerprints := formatFingerprints(uniqueFingerprints)
			fmt.Printf("%s[%s]%s %s %s%s%s\n", common.Blue, common.GetCurrentTime(), common.Reset, url, common.Green, formattedFingerprints, common.Reset)
		}(line)
	}

	// 等待所有 goroutines 完成
	wg.Wait()
}

func collectFingerprints(input string) []string {
	var fingerprints []string

	// 按行分割输入字符串
	lines := strings.Split(input, "\n")

	for _, line := range lines {
		// 按 "||" 分割每一行
		parts := strings.Split(line, "||")
		for _, part := range parts {
			// 去掉左右空格，并检查是否含有 ":"
			part = strings.TrimSpace(part)
			if strings.Contains(part, ":") {
				// 分割 ":"
				splitPart := strings.Split(part, ":")
				if len(splitPart) == 2 {
					// 获取名称
					name := strings.TrimSpace(splitPart[0])
					fingerprints = append(fingerprints, name)
				} else {
					// 如果没有来源，直接添加名称
					fingerprints = append(fingerprints, strings.TrimSpace(part))
				}
			} else {
				// 如果没有 ":"
				fingerprints = append(fingerprints, part)
			}
		}
	}

	return fingerprints
}

// unique 去重函数
func unique(input []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	for _, entry := range input {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// formatFingerprints 将每个指纹信息用方括号括起来
func formatFingerprints(fingerprints []string) string {
	for i, fp := range fingerprints {
		fingerprints[i] = "[" + fp + "]"
	}
	return strings.Join(fingerprints, " ")
}
