package chat

// ProblemContext 一道题的上下文快照
type ProblemContext struct {
	Site        string `json:"site" binding:"required"`
	Title       string `json:"title" binding:"required"`
	URL         string `json:"url" binding:"required"`
	Description string `json:"description"`
	Language    string `json:"language" binding:"required"`
	Code        string `json:"code"`
}

// ChatRequest 多轮对话请求
type ChatRequest struct {
	ConversationID *int64          `json:"conversationId"`
	Message        string          `json:"message" binding:"required"`
	Context        *ProblemContext `json:"context"`
}

// ChatMessage 返回给前端的消息
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatData 成功时的业务数据（混合：主回复文本 + 可选结构）
type ChatData struct {
	ConversationID    int64         `json:"conversationId"`
	Reply             string        `json:"reply"`
	Hints             []string      `json:"hints,omitempty"`
	Issues            []string      `json:"issues,omitempty"`
	ComplexityComment string        `json:"complexityComment,omitempty"`
	Messages          []ChatMessage `json:"messages"`
}

// ChatResponse 统一响应
type ChatResponse struct {
	Result bool     `json:"result"`
	Error  string   `json:"error,omitempty"`
	Data   ChatData `json:"data"`
}
