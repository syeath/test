package httpx

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"context"
	"github.com/chainreactors/fingers"
)

// GetTitle httpx匹配出标题
func GetTitle(args []string) {
	fingerCmd := flag.NewFlagSet("finger", flag.ExitOnError)
	// 定义命令行参数
	fingerInputFile := fingerCmd.String("f", "", "批量请求的url文件")
	fingerHelp := fingerCmd.Bool("h", false, "显示帮助信息")

	fingerCmd.Parse(args)

	if *fingerHelp || *fingerInputFile == "" {
		fmt.Println("Usage of finger:")
		fingerCmd.PrintDefaults()
		return
	}

	// 获取保存文件的绝对路径
	executablePath, _ := os.Executable()
	outputPath := fmt.Sprintf("%s\\%s", filepath.Dir(executablePath), *fingerInputFile)

	// 创建一个新的指纹检测器
	fingerprints, err := fingers.NewFingerprints()
	if err != nil {
		log.Fatalf("无法创建指纹检测器: %v", err)
	}

	fileCtx, err := file.ReadText(outputPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	domains := strings.Split(fileCtx, "\n")

	// 使用上下文进行指纹识别
	ctx := context.Background()

	// 进行域名指纹识别
	for _, domain := range domains {
		fmt.Printf("检测域名: %s\n", domain)

		// 获取所有指纹库的指纹信息
		allFingerprints, err := fingerprints.All(ctx, domain)
		if err != nil {
			log.Printf("域名 %s 指纹识别失败: %v", domain, err)
			continue
		}

		// 打印指纹信息
		fmt.Printf("域名 %s 的指纹信息:\n", domain)
		for _, fp := range allFingerprints {
			fmt.Printf(" - %s: %s\n", fp.Key, fp.Value)
		}
		fmt.Println()
	}


	//// 执行命令
	//args, ok := common.HttpxConfig["url"].([]string)
	//if !ok {
	//	log.Fatalf("无法获取 httpx 参数")
	//}
	//
	//args = append(args, "-l")
	//args = append(args, outputPath)
	//// 创建 exec.Command 对象
	//cmd := exec.Command("httpx", args...)
	//// 创建 Stdout 和 Stderr 管道
	//stdoutPipe, err := cmd.StdoutPipe()
	//if err != nil {
	//	log.Fatalf("创建 Stdout 管道失败: %v", err)
	//}
	//stderrPipe, err := cmd.StderrPipe()
	//if err != nil {
	//	log.Fatalf("创建 Stderr 管道失败: %v", err)
	//}
	//
	//// 启动命令
	//if err := cmd.Start(); err != nil {
	//	log.Fatalf("启动命令失败: %v", err)
	//}
	//
	//// 实时读取 Stdout
	//go func() {
	//	scanner := bufio.NewScanner(stdoutPipe)
	//	for scanner.Scan() {
	//		line := scanner.Text()
	//		// 过滤输出行，确保只输出以 http:// 或 https:// 开头的行
	//		if strings.HasPrefix(line, "http://") || strings.HasPrefix(line, "https://") {
	//			fmt.Println(line)
	//		}
	//	}
	//	if err := scanner.Err(); err != nil {
	//		fmt.Fprintf(os.Stderr, "读取 Stdout 时出错: %v\n", err)
	//	}
	//}()
	//
	//// 实时读取 Stderr
	//go func() {
	//	scanner := bufio.NewScanner(stderrPipe)
	//	for scanner.Scan() {
	//		line := scanner.Text()
	//		// 过滤输出行，确保只输出以 http:// 或 https:// 开头的行
	//		if strings.HasPrefix(line, "http://") || strings.HasPrefix(line, "https://") {
	//			fmt.Println(line)
	//		}
	//	}
	//	if err := scanner.Err(); err != nil {
	//		fmt.Fprintf(os.Stderr, "读取 Stderr 时出错: %v\n", err)
	//	}
	//}()
	//
	//// 等待命令执行完成
	//if err := cmd.Wait(); err != nil {
	//	log.Fatalf("执行命令时发生错误: %v", err)
	//}
}
