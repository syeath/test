package main

import "web/tools"

func main() {

	tools.Finger()

	//// 检查是否有子命令
	//if len(os.Args) < 2 {
	//	fmt.Println()
	//	fmt.Println("Usage: main.exe <command> [<args>]")
	//	fmt.Println()
	//	fmt.Println("  extract	提取域名和ip")
	//	fmt.Println("  finger	指纹识别")
	//	fmt.Println("  domain	子域名、端口收集")
	//	fmt.Println("  path		目录爆破")
	//	fmt.Println("  poc		poc识别")
	//	return
	//}
	//
	//// 根据子命令选择解析和执行  args []string
	//switch os.Args[1] {
	//case "extract":
	//	tools.ParseDomainIP(os.Args[2:])
	//case "finger":
	//	tools.GetTitle(os.Args[2:])
	//case "domain":
	//	domainCmd := flag.NewFlagSet("domain", flag.ExitOnError)
	//	domainInputUrl := domainCmd.String("u", "", "输入单个域名")
	//	domainInputFile := domainCmd.String("f", "", "输入文件")
	//	domainOutputFile := domainCmd.String("o", "", "输出文件")
	//	domainHelp := domainCmd.Bool("h", false, "显示帮助信息")
	//
	//	domainCmd.Parse(os.Args[2:])
	//
	//	if *domainHelp || *domainInputFile == "" || *domainOutputFile == "" || *domainInputUrl == "" {
	//		fmt.Println("Usage of domain:")
	//		domainCmd.PrintDefaults()
	//		return
	//	}
	//
	//	// 执行 domain 子命令的逻辑
	//
	//case "path":
	//	//tools.BlastDomain()
	//	//pathCmd := flag.NewFlagSet("path", flag.ExitOnError)
	//	//pathInputUrl := pathCmd.String("u", "", "输入单个域名")
	//	//pathInputFile := pathCmd.String("f", "", "输入文件")
	//	//pathOutputFile := pathCmd.String("o", "", "输出文件")
	//	//pathHelp := pathCmd.Bool("h", false, "显示帮助信息")
	//	//
	//	//pathCmd.Parse(os.Args[2:])
	//	//
	//	//if *pathHelp || *pathInputFile == "" || *pathOutputFile == "" && *pathInputUrl == "" {
	//	//	fmt.Println("Usage of domain:")
	//	//	pathCmd.PrintDefaults()
	//	//	return
	//	//}
	//	// 执行 domain 子命令的逻辑
	//
	//default:
	//	fmt.Println()
	//	fmt.Println("Usage: main.exe <command> [<args>]")
	//	fmt.Println()
	//	fmt.Println("  extract	提取域名和ip")
	//	fmt.Println("  finger	指纹识别")
	//	fmt.Println("  domain	子域名收集")
	//	fmt.Println("  ip		端口收集")
	//	fmt.Println("  path		目录爆破")
	//	fmt.Println("  poc		poc识别")
	//}
}
