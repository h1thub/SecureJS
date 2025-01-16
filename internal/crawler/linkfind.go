package crawler

import (
	"SecureJS/internal/parser"
	"SecureJS/internal/utils"
	"regexp"
)

func CollectLinksFromBody(urls []string, threads int, uniqueLinks map[string]struct{}, toParse *[]string) error {
    // 解析所有 URL 的内容
    parsedResult, err := parser.ParseAll(urls, threads)
    if err != nil {
        //return fmt.Errorf("解析失败: %v", err)
    }

    // 正则表达式匹配 http 或 https 链接
    urlRegex := regexp.MustCompile(`https?://[^\s"']+`)

    for _, parsed := range parsedResult {
        if parsed.Error != nil {
            //fmt.Printf("[!] 解析 URL: %s, 错误: %v\n", parsed.URL, parsed.Error)
            continue
        }

        // 查找所有匹配的链接
        foundURLs := urlRegex.FindAllString(parsed.Body, -1)

        for _, extractedURL := range foundURLs {
            // 调用 HasSkipExtension 函数判断是否过滤该链接
            if utils.HasSkipExtension(extractedURL) {
                continue
            }
            // 调用 Skip 函数判断是否过滤该链接
            if utils.Skip(extractedURL) {
                continue
            }

            if _, exists := uniqueLinks[extractedURL]; !exists {
                uniqueLinks[extractedURL] = struct{}{}
                *toParse = append(*toParse, extractedURL)
            }
        }
    }

    return nil
}