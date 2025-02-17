# SecureJS

SecureJS 是一个强大的工具，旨在从目标网站收集所有相关链接，对这些链接（主要是 JavaScript 文件）执行请求，并扫描敏感信息，如令牌、密钥、密码、AKSK 等，后引入 DeepSeek 对结果进行分析。

## 目录

- [SecureJS](#securejs)
  - [目录](#目录)
  - [使用方法](#使用方法)
    - [帮助信息](#帮助信息)
    - [示例](#示例)
  - [配置](#配置)
  - [项目结构](#项目结构)
  - [免责声明](#免责声明)

## 使用方法

### 帮助信息

```
Usage:
  SecureJS [flags]

Flags:
  -a, --ai string            true/false. Enable AI Analytics. If not set, will use false (default "false")
  -b, --browser string       Path to Chrome/Chromium executable (optional). If not set, will use Rod's default.
  -c, --config string        Path to config file (e.g. config.yaml) (default "config/config.yaml")
  -H, --header stringArray   Add custom request headers. (e.g. -H 'Key: Value')
  -h, --help                 help for SecureJS
  -i, --id string            YOUR_ENDPOINT_ID
  -k, --key string           ARK_API_KEY
  -l, --list string          File containing target URLs (one per line)
  -o, --output string        Output file (supports .txt, .csv, .json)
  -p, --proxy string         Proxy to use (e.g. http://127.0.0.1:8080)
  -t, --threads int          Number of concurrent threads for scanning (default 20)
  -u, --url string           Single target URL to scan (e.g. https://example.com)
```

### 示例


## 配置

SecureJS 使用 `config/config.yaml` 文件来定义自定义匹配规则和其他项目级配置。如不存在，首次运行后将自动生成该文件。另外，规则将进行尽可能的匹配，因为后续可进行AI分析，但这也会导致AI分析前误报结果高。

## 项目结构

```
SecureJS/
├── cmd/
│   └── root.go             # 处理命令行参数（-u、-l、-t 等）的入口点
│
├── internal/
│   ├── crawler/
│   │   └── ai.go      # 引入 DeepSeek 对结果二次分析
│   │
│   ├── crawler/
│   │   ├── crawler.go      # 爬虫逻辑，模拟浏览器访问，收集所有链接和 JS 文件
│   │   └── linkfind.go     # 从目标页面的响应体中提取所有链接和 JS
│   │
│   ├── parser/
│   │   └── parser.go       # 对所有收集的链接和 JS 文件执行二次请求
│   │
│   ├── matcher/
│   │   └── matcher.go      # 从 config.yaml 中读取并解析自定义规则，并与响应体匹配
│   │
│   └── output/
│       └── output.go       # 将结果输出为 CSV、JSON 或文本格式的文件
│
├── config/
│   ├── config.go           # 处理配置文件（config.yaml）的加载和解析
│   └── config.yaml         # 自定义规则和其他项目级配置
│
├── go.mod                  # Go Modules 管理文件
├── go.sum                  # Go Modules 校验文件
└── main.go                 # 主程序入口点，初始化并启动应用程序
```

## 免责声明

本工具仅用于安全研究与合法测试目的。请确保遵守相关法律法规，不得将本工具用于任何非法或未经授权的行为。作者及项目维护者对因使用或滥用本工具而导致的任何损失或损害，不承担任何责任。