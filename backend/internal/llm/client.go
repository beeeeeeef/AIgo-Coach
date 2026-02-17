package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// 数据结构
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

type GeminiResponse struct {
	Candidates []struct {
		Content Content `json:"content"`
	} `json:"candidates"`
}

// 核心函数
func ChatWithGemini(userCode string, userQuestion string) (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	apiUrl := os.Getenv("GEMINI_API_URL")

	apiUrl = strings.TrimSpace(apiUrl)
	fullUrl := fmt.Sprintf("%s?key=%s", apiUrl, apiKey)

	// 准备 Prompt
	systemPrompt := `你是一位严厉但循循善诱的 ACM 算法竞赛教练。
	学生会发给你代码和问题。
	1. 请不要直接给出完整代码。
	2. 先分析代码的时间/空间复杂度。
	3. 指出逻辑漏洞 (TLE, WA) 或边界情况。
	4. 用苏格拉底式提问引导学生自己修改。
	5. 如果学生提的问题太笨或者反复无常，请严厉批评。
	
	以下是学生的内容：
	`
	// 组合最终发给 AI 的文本
	finalText := fmt.Sprintf("%s\n代码:\n%s\n\n问题: %s", systemPrompt, userCode, userQuestion)

	reqBodyData := GeminiRequest{
		Contents: []Content{
			{
				Role: "user",
				Parts: []Part{
					{Text: finalText},
				},
			},
		},
	}
	jsonData, _ := json.Marshal(reqBodyData)

	fmt.Println("--------------------------------")
	fmt.Println("🚀 [Client] 正在请求:", fullUrl)
	fmt.Println("--------------------------------")

	// 使用自定义 Client 设置代理
	// 创建请求对象
	req, _ := http.NewRequest("POST", fullUrl, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	proxyStr := "http://127.0.0.1:7890"
	proxyURL, _ := url.Parse(proxyStr)

	client := &http.Client{
		Timeout: 30 * time.Second, // 设置超时时间，防止死等
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL), // 强制走代理
		},
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		// 这里的报错通常是网络不通
		return "", fmt.Errorf("网络请求发送失败 (请检查代理端口 %s): %v", proxyStr, err)
	}
	defer resp.Body.Close()

	// 解析 Body
	body, _ := io.ReadAll(resp.Body)

	// 更详细的错误日志
	if resp.StatusCode != 200 {
		fmt.Println("❌ Google API 报错详情:", string(body))
		return "", fmt.Errorf("Google API 状态码 %d", resp.StatusCode)
	}

	var geminiResp GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return "", fmt.Errorf("JSON 解析失败: %v", err)
	}

	// 安全获取内容
	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		return geminiResp.Candidates[0].Content.Parts[0].Text, nil
	}

	return "", fmt.Errorf("Gemini 返回了空内容 (可能是触发了安全拦截)")
}
