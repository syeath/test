package tools

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func RunCmd() {

	//	fmt.Println("  extract	提取域名和ip探测存活")
	//	fmt.Println("  finger	指纹识别")
	//	fmt.Println("  domain	子域名收集")
	//	fmt.Println("  port		端口收集")
	//	fmt.Println("  path		目录爆破")
	//	fmt.Println("  poc		poc识别")

	var inputFile string
	var outputFile string
	var isSurvive bool
	var proxy string
	var header string
	var userAgent bool
	var thread int

	var extractCmd = &cobra.Command{
		Use:   "extract",
		Short: "提取文件中的url",
		Long:  `提取文件中的url`,
		Run: func(cmd *cobra.Command, args []string) {
			// 检查是否提供了必需的参数
			if inputFile == "" || outputFile == "" {
				fmt.Println("请提供输入文件和输出文件")
				cmd.Help()
				os.Exit(1)
			}

			// 提取逻辑
			ParseDomainIP(inputFile, outputFile, isSurvive)
		},
	}
	extractExample := `  %s extract -f urls.txt -o result.txt
  %s extract -f urls.txt -o result.txt --is-survive`
	extractCmd.Example = fmt.Sprintf(extractExample, os.Args[0], os.Args[0])

	extractCmd.Flags().StringVarP(&inputFile, "input-file", "f", "", "要提取的url文件")
	extractCmd.Flags().StringVarP(&outputFile, "output-file", "o", "", "保存的文件名 [result.txt]")
	extractCmd.Flags().BoolVar(&isSurvive, "is-survive", false, "探测url是否存活")
	extractCmd.Flags().BoolP("help", "h", false, "显示帮助信息")

	var fingerCmd = &cobra.Command{
		Use:   "finger",
		Short: "批量指纹识别",
		Long:  "批量指纹识别",
		Run: func(cmd *cobra.Command, args []string) {
			// 检查是否提供了必需的参数
			if inputFile == "" {
				fmt.Println("请提供输入文件")
				cmd.Help()
				os.Exit(1)
			}
			// 执行指纹识别逻
			Finger(inputFile, proxy, header, userAgent, thread)
		},
	}
	fingerCmd.Example = "  # 从url文件中批量识别指纹\n  ./main finger -f urls.txt\n\n" +
		"  # 使用代理和自定义请求头\n  ./main finger -f urls.txt --proxy http://proxy:8080 --header 'X-My-Header: Value'"
	fingerExample := `  %s finger -f urls.txt
  %s finger -f urls.txt --proxy 127.0.0.1:7890 --header 'X-My-Header: Value;Referer: 127.0.0.1'
  %s finger -f urls.txt --proxy 127.0.0.1:7890 --ua'`
	fingerCmd.Example = fmt.Sprintf(fingerExample, os.Args[0], os.Args[0], os.Args[0])

	fingerCmd.Flags().StringVarP(&inputFile, "input-file", "f", "", "url文件")
	fingerCmd.Flags().IntVarP(&thread, "thread", "t", 10, "设置线程 [默认10]")
	fingerCmd.Flags().StringVar(&proxy, "proxy", "", "设置代理")
	fingerCmd.Flags().StringVar(&header, "header", "", "设置请求头")
	fingerCmd.Flags().BoolVar(&userAgent, "ua", false, "随机的ua请求头")
	fingerCmd.Flags().BoolP("help", "h", false, "显示帮助信息")

	var rootCmd = &cobra.Command{
		Use: os.Args[0],
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}

	// 定义自定义的帮助命令
	var rootHelpCmd = &cobra.Command{
		Use:   "help",
		Short: "显示帮助信息",
		Long:  "显示该工具的帮助信息。",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				// 显示根命令的帮助信息
				rootCmd.Help()
			} else {
				// 查找指定的子命令并显示其帮助信息
				cmd, _, err := rootCmd.Find(args)
				if err != nil {
					fmt.Printf("未找到命令: %s\n", args[0])
					return
				}
				cmd.Help()
			}
		},
	}

	// 设置自定义帮助命令
	rootCmd.SetHelpCommand(rootHelpCmd)

	rootCmd.AddCommand(extractCmd, fingerCmd)
	rootCmd.Execute()
}
