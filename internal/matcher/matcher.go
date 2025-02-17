package matcher

import (
	"fmt"
	"regexp"
	"strings"

	"SecureJS/config"
	"SecureJS/internal/parser"
)

// MatchItem 表示单条命中结果
type MatchItem struct {
	RuleName    string // 命中的规则名称
	MatchedText string // 实际匹配到的敏感信息片段
}

// MatchResult 表示对某个 URL 的匹配结果
type MatchResult struct {
	URL   string      // 目标URL
	Items []MatchItem // 命中的所有结果
	Error error       // 如果在匹配过程中有什么错误，可记录在这里（一般不会有）
}

// compiledRule 用于保存编译后的正则，避免重复编译
type compiledRule struct {
	Name  string
	Regex *regexp.Regexp
}

// MatchAll 对从 parser 获得的一组响应内容进行匹配，
// 返回每个 URL 对应的匹配情况。
func MatchAll(rules []config.Rule, parseResults []*parser.ParseResult) ([]*MatchResult, error) {
	// 1) 先编译所有规则（减少重复编译）
	compiledRules := make([]compiledRule, 0, len(rules))
	for _, r := range rules {
		// 如果 f_regex 为空或无效，可以跳过
		if r.FRegex == "" {
			continue
		}
		re, err := regexp.Compile(r.FRegex)
		if err != nil {
			return nil, fmt.Errorf("failed to compile regex for rule '%s': %w", r.Name, err)
		}
		compiledRules = append(compiledRules, compiledRule{
			Name:  r.Name,
			Regex: re,
		})
	}

	// 2) 对每个 parseResult 的 Body 做匹配
	results := make([]*MatchResult, 0, len(parseResults))
	for _, pr := range parseResults {
		if pr.Error != nil {
			// 如果 parse 出错了，这里就直接记录错误
			results = append(results, &MatchResult{
				URL:   pr.URL,
				Error: pr.Error,
			})
			continue
		}

		//去重
		uniqueMatches := make(map[string]bool)
		// 准备收集此 URL 下所有命中项
		var matchedItems []MatchItem
		body := pr.Body

		// 对所有规则匹配
		for _, cr := range compiledRules {
			allMatches := cr.Regex.FindAllString(body, -1)
			for _, matchStr := range allMatches {
				if _, exists := uniqueMatches[matchStr]; !exists {
					uniqueMatches[matchStr] = true
					matchedItems = append(matchedItems, MatchItem{
						RuleName:    cr.Name,
						MatchedText: matchStr,
					})
				}
			}
		}
		// 过滤匹配项
		matchedItems = filterMatchedItems(matchedItems)

		// 如果有匹配项，记录结果
		if len(matchedItems) > 0 {
			results = append(results, &MatchResult{
				URL:   pr.URL,
				Items: matchedItems,
				Error: nil,
			})
		}
	}
	// 输出匹配结果
	// var lastURL string
	// for _, result := range results {
	// 	for _, item := range result.Items {
	// 		if item.MatchedText != "" {
	// 			if lastURL != result.URL {
	// 				if lastURL != "" {
	// 					fmt.Println() // 为了分隔不同的URL输出
	// 				}
	// 				fmt.Printf("URL: %s\n", result.URL)
	// 				lastURL = result.URL
	// 			}
	// 			fmt.Printf("  Rule: %s, Matched: %s\n", item.RuleName, item.MatchedText)
	// 		}
	// 	}
	// }
	return results, nil
}

// filterMatchedItems 只会针对“找到的第一个 : 或 =”进行拆分。
// 拆分后的 key/value 如果包含了指定的关键词或中文字符，就过滤掉。
func filterMatchedItems(matchedItems []MatchItem) []MatchItem {
    // 1) 设定普通过滤关键词（子串匹配、忽略大小写）
    filterKeywords := []string{
        "xml",
    }

    // 2) 前置过滤关键词（只要 key 包含这些，就过滤）
    preFilterKeywords := []string{
        "passive",
    }

    // 3) 编译对应的正则，用于子串匹配和忽略大小写
    keywordRegex := regexp.MustCompile(`(?i)(` + strings.Join(filterKeywords, "|") + `)`)
    preKeywordRegex := regexp.MustCompile(`(?i)(` + strings.Join(preFilterKeywords, "|") + `)`)

    // 4) 检查是否有中文字符
    chineseRegex := regexp.MustCompile(`[\p{Han}]`)

    var filteredItems []MatchItem

    for _, item := range matchedItems {
        // 默认不过滤
        shouldFilter := false

        // 取出整条文本
        text := strings.TrimSpace(item.MatchedText)

        // 找到第一个 ':' 或 '='
        idxColon := strings.IndexRune(text, ':')
        idxEqual := strings.IndexRune(text, '=')

        var splitPos int
        switch {
        case idxColon < 0 && idxEqual < 0:
            // 连一个 ':' 或 '=' 都没有，说明无法拆分
            // 你可以选择：直接保留 (shouldFilter=false)，也可以选择直接过滤
            splitPos = -1
        case idxColon < 0:
            // 没有冒号，只找到等号
            splitPos = idxEqual
        case idxEqual < 0:
            // 没有等号，只找到冒号
            splitPos = idxColon
        default:
            // 同时存在 ':' 和 '='，取最小位置 => 谁先出现
            if idxColon < idxEqual {
                splitPos = idxColon
            } else {
                splitPos = idxEqual
            }
        }

        if splitPos != -1 {
            // 说明找到了分隔符，进行拆分
            key := text[:splitPos]
            value := text[splitPos+1:] // 从分隔符的下一个字符开始到末尾都算 value

            // 转小写并去除首尾的引号、空格
            key = strings.ToLower(strings.Trim(key, `"' `))
            value = strings.ToLower(strings.Trim(value, `"' `))

            // ① 如果 key 命中“前置过滤关键词”，则过滤
            if preKeywordRegex.MatchString(key) {
                shouldFilter = true
            }

            // ② value 中出现“普通过滤关键词”，则过滤
            if keywordRegex.MatchString(value) {
                shouldFilter = true
            }

            // ③ value 中包含中文字符，也过滤
            if chineseRegex.MatchString(value) {
                shouldFilter = true
            }
        }

        if !shouldFilter {
            filteredItems = append(filteredItems, item)
        }
    }

    return filteredItems
}
