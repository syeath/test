package main

import (
	"flag"
	"fmt"
	"os"
	"web/utils/httpx"
	"web/utils/tools"
)

func main() {

	// 检查是否有子命令
	if len(os.Args) < 2 {
		fmt.Println()
		fmt.Println("Usage: main.exe <command> [<args>]")
		fmt.Println("Available commands are:")
		fmt.Println("  extract	提取域名和ip")
		fmt.Println("  domain	子域名字典爆破")
		fmt.Printf("  path		路径爆破\n\n")
		return
	}

	// 根据子命令选择解析和执行  args []string
	switch os.Args[1] {
	case "extract":
		tools.ParseDomainIP(os.Args[2:])
	case "finger":
		httpx.GetTitle(os.Args[2:])
	case "domain":
		domainCmd := flag.NewFlagSet("domain", flag.ExitOnError)
		domainInputUrl := domainCmd.String("u", "", "输入单个域名")
		domainInputFile := domainCmd.String("f", "", "输入文件")
		domainOutputFile := domainCmd.String("o", "", "输出文件")
		domainHelp := domainCmd.Bool("h", false, "显示帮助信息")

		domainCmd.Parse(os.Args[2:])

		if *domainHelp || *domainInputFile == "" || *domainOutputFile == "" || *domainInputUrl == "" {
			fmt.Println("Usage of domain:")
			domainCmd.PrintDefaults()
			return
		}

		// 执行 domain 子命令的逻辑

	case "path":
		httpx.BlastDomain()
		pathCmd := flag.NewFlagSet("path", flag.ExitOnError)
		pathInputUrl := pathCmd.String("u", "", "输入单个域名")
		pathInputFile := pathCmd.String("f", "", "输入文件")
		pathOutputFile := pathCmd.String("o", "", "输出文件")
		pathHelp := pathCmd.Bool("h", false, "显示帮助信息")

		pathCmd.Parse(os.Args[2:])

		if *pathHelp || *pathInputFile == "" || *pathOutputFile == "" && *pathInputUrl == "" {
			fmt.Println("Usage of domain:")
			pathCmd.PrintDefaults()
			return
		}
		// 执行 domain 子命令的逻辑

	default:
		fmt.Println()
		fmt.Println("Usage: main.exe <command> [<args>]")
		fmt.Println("Available commands are:")
		fmt.Println("  extract	提取域名和ip")
		fmt.Println("  finger	指纹识别")
		fmt.Println("  domain	子域名字典爆破")
		fmt.Printf("  path		路径爆破\n\n")
	}

	//// 提取域名和ip
	//extractCmd := flag.NewFlagSet("extract", flag.ExitOnError)
	//// 域名爆破
	//domainCmd := flag.NewFlagSet("domain", flag.ExitOnError)
	//// 资产探活识别
	//fingerCmd := flag.NewFlagSet("finger", flag.ExitOnError)
	//
	//// domain 命令行参数
	//domainInputFile := domainCmd.String("f", "", "输入文件")
	//domainOutputFile := domainCmd.String("o", "", "输出文件")
	//
	//// gettitle 命令行参数
	//getTitleOutputFile := getTitleCmd.String("o", "", "输出文件")
	//
	//if len(os.Args) < 2 {
	//	fmt.Println("expected 'domain' or 'gettitle' subcommands")
	//	os.Exit(1)
	//}

	/*
		request		批量请求后台
		parse		从文件中提取域名或IP
		place		批量获取归属地
	*/

	// 批量请求后台
	//Tools.RequestBackend()
	// 文件提取IP或域名
	//Tools.ParseDomainIP()
}
