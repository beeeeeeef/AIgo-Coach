// 题目上下文（新建会话必填）
export interface ProblemContext {
    site: string;
    title: string;
    url: string;
    description: string;
    language: string;
    code: string;
}

export interface ChatMessage {
    role: 'user' | 'assistant' | 'system';
    content: string;
}

export interface ChatRequest {
    conversationId?: number | null;
    message: string;
    context?: ProblemContext;
}

export interface ChatResponse {
    result: boolean;
    error?: string;
    data: {
        conversationId: number;
        reply: string;
        hints?: string[];
        issues?: string[];
        complexityComment?: string;
        messages: ChatMessage[];
    };
}

// 插件主路径：多轮对话
export async function chat(request: ChatRequest): Promise<ChatResponse> {
    const response = await fetch('http://localhost:8080/api/chat', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            conversationId: request.conversationId ?? null,
            message: request.message,
            context: request.context,
        }),
    });

    const data: ChatResponse = await response.json();
    if (!response.ok || !data.result) {
        throw new Error(data.error || `HTTP error! status: ${response.status}`);
    }
    return data;
}
