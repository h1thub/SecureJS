package cmd

import (
	"fmt"
	"log"
	"os"

	"SecureJS/config"
	"SecureJS/internal/analyze"
	"SecureJS/internal/crawler"
	"SecureJS/internal/matcher"
	"SecureJS/internal/output"
	"SecureJS/internal/parser"
	"SecureJS/internal/utils"

	"github.com/spf13/cobra"
)

var (
	singleURL  string
	listFile   string
	threads    int
	configPath string
	outputFile string
	browserPath string
	customHeaders []string
	proxy string
	ai string
	ARK_API_KEY string
	Model_ENDPOINT_ID string
)

func init() {
	rootCmd.Flags().StringVarP(&singleURL, "url", "u", "", "Single target URL to scan (e.g. https://example.com)")
	rootCmd.Flags().StringVarP(&listFile, "list", "l", "", "File containing target URLs (one per line)")
	rootCmd.Flags().IntVarP(&threads, "threads", "t", 20, "Number of concurrent threads for scanning")
	rootCmd.Flags().StringVarP(&configPath, "config", "c", "config/config.yaml", "Path to config file (e.g. config.yaml)")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file (supports .txt, .csv, .json)")
	rootCmd.Flags().StringVarP(&browserPath, "browser", "b", "", "Path to Chrome/Chromium executable (optional). If not set, will use Rod's default.")
	rootCmd.Flags().StringArrayVarP(&customHeaders, "header", "H", nil, "Add custom request headers. (e.g. -H 'Key: Value')")
	rootCmd.Flags().StringVarP(&proxy, "proxy", "p", "", "Proxy to use (e.g. http://127.0.0.1:8080)")
	rootCmd.Flags().StringVarP(&ai, "ai", "a", "false", "true/false. Enable AI Analytics. If not set, will use false")
	rootCmd.Flags().StringVarP(&Model_ENDPOINT_ID, "id", "i", "", "YOUR_ENDPOINT_ID")
	rootCmd.Flags().StringVarP(&ARK_API_KEY, "key", "k", "", "ARK_API_KEY")
}

var rootCmd = &cobra.Command{
	Use:   "SecureJS",
	Short: "A tool to crawl websites, parse links/JS, and match sensitive patterns based on custom config.",
	Long:  `...`,

	Run: func(cmd *cobra.Command, args []string) {
		// 1) 收集目标 URL
		var urls []string
		if singleURL != "" { // 使用 -u 参数
			urls = append(urls, singleURL)
		} else if listFile != "" { //使用 -l 参数
			var err error
			urls, err = utils.ReadURLs(listFile, urls)
			if err != nil {
				log.Fatalf("[!] Failed to read URLs from file %s: %v\n", listFile, err)
			}
		} else {
			fmt.Println("[!] Please provide either a single -u <URL> or a -l <file> containing URLs.")
			os.Exit(1)
		}

		// 2) 爬取链接（这里两个思路）2.1 && 2.2
		uniqueLinks := make(map[string]struct{}) // 使用 map 来跟踪已添加的链接，实现去重
		var toParse []string // 所有捕获的链接放入 toParse

		// 	2.1 收集加载某个目标 url 后（使用无头浏览器），默认加载的所有其他链接（js等）并放入 toParse
		err := crawler.CollectLinks(urls, threads, uniqueLinks, &toParse, browserPath, customHeaders, proxy)
		if err != nil {
			log.Fatalf("[!] Error collecting links: %v", err)
		}

		// 	2.2 收集加载某个目标 url 后其 body 中的链接（js等）并放入 toParse
		err = crawler.CollectLinksFromBody(urls, threads, uniqueLinks, &toParse, customHeaders, proxy)
		if err != nil {
			log.Fatalf("[!] Error collecting links from body: %v", err)
		}

		// for _, jsurl := range toParse {
		// 	fmt.Println(jsurl)
		// }

		// 3) 对所有收集到的链接进行二次请求
		parseResults, err := parser.ParseAll(toParse, threads, customHeaders, proxy)
		if err != nil {
			log.Fatalf("[!] Failed to parseAll: %v\n", err)
		}

		// 4) 加载 config.yaml 中的敏感信息正则匹配规则
		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			log.Fatalf("[!] Failed to load config: %v\n", err)
		}

		// 5) 对所有收集到的链接进行二次请求后的 body 中进行敏感信息的匹配
		matchResults, err := matcher.MatchAll(cfg.Rules, parseResults)
		if err != nil {
			log.Fatalf("[!] Failed to matchAll: %v\n", err)
		}

		// 6) 输出
		if outputFile == "" {
			if ai == "true" {
				resultString := analyze.FormatResultsToString(matchResults)
				analyze.Analyze(resultString, ARK_API_KEY, Model_ENDPOINT_ID)
			} else {
				output.PrintResultsToConsole(matchResults)
			}
		} else {
			err := output.WriteResultsToFile(matchResults, outputFile)
			if err != nil {
				log.Fatalf("[!] Failed to write results to file: %v\n", err)
			}
		}
	},
}

func Execute() error {
	return rootCmd.Execute()
}
