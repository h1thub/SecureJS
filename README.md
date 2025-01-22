# SecureJS

SecureJS 是一个强大的工具，旨在从目标网站收集所有相关链接，对这些链接（主要是 JavaScript 文件）执行请求，并扫描敏感信息，如令牌、密钥、密码、AKSK 等。

## 目录

- [SecureJS](#securejs)
  - [目录](#目录)
  - [功能](#功能)
  - [使用方法](#使用方法)
    - [示例](#示例)
  - [配置](#配置)
  - [项目结构](#项目结构)

## 功能

- **全面爬取**：模拟浏览器访问，收集目标网站的所有链接和 JavaScript 文件。
- **二次请求**：对收集的资源执行额外的请求以进行更深入的分析。
- **可自定义的匹配规则**：支持在 `config.yaml` 中定义的自定义规则，以识别敏感信息。
- **灵活的输出格式**：将结果输出为 CSV、JSON 或纯文本格式。
- **简易配置**：通过配置文件简化设置和自定义过程。

## 使用方法

SecureJS 可以通过命令行执行，并提供各种选项以自定义其行为。

### 示例

```bash
./SecureJS -u https://example.com -o results.csv
```

```bash
./SecureJS -l targets.txt -o results.csv -t 30
```

## 配置

SecureJS 使用 `config/config.yaml` 文件来定义自定义匹配规则和其他项目级配置。

## 项目结构

```
SecureJS/
├── cmd/
│   └── root.go             # 处理命令行参数（-u、-l、-t 等）的入口点
│
├── internal/
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

**免责声明**: 

本工具仅用于安全研究与合法测试目的。请确保遵守相关法律法规，不得将本工具用于任何非法或未经授权的行为。作者及项目维护者对因使用或滥用本工具而导致的任何损失或损害，不承担任何责任。