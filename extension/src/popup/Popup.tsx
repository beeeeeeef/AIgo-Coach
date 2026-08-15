import { useState } from 'react';
import { analyzeCode, AnalyzeRequest, AnalyzeResponse } from '../shared/api';
import './Popup.css';

export default function Popup() {
    const [code, setCode] = useState('');
    const [language, setLanguage] = useState('cpp');
    const [title, setTitle] = useState('');
    const [url, setUrl] = useState('');
    const [description, setDescription] = useState('');
    const [result, setResult] = useState<AnalyzeResponse | null>(null);
    const [loading, setLoading] = useState(false);

    const handleAnalyze = async () => {
        try {
            setLoading(true);
            const request: AnalyzeRequest = {
                site: 'leetcode-cn',
                title,
                url,
                description,
                language,
                code,
                mode: 'coach'
            };
            const response = await analyzeCode(request);
            setResult(response);
        } catch (error) {
            alert('分析失败：' + error);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="popup-container">
            <h2 className="popup-title">Algo-Coach</h2>

            <div className="form-group">
                <label>题目标题：</label>
                <input
                    type="text"
                    value={title}
                    onChange={(e) => setTitle(e.target.value)}
                    className="form-input"
                />
            </div>

            <div className="form-group">
                <label>题目 URL：</label>
                <input
                    type="text"
                    value={url}
                    onChange={(e) => setUrl(e.target.value)}
                    className="form-input"
                />
            </div>

            <div className="form-group">
                <label>题目描述：</label>
                <textarea
                    value={description}
                    onChange={(e) => setDescription(e.target.value)}
                    rows={3}
                    className="form-textarea"
                />
            </div>

            <div className="form-group">
                <label>语言：</label>
                <select
                    value={language}
                    onChange={(e) => setLanguage(e.target.value)}
                    className="form-select"
                >
                    <option value="cpp">C++</option>
                    <option value="python">Python</option>
                    <option value="java">Java</option>
                    <option value="go">Go</option>
                </select>
            </div>

            <div className="form-group">
                <label>代码：</label>
                <textarea
                    value={code}
                    onChange={(e) => setCode(e.target.value)}
                    rows={10}
                    placeholder="粘贴你的代码..."
                    className="form-textarea code-input"
                />
            </div>

            <button
                onClick={handleAnalyze}
                disabled={loading || !code || !title}
                className="analyze-button"
            >
                {loading ? '分析中...' : '分析代码'}
            </button>

            {result && result.result && (
                <div className="result-container">
                    <h3 className="result-title">分析结果</h3>

                    <div className="result-section">
                        <div className="result-section-title">总结：</div>
                        <p className="result-text">{result.data.summary}</p>
                    </div>

                    {result.data.problem.length > 0 && (
                        <div className="result-section">
                            <div className="result-section-title">问题：</div>
                            <ul className="result-list">
                                {result.data.problem.map((item, index) => (
                                    <li key={index}>{item}</li>
                                ))}
                            </ul>
                        </div>
                    )}

                    {result.data.hints.length > 0 && (
                        <div className="result-section">
                            <div className="result-section-title">提示：</div>
                            <ul className="result-list">
                                {result.data.hints.map((item, index) => (
                                    <li key={index}>{item}</li>
                                ))}
                            </ul>
                        </div>
                    )}

                    <div className="result-section">
                        <div className="result-section-title">复杂度分析：</div>
                        <p className="result-text">{result.data.complexity}</p>
                    </div>
                </div>
            )}

            {result && !result.result && (
                <div className="error-message">
                    错误：{result.error}
                </div>
            )}
        </div>
    );
}