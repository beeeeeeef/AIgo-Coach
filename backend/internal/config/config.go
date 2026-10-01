package config

type Config struct {
	AIs []AIConfig  `json:"ais"`  
}

type AIConfig struct {
	// 基础信息
	Name     string `json:"name" yaml:"name"`         // 标识名称，如 "openai", "deepseek"
	Provider string `json:"provider" yaml:"provider"` // 提供商类型，用于选择不同的客户端实现

	// 认证
	APIKey  string `json:"api_key" yaml:"api_key"`
	BaseURL string `json:"base_url" yaml:"base_url"` // 支持自定义端点或代理

	// 模型配置
	Model       string  `json:"model" yaml:"model"`             // 如 "gpt-4", "deepseek-coder"
	Temperature float64 `json:"temperature" yaml:"temperature"` // 0.0-1.0
	MaxTokens   int     `json:"max_tokens" yaml:"max_tokens"`
	Enabled     bool    `json:"enabled" yaml:"enabled"`
}
