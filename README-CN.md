# SecureJS

SecureJS 是一个强大的工具，旨在从目标网站收集所有相关链接，对这些链接（主要是 JavaScript 文件）执行请求，并扫描敏感信息，如令牌、密钥、密码、AKSK 等。

## 目录

- [SecureJS](#securejs)
  - [目录](#目录)
  - [功能](#功能)
  - [安装](#安装)
    - [前提条件](#前提条件)
    - [步骤](#步骤)
  - [使用方法](#使用方法)
    - [示例](#示例)
  - [配置](#配置)
    - [示例 `config.yaml`](#示例-configyaml)
    - [加载配置](#加载配置)
  - [项目结构](#项目结构)

## 功能

- **全面爬取**：模拟浏览器访问，收集目标网站的所有链接和 JavaScript 文件。
- **二次请求**：对收集的资源执行额外的请求以进行更深入的分析。
- **可自定义的匹配规则**：支持在 `config.yaml` 中定义的自定义规则，以识别敏感信息。
- **灵活的输出格式**：将结果输出为 CSV、JSON 或纯文本格式。
- **简易配置**：通过配置文件简化设置和自定义过程。

## 安装

### 前提条件

- [Go](https://golang.org/dl/) 1.16 或更高版本

### 步骤

1. **克隆仓库**

   ```bash
   git clone
   cd SecureJS
   ```

2. **构建应用程序**

   ```bash
   go build
   ```

3. **验证安装**

   ```bash
   ./SecureJS -h
   ```

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

SecureJS 使用 `config.yaml` 文件来定义自定义匹配规则和其他项目级配置。

### 示例 `config.yaml`

```yaml
rules:
  - name: Sensitive Field
    f_regex: (?i)\[?["']?[0-9A-Za-z_-]{0,15}(?:key|secret|token|config|auth|access|admin|ticket)[0-9A-Za-z_-]{0,15}["']?\]?\s*(?:=|:|\)\.val\()\s*\[?\{?["']([^"']{8,256})["']?(?::|,)?

  - name: Password Field
    f_regex: ((|\\)(|'|")(|[\w]{1,10})([p](ass|wd|asswd|assword))(|[\w]{1,10})(|\\)(|'|")(:|=|\)\.val\()(|)(|\\)('|")([^'"]+?)(|\\)('|")(|,|\)))

  - name: JSON Web Token
    f_regex: (eyJ[A-Za-z0-9_-]{10,}\.[A-Za-z0-9._-]{10,}|eyJ[A-Za-z0-9_\/+-]{10,}\.[A-Za-z0-9._\/+-]{10,})

  - name: Cloud Key
    f_regex: (?i)(?:AWSAccessKeyId=[A-Z0-9]{16,32}|access[-_]?key[-_]?(?:id|secret)|LTAI[a-z0-9]{12,20})
```

### 加载配置

配置会自动从 `config/config.yaml` 文件中加载。请确保您的自定义规则已正确定义，以匹配您希望识别的敏感信息。

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