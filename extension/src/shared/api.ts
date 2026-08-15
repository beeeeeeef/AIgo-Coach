//请求类型定义
export interface AnalyzeRequest {
    site: string;
    title: string;
    url: string;
    description: string;
    language: string;
    code: string;
    mode: string;
}
//响应类型定义
export interface AnalyzeResponse {
    result: boolean;
    error?: string;
    data: {
        summary: string;
        problem: string[];
        hints: string[];
        complexity: string;
    }
}

//调用后端分析接口
export async function analyzeCode(request: AnalyzeRequest): Promise<AnalyzeResponse> {
    const response = await fetch('http://localhost:8080/api/analyze', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(request),
    });
    if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
    }
    const data: AnalyzeResponse = await response.json();
    return data;
}

