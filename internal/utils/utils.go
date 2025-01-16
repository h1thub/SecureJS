package utils

import (
    "bufio"
    "os"
    "strings"
)

// ReadURLs 从指定文件中读取非空行并将其添加到传入的 urls 切片中
func ReadURLs(filePath string, urls []string) ([]string, error) {
    // 打开文件
    file, err := os.Open(filePath)
    if err != nil {
        return urls, err
    }
    defer file.Close()

    // 使用 bufio.Scanner 逐行读取文件
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line != "" {
            urls = append(urls, line)
        }
    }

    // 检查扫描过程中是否有错误
    if err := scanner.Err(); err != nil {
        return urls, err
    }

    return urls, nil
}


//定义一个包含所有需要过滤的子字符串的切片，Skip 检查 URL 是否包含任何需要跳过的子字符串
var skipSubstrings = []string{
    "static.ocecdn.oraclecloud.com",
    "google",
    "baidu",
    "data:image",
    "www.youtube.com",
    "www.facebook.com",
    "www.w3.org",
    "twitter.com",
}
func Skip(url string) bool {
    for _, substr := range skipSubstrings {
        if strings.Contains(url, substr) {
            return true
        }
    }
    return false
}

// skipExtensions 是需要过滤掉的所有后缀(小写)，HasSkipExtension 函数判断请求 URL 是否包含需要过滤的后缀
var skipExtensions = map[string]bool{
	".3g2": true, ".3gp": true, ".7z": true, ".aac": true,
	".abw": true, ".aif": true, ".aifc": true, ".aiff": true,
	".apk": true, ".arc": true, ".au": true, ".avi": true,
	".azw": true, ".bat": true, ".bin": true, ".bmp": true,
	".bz": true, ".bz2": true, ".cmd": true, ".cmx": true,
	".cod": true, ".com": true, ".csh": true, ".css": true,
	".csv": true, ".dll": true, ".doc": true, ".docx": true,
	".ear": true, ".eot": true, ".epub": true, ".exe": true,
	".flac": true, ".flv": true, ".gif": true, ".gz": true,
	".ico": true, ".ics": true, ".ief": true, ".jar": true,
	".jfif": true, ".jpe": true, ".jpeg": true, ".jpg": true,
	".less": true, ".m3u": true, ".mid": true, ".midi": true,
	".mjs": true, ".mkv": true, ".mov": true, ".mp2": true,
	".mp3": true, ".mp4": true, ".mpa": true, ".mpe": true,
	".mpeg": true, ".mpg": true, ".mpkg": true, ".mpp": true,
	".mpv2": true, ".odp": true, ".ods": true, ".odt": true,
	".oga": true, ".ogg": true, ".ogv": true, ".ogx": true,
	".otf": true, ".pbm": true, ".pdf": true, ".pgm": true,
	".png": true, ".pnm": true, ".ppm": true, ".ppt": true,
	".pptx": true, ".ra": true, ".ram": true, ".rar": true,
	".ras": true, ".rgb": true, ".rmi": true, ".rtf": true,
	".scss": true, ".sh": true, ".snd": true, ".svg": true,
	".swf": true, ".tar": true, ".tif": true, ".tiff": true,
	".ttf": true, ".vsd": true, ".war": true, ".wav": true,
	".weba": true, ".webm": true, ".webp": true, ".wmv": true,
	".woff": true, ".woff2": true, ".xbm": true, ".xls": true,
	".xlsx": true, ".xpm": true, ".xul": true, ".xwd": true,
	".zip": true,
}
func HasSkipExtension(reqURL string) bool {
	beforeQuery := strings.Split(reqURL, "?")[0]
	beforeHash := strings.Split(beforeQuery, "#")[0]

	slashIndex := strings.LastIndex(beforeHash, "/")
	filename := beforeHash
	if slashIndex != -1 {
		filename = beforeHash[slashIndex+1:]
	}

	for ext := range skipExtensions {
		if strings.HasSuffix(filename, ext) {
			return true
		}
	}
	return false
}