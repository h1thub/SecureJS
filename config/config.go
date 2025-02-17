package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Rule 表示一条匹配规则，对应 config.yaml 里的每个条目
type Rule struct {
	Name   string `yaml:"name"`
	FRegex string `yaml:"f_regex"`
}

// Config 表示整个配置文件内容，里面是若干 Rule
type Config struct {
	Rules []Rule `yaml:"rules"`
}

// LoadConfig 从指定路径的 YAML 文件中加载配置，若文件不存在则创建并写入默认配置，返回 *Config
func LoadConfig(path string) (*Config, error) {
	// 1. 尝试读取文件
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// 文件不存在，创建并写入默认配置

			// 获取文件的目录路径
			dir := filepath.Dir(path)

			// 创建目录（如果不存在）
			if err := os.MkdirAll(dir, 0755); err != nil {
				return nil, fmt.Errorf("无法创建目录 '%s': %w", dir, err)
			}

			// 定义默认配置内容
			defaultContent := `
rules:
  - name: Extended Sensitive Field
    f_regex: "(?i)([\"']?[\\w-]{0,15}(?:key|secret|token|config|auth|access|admin|ticket|api_key|client_secret|private_key|public_key|bearer|session|cookie|license|cert|ssh|salt|pepper)[\\w-]{0,15}[\"']?)\\s*(?:=|:|\\)\\.val\\()\\s*\\[?\\{?(?:'([^']{8,500})'|\"([^\\\"]{8,500})\")(?:[:;,\\}\\]]?)?"

  - name: Extended Password Field
    f_regex: "(?i)((|\\\\)(?:'|\")[\\w-]{0,10}(?:p(?:ass|wd|asswd|assword|asscode|assphrase)|secret)[\\w-]{0,10}(|\\\\)(?:'|\"))\\s*(?:=|:|\\)\\.val\\()(|)(|\\\\)(?:'|\")([^'\"]+?)(|\\\\)(?:'|\")(?:|,|\\)|;)?"

  - name: Extended JSON Web Token
    f_regex: "(?i)(eyJ[A-Za-z0-9_-]{5,}\\.[A-Za-z0-9._-]{5,}\\.[A-Za-z0-9._-]{5,})"

  - name: Extended Cloud Key
    f_regex: "(?i)(AWSAccessKeyId=[A-Z0-9]{16,32}|access[-_]?key[-_]?(?:id|secret)|LTAI[a-z0-9]{12,20}|(?:AKIA|ABIA|ACCA|ASIA)[A-Z0-9]{16}|aws_secret_access_key\\s*=\\s*[\"'][^\"']{8,100}[\"'])"

  - name: Azure Key
    f_regex: "(?i)(AZURE_STORAGE[_-]?ACCOUNT[_-]?KEY|AZURE_STORAGE_KEY|AZURE_KEY_VAULT|azure_tenant_id)\\s*=\\s*[\"']([^\"']{8,100})[\"']"

  - name: GCP Service Account
    f_regex: "(?s)(\"type\"\\s*:\\s*\"service_account\".*?\"private_key_id\"\\s*:\\s*\"([a-z0-9]{10,})\".*?\"private_key\"\\s*:\\s*\"-----BEGIN PRIVATE KEY-----.*?-----END PRIVATE KEY-----\")"

  - name: Private Key
    f_regex: "(?s)-----BEGIN\\s+(?:RSA|EC|DSA|OPENSSH)?\\s*PRIVATE\\s+KEY-----.*?-----END\\s+(?:RSA|EC|DSA|OPENSSH)?\\s*PRIVATE\\s+KEY-----"
`

			// 创建文件并写入默认内容
			err = os.WriteFile(path, []byte(defaultContent), 0644)
			if err != nil {
				return nil, fmt.Errorf("无法创建默认配置文件 '%s': %w", path, err)
			}

			// 重新读取刚创建的文件
			data, err = os.ReadFile(path)
			if err != nil {
				return nil, fmt.Errorf("创建后无法读取配置文件 '%s': %w", path, err)
			}
		} else {
			// 其他读取错误
			return nil, fmt.Errorf("读取配置文件 '%s' 失败: %w", path, err)
		}
	}

	// 2. 解析 YAML
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析 YAML 失败: %w", err)
	}

	return &cfg, nil
}
