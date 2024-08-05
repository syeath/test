package httpx

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"web/common"
)

// GetTitle httpx匹配出标题
func GetTitle(args []string) {
	fileName := "123"
	// 获取保存文件的绝对路径
	executablePath, _ := os.Executable()
	outputPath := fmt.Sprintf("%s\\%s", filepath.Dir(executablePath), fileName)

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
