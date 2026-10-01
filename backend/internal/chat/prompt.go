package chat

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ChatCoachSystemPrompt 多轮教练：主输出自然语言，可选附带结构化 JSON
const ChatCoachSystemPrompt = `你是一名专业的算法教练，通过多轮对话帮助学生理解和解决 LeetCode 算法题。

## 职责
1. 根据题目、代码和对话历史，给出下一步思考方向
2. 用苏格拉底式提问引导学生，而不是直接给完整题解
3. 可以指出逻辑问题、边界情况和复杂度问题
4. 回复使用自然中文，适合聊天气泡阅读

## 禁止
- 不要直接贴完整 AC 代码
- 不要一次性给出最优完整解法
- 不要忽略用户的上一轮问题

## 输出格式（重要）
1. 先写给用户看的自然语言回复（必有）
2. 若本轮适合补充结构化信息，在回复末尾单独追加一个 JSON 代码块，格式如下：

` + "```json" + `
{"hints":["..."],"issues":["..."],"complexityComment":"..."}
` + "```" + `

3. JSON 可选；没有把握就只输出自然语言
4. 不要把 JSON 当作唯一输出`

// BuildProblemCard 把题目上下文压成一条稳定的上下文消息
func BuildProblemCard(c ProblemContext) string {
	var sb strings.Builder
	sb.WriteString("以下是当前题目上下文，后续对话都围绕它进行：\n")
	sb.WriteString(fmt.Sprintf("标题: %s\n", c.Title))
	sb.WriteString(fmt.Sprintf("链接: %s\n", c.URL))
	sb.WriteString(fmt.Sprintf("站点: %s\n", c.Site))
	sb.WriteString(fmt.Sprintf("语言: %s\n", c.Language))
	if c.Description != "" {
		sb.WriteString("\n题目描述:\n")
		sb.WriteString(c.Description)
		sb.WriteString("\n")
	}
	if c.Code != "" {
		sb.WriteString("\n学生当前代码:\n```")
		sb.WriteString(c.Language)
		sb.WriteString("\n")
		sb.WriteString(c.Code)
		sb.WriteString("\n```\n")
	}
	return sb.String()
}

// AssistantMeta 可选结构化字段
type AssistantMeta struct {
	Hints             []string `json:"hints,omitempty"`
	Issues            []string `json:"issues,omitempty"`
	ComplexityComment string   `json:"complexityComment,omitempty"`
}

// SplitAssistantOutput 拆出自然语言主回复 + 可选 JSON meta
func SplitAssistantOutput(raw string) (reply string, meta *AssistantMeta) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", nil
	}

	start := strings.LastIndex(raw, "```json")
	if start < 0 {
		start = strings.LastIndex(raw, "```")
	}
	if start < 0 {
		return raw, nil
	}

	rest := raw[start:]
	rest = strings.TrimPrefix(rest, "```json")
	rest = strings.TrimPrefix(rest, "```")
	end := strings.Index(rest, "```")
	jsonPart := rest
	if end >= 0 {
		jsonPart = rest[:end]
	}
	jsonPart = strings.TrimSpace(jsonPart)

	var parsed AssistantMeta
	if err := json.Unmarshal([]byte(jsonPart), &parsed); err != nil {
		return raw, nil
	}

	reply = strings.TrimSpace(raw[:start])
	if reply == "" {
		reply = raw
	}
	if len(parsed.Hints) == 0 && len(parsed.Issues) == 0 && parsed.ComplexityComment == "" {
		return reply, nil
	}
	return reply, &parsed
}
