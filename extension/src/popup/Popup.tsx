import { useState } from 'react';
import { chat, ChatMessage, ProblemContext } from '../shared/api';

type ExtractProblemResponse =
    | {
        ok: true;
        data: {
            title: string;
            url: string;
            description: string;
            code: string;
            language: string;
        };
    }
    | {
        ok: false;
        error: string;
    };

const fieldClass =
    'w-full rounded-md border border-slate-200 bg-white px-2.5 py-2 text-[13px] text-slate-800 outline-none transition focus:border-blue-500 focus:ring-2 focus:ring-blue-100';
const labelClass =
    'mb-1.5 block text-[12px] font-semibold uppercase tracking-wide text-slate-500';
const cardClass =
    'mb-2.5 rounded-md bg-white p-2.5 shadow-sm shadow-slate-200/80';

function Spinner() {
    return (
        <span className="mr-2 inline-block h-3.5 w-3.5 animate-spin rounded-full border-2 border-white/30 border-t-white align-middle" />
    );
}

export default function Popup() {
    const [code, setCode] = useState('');
    const [language, setLanguage] = useState('cpp');
    const [title, setTitle] = useState('');
    const [url, setUrl] = useState('');
    const [description, setDescription] = useState('');
    const [conversationId, setConversationId] = useState<number | null>(null);
    const [messages, setMessages] = useState<ChatMessage[]>([]);
    const [input, setInput] = useState('');
    const [loading, setLoading] = useState(false);
    const [extracting, setExtracting] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [showContext, setShowContext] = useState(true);
    const [extraHints, setExtraHints] = useState<string[]>([]);
    const [extraIssues, setExtraIssues] = useState<string[]>([]);
    const [extraComplexity, setExtraComplexity] = useState('');
    const [codeExtracted, setCodeExtracted] = useState(false);

    const buildContext = (): ProblemContext => ({
        site: 'leetcode-cn',
        title,
        url,
        description,
        language,
        code,
    });

    const handleExtract = async () => {
        try {
            setExtracting(true);
            setError(null);

            const [tab] = await chrome.tabs.query({
                active: true,
                currentWindow: true,
            });

            if (!tab?.id) {
                setError('未找到当前标签页');
                return;
            }

            if (!tab.url || !tab.url.includes('leetcode.cn/problems/')) {
                setError('请先打开 LeetCode 中文站题目页');
                return;
            }

            const response = await chrome.tabs.sendMessage(tab.id, {
                type: 'EXTRACT_PROBLEM',
            }) as ExtractProblemResponse;

            if (!response?.ok) {
                setError('抓取失败：' + (response?.error || '未知错误'));
                return;
            }

            setTitle(response.data.title);
            setUrl(response.data.url);
            setDescription(response.data.description);
            setCode(response.data.code);
            setLanguage(response.data.language);
            setCodeExtracted(response.data.code.length > 0);
            setConversationId(null);
            setMessages([]);
            setExtraHints([]);
            setExtraIssues([]);
            setExtraComplexity('');
            setShowContext(true);
        } catch (err) {
            setError(
                '抓取失败：' +
                (err instanceof Error ? err.message : String(err)) +
                '。请确认已打开题目页，并已重新加载插件后刷新页面。'
            );
        } finally {
            setExtracting(false);
        }
    };

    const sendMessage = async (text: string) => {
        const content = text.trim();
        if (!content || loading) return;

        if (!title) {
            setError('请先抓取或填写题目标题');
            return;
        }

        try {
            setLoading(true);
            setError(null);
            setInput('');

            const response = await chat({
                conversationId,
                message: content,
                context: buildContext(),
            });

            setConversationId(response.data.conversationId);
            setMessages(response.data.messages);
            setExtraHints(response.data.hints || []);
            setExtraIssues(response.data.issues || []);
            setExtraComplexity(response.data.complexityComment || '');
            setShowContext(false);
        } catch (err) {
            const errorMsg = err instanceof Error ? err.message : String(err);
            if (errorMsg.includes('Failed to fetch') || errorMsg.includes('NetworkError')) {
                setError('无法连接后端。请确认 PostgreSQL 与后端已启动（localhost:8080）');
            } else {
                setError(errorMsg);
            }
        } finally {
            setLoading(false);
        }
    };

    const handleSend = async () => {
        await sendMessage(input);
    };

    const handleQuickAnalyze = async () => {
        await sendMessage('请先看看我的代码，指出主要问题和下一步思考方向。');
    };

    const handleNewChat = () => {
        setConversationId(null);
        setMessages([]);
        setExtraHints([]);
        setExtraIssues([]);
        setExtraComplexity('');
        setError(null);
        setShowContext(true);
    };

    return (
        <div className="h-[640px] w-[520px] overflow-y-auto bg-slate-100 p-4 font-sans text-slate-900">
            <h2 className="mb-3 border-b-2 border-blue-500 pb-3 text-center text-[22px] font-bold text-slate-900">
                Algo-Coach
            </h2>

            <div className="mb-3 flex gap-2">
                <button
                    onClick={handleExtract}
                    disabled={extracting}
                    className="flex-1 rounded-md bg-gradient-to-br from-green-600 to-emerald-500 px-3 py-2.5 text-sm font-semibold text-white shadow-md shadow-green-600/30 transition hover:from-green-700 hover:to-emerald-600 disabled:cursor-not-allowed disabled:bg-none disabled:bg-slate-300 disabled:shadow-none"
                >
                    {extracting && <Spinner />}
                    {extracting ? '抓取中...' : '🎯 抓取当前题目'}
                </button>
                <button
                    onClick={handleNewChat}
                    className="rounded-md bg-slate-200 px-3 py-2.5 text-sm font-semibold text-slate-700 transition hover:bg-slate-300"
                >
                    新对话
                </button>
            </div>

            {error && (
                <div className="mb-3 rounded-md border-l-4 border-amber-400 bg-amber-50 px-3 py-3 text-[13px] leading-relaxed text-amber-800">
                    ⚠️ {error}
                </div>
            )}

            {title && (
                <div className="mb-3 flex items-center justify-between gap-2 rounded-md bg-white px-3 py-2.5 shadow-sm shadow-slate-200/80">
                    <div className="text-sm font-semibold text-slate-800">{title}</div>
                    <button
                        className="shrink-0 border-0 bg-transparent text-xs text-blue-500 hover:text-blue-600"
                        onClick={() => setShowContext((v) => !v)}
                    >
                        {showContext ? '收起题目/代码' : '展开题目/代码'}
                    </button>
                </div>
            )}

            {showContext && (
                <>
                    <div className={cardClass}>
                        <label className={labelClass}>题目标题</label>
                        <input
                            type="text"
                            value={title}
                            onChange={(e) => setTitle(e.target.value)}
                            className={fieldClass}
                        />
                    </div>

                    <div className={cardClass}>
                        <label className={labelClass}>题目 URL</label>
                        <input
                            type="text"
                            value={url}
                            onChange={(e) => setUrl(e.target.value)}
                            className={fieldClass}
                        />
                    </div>

                    <div className={cardClass}>
                        <label className={labelClass}>题目描述</label>
                        <textarea
                            value={description}
                            onChange={(e) => setDescription(e.target.value)}
                            rows={3}
                            className={fieldClass}
                        />
                    </div>

                    {codeExtracted ? (
                        <div className={cardClass}>
                            <label className={labelClass}>已抓取代码（{language}）</label>
                            <div className="max-h-[200px] overflow-y-auto whitespace-pre-wrap break-all rounded border border-zinc-700 bg-zinc-900 p-3 font-mono text-[13px] text-zinc-300">
                                {code || '（空）'}
                            </div>
                            <button
                                onClick={() => setCodeExtracted(false)}
                                className="mt-1 border-0 bg-transparent text-xs text-blue-500 hover:text-blue-600"
                            >
                                切换到手动编辑
                            </button>
                        </div>
                    ) : (
                        <>
                            <div className={cardClass}>
                                <label className={labelClass}>语言</label>
                                <select
                                    value={language}
                                    onChange={(e) => setLanguage(e.target.value)}
                                    className={fieldClass}
                                >
                                    <option value="cpp">C++</option>
                                    <option value="python">Python</option>
                                    <option value="java">Java</option>
                                    <option value="go">Go</option>
                                    <option value="javascript">JavaScript</option>
                                    <option value="rust">Rust</option>
                                    <option value="csharp">C#</option>
                                </select>
                            </div>

                            <div className={cardClass}>
                                <label className={labelClass}>代码</label>
                                <textarea
                                    value={code}
                                    onChange={(e) => setCode(e.target.value)}
                                    rows={8}
                                    placeholder="粘贴你的代码，或点击抓取自动获取..."
                                    className={`${fieldClass} font-mono`}
                                />
                            </div>
                        </>
                    )}

                    <button
                        onClick={handleQuickAnalyze}
                        disabled={loading || !title}
                        className="mb-3 mt-1 w-full rounded-md bg-gradient-to-br from-blue-500 to-blue-700 px-3 py-3 text-sm font-semibold text-white shadow-md shadow-blue-500/30 transition hover:from-blue-600 hover:to-blue-800 disabled:cursor-not-allowed disabled:bg-none disabled:bg-slate-300 disabled:shadow-none"
                    >
                        {loading && <Spinner />}
                        {loading ? 'AI 思考中...' : '🤖 先分析我的代码'}
                    </button>
                </>
            )}

            <div className="rounded-lg bg-white p-3 shadow-md shadow-slate-200/80">
                <div className="mb-3 flex max-h-[260px] flex-col gap-2.5 overflow-y-auto">
                    {messages.length === 0 && (
                        <div className="px-2 py-6 text-center text-[13px] text-slate-400">
                            抓取题目后，开始多轮对话吧。
                        </div>
                    )}
                    {messages.map((msg, index) => (
                        <div
                            key={index}
                            className={`max-w-[90%] rounded-[10px] px-3 py-2.5 text-[13px] leading-relaxed ${msg.role === 'user'
                                    ? 'self-end bg-blue-500 text-white'
                                    : 'self-start bg-slate-100 text-slate-800'
                                }`}
                        >
                            <div className="mb-1 text-[11px] font-semibold opacity-80">
                                {msg.role === 'user' ? '你' : '教练'}
                            </div>
                            <div className="whitespace-pre-wrap break-words">{msg.content}</div>
                        </div>
                    ))}
                </div>

                {(extraIssues.length > 0 || extraHints.length > 0 || extraComplexity) && (
                    <div className="mb-3 border-t border-dashed border-slate-200 pt-2">
                        {extraIssues.length > 0 && (
                            <div className="mb-2.5 rounded-md bg-slate-50 p-2">
                                <div className="mb-1.5 text-[12px] font-bold text-blue-500">⚠️ 问题</div>
                                <ul className="m-0 list-disc pl-4 text-[13px] leading-relaxed text-slate-700">
                                    {extraIssues.map((item, i) => (
                                        <li key={i}>{item}</li>
                                    ))}
                                </ul>
                            </div>
                        )}
                        {extraHints.length > 0 && (
                            <div className="mb-2.5 rounded-md bg-slate-50 p-2">
                                <div className="mb-1.5 text-[12px] font-bold text-blue-500">💡 提示</div>
                                <ul className="m-0 list-disc pl-4 text-[13px] leading-relaxed text-slate-700">
                                    {extraHints.map((item, i) => (
                                        <li key={i}>{item}</li>
                                    ))}
                                </ul>
                            </div>
                        )}
                        {extraComplexity && (
                            <div className="mb-2.5 rounded-md bg-slate-50 p-2">
                                <div className="mb-1.5 text-[12px] font-bold text-blue-500">⚡ 复杂度</div>
                                <p className="m-0 text-[13px] leading-relaxed text-slate-700">{extraComplexity}</p>
                            </div>
                        )}
                    </div>
                )}

                <div className="flex items-end gap-2">
                    <textarea
                        value={input}
                        onChange={(e) => setInput(e.target.value)}
                        rows={2}
                        placeholder="继续追问，例如：哈希表应该存什么？"
                        className={`${fieldClass} min-h-[52px] resize-y font-mono`}
                        onKeyDown={(e) => {
                            if (e.key === 'Enter' && !e.shiftKey) {
                                e.preventDefault();
                                void handleSend();
                            }
                        }}
                    />
                    <button
                        onClick={handleSend}
                        disabled={loading || !input.trim() || !title}
                        className="min-w-16 rounded-md bg-blue-500 px-3.5 py-3 text-sm font-semibold text-white transition hover:bg-blue-600 disabled:cursor-not-allowed disabled:bg-slate-300"
                    >
                        {loading ? '...' : '发送'}
                    </button>
                </div>
            </div>
        </div>
    );
}
