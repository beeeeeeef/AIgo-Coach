package ai

import (
	"context"
)

// Message 代表一条对话消息
// 遵循 OpenAI/DeepSeek 标准格式
type Message struct {
	Role    string `json:"role"`    // system, user, assistant
	Content string `json:"content"` // 消息内容
}

// ChatConfig 控制模型行为的配置参数
type ChatConfig struct {
	Model       string  `json:"model"`       // 模型名称，如 deepseek-chat
	Temperature float32 `json:"temperature"` // 0-2，越高越随机
	MaxTokens   int     `json:"max_tokens"`  // 最大返回 token 数
}

// Client 定义 AI 客户端的通用接口
// 不同厂商（DeepSeek、OpenAI、Claude）都实现这个接口
type Client interface {
	// Chat 发起一次对话请求
	// messages: 包含 system prompt 和用户消息的完整上下文
	// config: 本次请求的模型配置
	// 返回: 模型回复内容和可能的错误
	Chat(ctx context.Context, messages []Message, config ChatConfig) (string, error)
}
