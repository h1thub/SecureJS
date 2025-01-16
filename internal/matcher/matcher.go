package matcher

import (
	"fmt"
	"regexp"

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

		results = append(results, &MatchResult{
			URL:   pr.URL,
			Items: matchedItems,
			Error: nil,
		})
	}

	return results, nil
}
