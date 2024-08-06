package tools

import (
	"fmt"
	"regexp"
	"strings"
	"web/common"
	"web/utils"
)

func ParseDomainIP(inoutFileName, outPutFileName string, isSurvive bool) {

	fileCtx, err := utils.ReadText(inoutFileName)
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
		utils.WriteText(outPutFileName, correctedMatch)
	}
	// 是否探测存活
	if isSurvive {
		GetTitle(outPutFileName)
	}
}
