package parser

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// ParseResult 用于存储对单个 URL 做二次请求的结果
type ParseResult struct {
	URL        string // 当前目标 URL
	StatusCode int    // HTTP 状态码
	Body       string // 响应内容（纯文本/HTML/JSON 等）
	Error      error  // 如果请求失败或解析失败，则记录错误
}

// ParseAll 并发请求一批 URLs，并返回每个 URL 的响应内容。
// concurrency 用于控制并发线程数。
func ParseAll(urls []string, concurrency int, customHeaders []string, proxy string) ([]*ParseResult, error) {
	if len(urls) == 0 {
		return nil, fmt.Errorf("no URLs to parse")
	}
	if concurrency <= 0 {
		concurrency = 1
	}

    // 准备结果通道
    resultChan := make(chan *ParseResult, len(urls))

    // 并发控制 - 使用有缓冲的通道作为信号量
    sem := make(chan struct{}, concurrency)

    // 初始化 WaitGroup
    var wg sync.WaitGroup

    // 启动 goroutine 处理每个 URL
    for _, targetURL := range urls {
        wg.Add(1)

        go func(url string) {
            defer wg.Done()

            // 获取一个信号，限制并发数量
            sem <- struct{}{}
            defer func() { <-sem }()

            // 执行 URL 处理
            res, err := parseOneURL(url, customHeaders, proxy)
            if err != nil {
                resultChan <- &ParseResult{
                    URL:   url,
                    Error: err,
                }
                return
            }
            resultChan <- res
        }(targetURL)
    }

    // 等待所有 goroutine 完成
    wg.Wait()
    close(resultChan)

	// 收集结果
	var results []*ParseResult
	for r := range resultChan {
		results = append(results, r)
	}

	return results, nil
}

// parseOneURL 对单个 URL 发起请求，获取响应内容。
func parseOneURL(urlStr string, customHeaders []string, proxy string) (*ParseResult, error) {
	// 自定义 Transport，忽略证书错误
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // 忽略证书错误
		},
	}

	// 如果 proxy != "" 就设置代理
	if proxy != "" {
		proxyURL, err := url.Parse(proxy)
		if err == nil {
			tr.Proxy = http.ProxyURL(proxyURL)
		}
	}

	// 使用自定义 Transport
	client := &http.Client{
		Timeout:   10 * time.Second, 
		Transport: tr,
	}

	// 手动创建请求，以便设置 UA 和其他伪装头
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for %s: %w", urlStr, err)
	}

	// 设置伪装头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) "+
		"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36")

	// 自定义请求头
	for _, h := range customHeaders {
		parts := strings.SplitN(h, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			req.Header.Set(key, val)
		}
	}

	// 发起请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to GET %s: %w", urlStr, err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Printf("[!] failed to close response body for %s: %v", urlStr, cerr)
		}
	}()

	// 读取响应体
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body from %s: %w", urlStr, err)
	}

	body := strings.TrimSpace(string(bodyBytes))

	return &ParseResult{
		URL:        urlStr,
		StatusCode: resp.StatusCode,
		Body:       body,
	}, nil
}