// extension/src/content/leetcode-cn.ts
// 职责：识别 leetcode.cn 题目页，抓取 title/url/description，响应 popup 消息
export type ExtractProblemRequest = {
    type: 'EXTRACT_PROBLEM'
};
export type ExtractProblemSuccess = {
    ok: true;
    data: {
        title: string;
        url: string;
        description: string;
    }
};

export type ExtractProblemFailure = {
    ok: false;
    error: string;
};

export type ExtractProblemResponse = ExtractProblemSuccess | ExtractProblemFailure;

// 监听来自 popup 的消息
function isLeetCodeProblemPage(href = location.href): boolean {
    try {
        const url = new URL(href);
        return (
            url.hostname === 'leetcode.cn' && url.pathname.includes('/problems/')
        );

    } catch {
        return false;
    }
}

function queryText(selectors: string[]): string {
    for (const selector of selectors) {
        const element = document.querySelector(selector);
        const text = element?.textContent?.trim();
        if (text) {
            return text;
        }
    }
    return '';
}

function queryDescription(selectors: string[]): string {
    for (const selector of selectors) {
        const element = document.querySelector(selector) as HTMLElement | null;
        const text = element?.textContent?.trim();
        if (text && text.length > 20) {
            return text;
        }
    }
    return '';
}

function extractProblem(): ExtractProblemResponse {
    if (!isLeetCodeProblemPage()) {
        return {
            ok: false,
            error: '当前页面不是 leetcode.cn 题目页'
        };
    }
    const title =
        queryText([
            '[data-cy="question-title"]',    // 官方测试标识
            'div[class*="text-title"]',       // 包含 text-title 的 class
            'h1',                              // 最后试试 h1
            '.question-title',
        ]) || document.title.replace(/\s*-\s*力扣.*$/, '').trim();
    // ↑ 如果都没找到，从 document.title 提取（去掉 "- 力扣"）
    const url = location.href;
    // 3. 提取描述
    const description = queryDescription([
        '[data-track-load="description_content"]',
        'div[class*="question-content"]',
        'div[class*="content__"]',
        '.elfjS',
        'div[data-key="description-content"]',
    ]);

    if (!title) {
        return {
            ok: false,
            error: '未能提取到题目标题'
        };
    }
    if (!description) {
        return {
            ok: false,
            error: '未找到题目描述，可稍后重试或检查是否在题目描述页',
        };
    }

    // 5. 返回成功结果
    return {
        ok: true,
        data: { title, url, description },
    };
}

chrome.runtime.onMessage.addListener((message, _sender, sendResponse) => {
    // 1. 检查消息类型
    if (!message || message.type !== 'EXTRACT_PROBLEM') {
        return;  // 忽略不相关的消息
    }

    try {
        // 2. 执行提取
        const result = extractProblem();
        sendResponse(result);
    } catch (error) {
        // 3. 异常处理
        sendResponse({
            ok: false,
            error: error instanceof Error ? error.message : '抓取失败',
        } satisfies ExtractProblemFailure);
    }

    // 4. 保持消息通道
    return true;
});