package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// DeepSeekClient 实现 Client 接口，调用 DeepSeek API
type DeepSeekClient struct {
	apiKey     string       // DeepSeek API Key
	baseURL    string       // API 基础 URL
	httpClient *http.Client // 复用 HTTP 连接池
}

// NewDeepSeekClient 创建 DeepSeek 客户端实例
func NewDeepSeekClient(apiKey string) *DeepSeekClient {
	return &DeepSeekClient{
		apiKey:  apiKey,
		baseURL: "https://api.deepseek.com/v1", // DeepSeek 官方 API 地址
		httpClient: &http.Client{
			Timeout: 60 * time.Second, // 60 秒超时，避免长时间挂起
		},
	}
}

// deepseekRequest 对应 DeepSeek API 请求体结构
// 参考: https://platform.deepseek.com/api-docs/
type deepseekRequest struct {
	Model       string    `json:"model"`       // 必填: deepseek-chat
	Messages    []Message `json:"messages"`    // 必填: 对话消息列表
	Temperature float32   `json:"temperature"` // 可选: 0-2
	MaxTokens   int       `json:"max_tokens"`  // 可选: 最大输出 token 数
}

// deepseekResponse 对应 DeepSeek API 响应体结构
type deepseekResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"` // 这是模型回复的核心内容
		} `json:"message"`
		FinishReason string `json:"finish_reason"` // stop 表示正常结束
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// deepseekErrorResponse 对应错误响应
type deepseekErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`
}

// Chat 实现 Client 接口的 Chat 方法
func (c *DeepSeekClient) Chat(ctx context.Context, messages []Message, config ChatConfig) (string, error) {
	// 1. 构造请求体
	reqBody := deepseekRequest{
		Model:       config.Model,
		Messages:    messages,
		Temperature: config.Temperature,
		MaxTokens:   config.MaxTokens,
	}

	// 2. 序列化为 JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	// 3. 创建 HTTP 请求
	url := c.baseURL + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建 HTTP 请求失败: %w", err)
	}

	// 4. 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey) // DeepSeek 使用 Bearer 鉴权

	// 5. 发送请求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送 HTTP 请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 6. 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	// 7. 处理非 200 状态码（API 错误）
	if resp.StatusCode != http.StatusOK {
		var errResp deepseekErrorResponse
		if err := json.Unmarshal(body, &errResp); err != nil {
			// 无法解析错误响应，返回原始内容
			return "", fmt.Errorf("API 错误 (HTTP %d): %s", resp.StatusCode, string(body))
		}
		return "", fmt.Errorf("DeepSeek API 错误: %s (type: %s, code: %s)",
			errResp.Error.Message, errResp.Error.Type, errResp.Error.Code)
	}

	// 8. 解析成功响应
	var successResp deepseekResponse
	if err := json.Unmarshal(body, &successResp); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	// 9. 提取模型回复内容
	if len(successResp.Choices) == 0 {
		return "", fmt.Errorf("响应中没有 choices")
	}

	content := successResp.Choices[0].Message.Content
	if content == "" {
		return "", fmt.Errorf("模型返回空内容")
	}

	return content, nil
}
