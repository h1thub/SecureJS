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
  - name: Sensitive Field
    f_regex: (?i)\[?["']?[0-9A-Za-z_-]{0,15}(?:key|secret|token|config|auth|access|admin|ticket)[0-9A-Za-z_-]{0,15}["']?\]?\s*(?:=|:|\)\.val\()\s*\[?\{?["']([^"']{8,256})["']?(?::|,)?

  - name: Password Field
    f_regex: ((|\\)(|'|")(|[\w]{1,10})([p](ass|wd|asswd|assword))(|[\w]{1,10})(|\\)(|'|")(:|=|\)\.val\()(|)(|\\)('|")([^'"]+?)(|\\)('|")(|,|\)))

  - name: JSON Web Token
    f_regex: (eyJ[A-Za-z0-9_-]{10,}\.[A-Za-z0-9._-]{10,}|eyJ[A-Za-z0-9_\/+-]{10,}\.[A-Za-z0-9._\/+-]{10,})

  - name: Cloud Key
    f_regex: (?i)(?:AWSAccessKeyId=[A-Z0-9]{16,32}|access[-_]?key[-_]?(?:id|secret)|LTAI[a-z0-9]{12,20})
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
