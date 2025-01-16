package crawler

import (
	"SecureJS/internal/utils" // 你自己项目里的包路径，如有不同需改
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"

	// 确保版本更新到包含 `Inject` 的版本
	//"github.com/go-rod/stealth"
)

type CrawlResult struct {
	URL         string
	AllRequests []string
	Error       error
}

// -----------------------------------------------------------
// 并发爬取多个链接
// -----------------------------------------------------------
func crawlAll(urls []string, concurrency int) ([]*CrawlResult, error) {
	if len(urls) == 0 {
		return nil, fmt.Errorf("no URLs provided")
	}
	if concurrency <= 0 {
		concurrency = 1
	}

	chromePath := launcher.NewBrowser().MustGet()
	u := launcher.New().
		Bin(chromePath).
		Headless(true). // 调试时可设置为 false
		Set("ignore-certificate-errors").
		Set("disable-blink-features", "AutomationControlled").
		Set("disable-infobars").
		MustLaunch()

	browser := rod.New().ControlURL(u).MustConnect()
	defer func() {
		if err := browser.Close(); err != nil {
			log.Printf("[!] failed to close browser: %v\n", err)
		}
	}()

	resultChan := make(chan *CrawlResult, len(urls))
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	for _, targetURL := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			// 最大重试次数，可自行调整
			const maxRetry = 3
			res, err := fetchOneURLWithRetry(browser, url, maxRetry)
			if err != nil {
				resultChan <- &CrawlResult{URL: url, Error: err}
				return
			}
			resultChan <- res
		}(targetURL)
	}

	wg.Wait()
	close(resultChan)

	var results []*CrawlResult
	for r := range resultChan {
		results = append(results, r)
	}

	return results, nil
}

// -----------------------------------------------------------
// 带重试的抓取逻辑
// -----------------------------------------------------------
func fetchOneURLWithRetry(browser *rod.Browser, url string, maxAttempts int) (*CrawlResult, error) {
	var lastErr error
	baseTime := 20 * time.Second

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		currentTimeout := time.Duration(attempt) * baseTime

		result, err := tryFetchOneURL(browser, url, currentTimeout)
		if err == nil {
			return result, nil
		}

		lastErr = err
		log.Printf("[Attempt %d/%d] Failed to fetch '%s' (timeout=%v) error: %v",
			attempt, maxAttempts, url, currentTimeout, err)

		if attempt < maxAttempts {
			time.Sleep(2 * time.Second)
		}
	}

	return nil, lastErr
}

// -----------------------------------------------------------
// 单次访问逻辑：在已有 page 上使用 stealth.Inject(page)
// -----------------------------------------------------------
func tryFetchOneURL(browser *rod.Browser, url string, timeout time.Duration) (*CrawlResult, error) {
	page := browser.MustPage("")
	defer page.Close()

	// // 注入 stealth
	// if err := stealth.Inject(page); err != nil {
	// 	return nil, fmt.Errorf("failed to inject stealth: %w", err)
	// }

	// 设置 User-Agent
	err := page.SetUserAgent(&proto.NetworkSetUserAgentOverride{
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) " +
			"AppleWebKit/537.36 (KHTML, like Gecko) " +
			"Chrome/95.0.4638.69 Safari/537.36",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to set user agent: %w", err)
	}

	// 设置超时
	page = page.Timeout(timeout)

	loadedMap := make(map[string]bool)
	loadedMap[url] = true

	stop := page.EachEvent(func(e *proto.NetworkRequestWillBeSent) {
		reqURL := e.Request.URL
		lowerURL := strings.ToLower(reqURL)

		// 过滤
		if utils.Skip(lowerURL) {
			return
		}
		if utils.HasSkipExtension(lowerURL) {
			return
		}

		normalized := strings.TrimSuffix(reqURL, "/")
		loadedMap[normalized] = true
	})

	// 导航
	if err := page.Navigate(url); err != nil {
		stop()
		return nil, fmt.Errorf("failed to navigate %s: %w", url, err)
	}

	// 等待页面空闲
	if err := page.WaitIdle(15 * time.Second); err != nil {
		stop()
		return nil, fmt.Errorf("failed to wait idle %s: %w", url, err)
	}

	stop()

	allRequests := make([]string, 0, len(loadedMap))
	for r := range loadedMap {
		allRequests = append(allRequests, r)
	}

	return &CrawlResult{
		URL:         url,
		AllRequests: allRequests,
	}, nil
}

// -----------------------------------------------------------
// 对外的接口，用于收集
// -----------------------------------------------------------
func CollectLinks(urls []string, threads int, uniqueLinks map[string]struct{}, toParse *[]string) error {
	results, err := crawlAll(urls, threads)
	if err != nil {
		return fmt.Errorf("failed to crawl: %v", err)
	}

	for _, result := range results {
		if result.Error != nil {
			fmt.Printf("[!] URL: %s, Error: %v\n", result.URL, result.Error)
			continue
		}
		for _, reqURL := range result.AllRequests {
			if _, exists := uniqueLinks[reqURL]; !exists {
				uniqueLinks[reqURL] = struct{}{}
				*toParse = append(*toParse, reqURL)
			}
		}
	}
	return nil
}