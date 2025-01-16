package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"SecureJS/internal/matcher"
)

// PrintResultsToConsole 在控制台打印结果。
// 如果某个URL没有命中任何信息，仍然会显示"No sensitive info found."
func PrintResultsToConsole(results []*matcher.MatchResult) {
	for _, mr := range results {
		if mr.Error != nil {
			fmt.Printf("\n[!] Parse error on %s: %v\n", mr.URL, mr.Error)
			continue
		}
		if len(mr.Items) == 0 {
			//fmt.Printf("\n[+] %s: no sensitive info found.\n", mr.URL)
			continue
		}
		fmt.Printf("\n[+] %s: found %d item(s)\n", mr.URL, len(mr.Items))
		for _, item := range mr.Items {
			fmt.Printf("    - Rule: %s, Matched: %s\n", item.RuleName, item.MatchedText)
		}
	}
}

// WriteResultsToFile 将匹配结果写入指定文件；如果没有敏感信息则跳过该URL，不写入。
// ext 可以是 ".txt" / ".csv" / ".json"，否则视为 ".txt"。
func WriteResultsToFile(results []*matcher.MatchResult, outPath string) error {
	ext := strings.ToLower(filepath.Ext(outPath))
	if ext == "" {
		ext = ".txt"
	}

	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("create file error: %w", err)
	}
	defer f.Close()

	switch ext {
	case ".txt":
		return writeTxt(results, f)
	case ".csv":
		return writeCSV(results, f)
	case ".json":
		return writeJSON(results, f)
	default:
		// 如果后缀不是这三个，默认按 txt 处理
		return writeTxt(results, f)
	}
}

// writeTxt：只写有敏感信息的记录
func writeTxt(results []*matcher.MatchResult, w io.Writer) error {
	for _, mr := range results {
		// 如果解析出错，可以考虑还是写出来提示一下
		if mr.Error != nil {
			_, _ = fmt.Fprintf(w, "[!] Parse error on %s: %v\n\n", mr.URL, mr.Error)
			continue
		}
		// 如果没有任何命中，直接跳过，不写
		if len(mr.Items) == 0 {
			continue
		}
		// 写有敏感信息的
		_, _ = fmt.Fprintf(w, "[+] %s: found %d item(s)\n", mr.URL, len(mr.Items))
		for _, item := range mr.Items {
			_, _ = fmt.Fprintf(w, "    - Rule: %s, Matched: %s\n", item.RuleName, item.MatchedText)
		}
		_, _ = fmt.Fprintln(w) // 空行分隔
	}
	return nil
}

// writeCSV：只写有敏感信息的记录
func writeCSV(results []*matcher.MatchResult, w io.Writer) error {
	csvWriter := csv.NewWriter(w)
	defer csvWriter.Flush()

	// 写表头
	_ = csvWriter.Write([]string{"URL", "Rule", "MatchedText", "Error"})

	for _, mr := range results {
		if mr.Error != nil {
			// 如果整个页面解析出错，也写一行记录
			_ = csvWriter.Write([]string{mr.URL, "", "", mr.Error.Error()})
			continue
		}
		// 如果没报错但也没有任何命中 -> 跳过
		if len(mr.Items) == 0 {
			continue
		}
		// 写出匹配的条目
		for _, item := range mr.Items {
			_ = csvWriter.Write([]string{mr.URL, item.RuleName, item.MatchedText, ""})
		}
	}
	return nil
}

// writeJSON：只写有敏感信息的记录
func writeJSON(results []*matcher.MatchResult, w io.Writer) error {
	// 先构造一个新的 slice，仅存有匹配的结果
	filtered := make([]*matcher.MatchResult, 0, len(results))
	for _, mr := range results {
		if mr.Error != nil {
			// 即便出错，也可以把它保留，供排查
			filtered = append(filtered, mr)
			continue
		}
		if len(mr.Items) > 0 {
			filtered = append(filtered, mr)
		}
	}
	// 如果全都没有敏感信息，也就会是个空数组 []

	data, err := json.MarshalIndent(filtered, "", "  ")
	if err != nil {
		return fmt.Errorf("json marshal error: %w", err)
	}
	_, _ = w.Write(data)
	return nil
}
