package analyze

import (
	"context"
	"encoding/json"
	"fmt"

	"aigo-coach/backend/internal/ai"
)

// AnalyzeService 封装代码分析的业务逻辑
type AnalyzeService struct {
	aiClient ai.Client // AI 客户端（可以是 DeepSeek、OpenAI 等）
}

// NewAnalyzeService 创建分析服务实例
func NewAnalyzeService(client ai.Client) *AnalyzeService {
	return &AnalyzeService{
		aiClient: client,
	}
}

// AnalyzeCode 分析用户代码并返回教练式建议
// 这是核心业务方法，协调 prompt 构建和 AI 调用
func (s *AnalyzeService) AnalyzeCode(ctx context.Context, req AnalyzeRequest) (AnalyzeResponse, error) {
	// 1. 构建教练模式的对话上下文
	messages := []ai.Message{
		{
			Role:    "system",
			Content: CoachSystemPrompt, // 定义教练角色和行为
		},
		{
			Role:    "user",
			Content: BuildCoachPrompt(req), // 题目 + 代码
		},
	}

	// 2. 配置模型参数
	config := ai.ChatConfig{
		Model:       "deepseek-chat", // DeepSeek 主力模型
		Temperature: 0.7,             // 适中的随机性，既有创造性又稳定
		MaxTokens:   2000,            // 足够返回完整分析
	}

	// 3. 调用 AI 客户端
	content, err := s.aiClient.Chat(ctx, messages, config)
	if err != nil {
		// AI 调用失败，返回错误而不是 mock 数据
		return AnalyzeResponse{
			Result: false,
			Data:   data{},
		}, fmt.Errorf("AI 调用失败: %w", err)
	}

	// 4. 解析 AI 返回的 JSON 格式
	var result struct {
		Summary           string   `json:"summary"`
		Issues            []string `json:"issues"`
		Hints             []string `json:"hints"`
		ComplexityComment string   `json:"complexityComment"`
	}

	// 清理可能的 markdown 代码块标记
	content = cleanMarkdownJSON(content)

	if err := json.Unmarshal([]byte(content), &result); err != nil {
		// 如果模型没按格式返回，降级处理：把原文当 summary
		return AnalyzeResponse{
			Result: true,
			Data: data{
				Summary:    "AI 返回了非结构化内容",
				Problem:    []string{},
				Hints:      []string{content}, // 把原文放到 hints 里
				Complexity: "",
			},
		}, nil
	}

	// 5. 组装返回结果
	return AnalyzeResponse{
		Result: true,
		Data: data{
			Summary:    result.Summary,
			Problem:    result.Issues,
			Hints:      result.Hints,
			Complexity: result.ComplexityComment,
		},
	}, nil
}

// cleanMarkdownJSON 移除模型可能添加的 markdown 代码块标记
// 例如: ```json\n{...}\n``` -> {...}
func cleanMarkdownJSON(s string) string {
	// 移除前后空格
	s = trimSpace(s)

	// 移除开头的 ```json 或 ```
	if len(s) > 7 && s[:7] == "```json" {
		s = s[7:]
	} else if len(s) > 3 && s[:3] == "```" {
		s = s[3:]
	}

	// 移除结尾的 ```
	s = trimSpace(s)
	if len(s) > 3 && s[len(s)-3:] == "```" {
		s = s[:len(s)-3]
	}

	return trimSpace(s)
}

// trimSpace 简易的前后空白字符去除（避免引入 strings 包依赖）
func trimSpace(s string) string {
	start := 0
	end := len(s)

	// 去掉开头空白
	for start < end && (s[start] == ' ' || s[start] == '\n' || s[start] == '\t' || s[start] == '\r') {
		start++
	}

	// 去掉结尾空白
	for end > start && (s[end-1] == ' ' || s[end-1] == '\n' || s[end-1] == '\t' || s[end-1] == '\r') {
		end--
	}

	return s[start:end]
}
