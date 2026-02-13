package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// --- Gemini 专属的数据结构 ---
type GeminiRequest struct {
	Contents []Content `json:"contents"`
}

type Content struct {
	Role  string `json:"role,omitempty"` 
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

// Gemini 返回的数据结构
type GeminiResponse struct {
	Candidates []struct {
		Content Content `json:"content"`
	} `json:"candidates"`
}

// 核心函数：发送给 Gemini
func ChatWithGemini(userCode string, userQuestion string) (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	apiUrl := os.Getenv("GEMINI_API_URL")

	// 准备 Prompt (还是那个严厉的教练人设)
	systemPrompt := `你是一位严厉但循循善诱的 ACM 算法竞赛教练。
	学生会发给你代码和问题。
	1. 请不要直接给出完整代码。
	2. 先分析代码的时间/空间复杂度。
	3. 指出逻辑漏洞 (TLE, WA) 或边界情况。
	4. 用苏格拉底式提问引导学生自己修改。
	5. 如果学生提的问题太笨或者反复无常，请严厉批评。
	
	以下是学生的内容：
	`

	fullMessage := fmt.Sprintf("%s\n代码:\n%s\n\n问题: %s", systemPrompt, userCode, userQuestion)

	// 2. 构造请求体 (Gemini 格式)
	reqBodyData := GeminiRequest{
		Contents: []Content{
			{
				Role: "user",
				Parts: []Part{
					{Text: fullMessage},
				},
			},
		},
	}
	jsonData, _ := json.Marshal(reqBodyData)

	// 3. 拼接 URL (Google 的 Key 是放在 URL 参数里的)
	fullUrl := fmt.Sprintf("%s?key=%s", apiUrl, apiKey)

	// 4. 发送请求
	resp, err := http.Post(fullUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 5. 解析结果
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("API 报错 (%d): %s", resp.StatusCode, string(body))
	}

	var geminiResp GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return "", fmt.Errorf("JSON 解析失败: %v", err)
	}

	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		return geminiResp.Candidates[0].Content.Parts[0].Text, nil
	}

	return "", fmt.Errorf("Gemini 没有返回内容")
}