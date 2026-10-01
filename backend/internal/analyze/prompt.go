package analyze

import (
	"fmt"
	"strings"
)

// CoachSystemPrompt 定义算法教练的角色和行为规则
// 这是发送给 AI 的 system 消息，控制它的回答风格
const CoachSystemPrompt = `你是一名专业的算法教练，帮助学生理解和解决 LeetCode 算法题。

## 你的职责
1. 分析学生的代码，找出逻辑问题、性能瓶颈和边界情况
2. 给出提示和思路引导，而不是直接给出完整答案
3. 用苏格拉底式提问，帮助学生自己思考出解决方案
4. 解释算法的时间复杂度和空间复杂度

## 你不应该做的
- 不要直接贴出完整的 AC 代码
- 不要一次性给出最优解
- 不要替学生完成整道题
- 不要在学生还没尝试时就给答案                                         

## 回答格式
请按以下 JSON 格式返回（必须是合法 JSON，不要用 markdown 代码块包裹）：
{
  "summary": "一句话总结当前代码的核心问题或思路",
  "issues": ["问题1", "问题2"],                                                                  
  "hints": ["提示1：引导性问题", "提示2：关键观察"],
  "complexityComment": "时间/空间复杂度分析和优化方向"
}

## 示例
学生代码：两层循环暴力枚举
你的回答：
{
  "summary": "你的思路是正确的，但两层循环时间复杂度较高",
  "issues": [
    "当前是 O(n²) 时间复杂度，数据量大时会超时",
    "没有利用「目标值固定」这个条件"
  ],
  "hints": [
    "遍历到 nums[i] 时，你真正需要找的另一个数是多少？",
    "如果能用 O(1) 时间判断某个数是否出现过，会怎样？",
    "哈希表可以做什么？"
  ],
  "complexityComment": "目标是 O(n) 时间，O(n) 空间"
}`

// BuildCoachPrompt 构造完整的教练对话上下文
// 将题目信息和用户代码组装成模型能理解的格式
func BuildCoachPrompt(req AnalyzeRequest) string {
	var sb strings.Builder

	// 1. 题目基本信息
	sb.WriteString("## 题目信息\n")
	sb.WriteString(fmt.Sprintf("**标题**: %s\n", req.Title))
	sb.WriteString(fmt.Sprintf("**链接**: %s\n", req.URL))
	sb.WriteString("\n")

	// 2. 题目描述（如果有）
	if req.Description != "" {
		sb.WriteString("## 题目描述\n")
		sb.WriteString(req.Description)
		sb.WriteString("\n\n")
	}

	// 3. 学生代码
	sb.WriteString("## 学生代码\n")
	sb.WriteString(fmt.Sprintf("**编程语言**: %s\n", req.Language))
	sb.WriteString("```" + req.Language + "\n")
	sb.WriteString(req.Code)
	sb.WriteString("\n```\n\n")

	// 4. 明确要求（强化 system prompt）
	sb.WriteString("## 要求\n")
	sb.WriteString("请分析上述代码，按 JSON 格式返回教练式的分析结果。")
	sb.WriteString("记住：给提示不给答案，引导学生自己思考。\n")

	return sb.String()
}

// ParseCoachResponse 解析模型返回的 JSON 字符串
// 如果模型返回的不是合法 JSON，做降级处理
func ParseCoachResponse(content string) (AnalyzeResponse, error) {
	// 移除可能的 markdown 代码块标记（有些模型会加 ```json）
	content = strings.TrimSpace(content)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	// 尝试解析为结构化响应
	var result struct {
		Summary           string   `json:"summary"`
		Issues            []string `json:"issues"`
		Hints             []string `json:"hints"`
		ComplexityComment string   `json:"complexityComment"`
	}

	// 这里简化处理：直接用 json.Unmarshal
	// 实际项目中可以用 encoding/json
	// 为了不依赖外部，这里先返回一个占位实现
	// TODO: 在调用处用 json.Unmarshal 解析

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
