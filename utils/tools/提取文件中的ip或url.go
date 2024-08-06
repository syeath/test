package tools

import (
	"flag"
	"fmt"
	"regexp"
	"strings"
	"web/common"
	"web/utils/file"
)

func ParseDomainIP(args []string) {

	extractCmd := flag.NewFlagSet("domain", flag.ExitOnError)
	// 定义命令行参数
	extractInputFile := extractCmd.String("f", "", "批量请求的url文件")
	extractOutputFile := extractCmd.String("o", "", "保存到txt文件中")
	extractHelp := extractCmd.Bool("h", false, "显示帮助信息")

	extractCmd.Parse(args)

	if *extractHelp || *extractInputFile == "" {
		fmt.Println("Usage of domain:")
		extractCmd.PrintDefaults()
		return
	}

	fileCtx, err := file.ReadText(*extractInputFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	lines := strings.Split(fileCtx, "\n")

	// 定义正则
	urlRegex := regexp.MustCompile(common.URLPattern)

	var results []string

	for _, line := range lines {
		// 提取 url
		urls := urlRegex.FindAllString(line, -1)
		for _, url := range urls {
			results = append(results, url)
		}
	}

	for _, v := range results {
		correctedMatch := strings.ReplaceAll(v, "：", ":")
		// 检查是否已包含 http 或 https 协议
		if !strings.HasPrefix(correctedMatch, "http://") && !strings.HasPrefix(correctedMatch, "https://") {
			// 添加 https 协议
			correctedMatch = "https://" + correctedMatch
		}
		file.WriteText(*extractOutputFile, correctedMatch)
	}
}
