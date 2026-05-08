import { useState } from 'react';
import './App.css';
import { useMutation } from '@tanstack/react-query'; // 引入 TanStack
import axios from 'axios';
import { useChatStore } from './store'; // 引入 Zustand

interface ChatRequestPayload {
  code: string;
  question: string;
}

interface ChatResponseData {
  reply: string;
}

interface ApiErrorPayload {
  code: string;
  message: string;
}

interface ApiResponse<T> {
  success: boolean;
  data: T | null;
  error: ApiErrorPayload | null;
}

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080';

function App() {
  const [code, setCode] = useState('');
  const [question, setQuestion] = useState('');
  
  // 1. 从 Zustand 取出状态和方法
  const { messages, addMessage } = useChatStore();

  // 2. 使用 TanStack Query 定义请求逻辑
  const chatMutation = useMutation({
    mutationFn: async (payload: ChatRequestPayload) => {
      // 发送请求
      const res = await axios.post<ApiResponse<ChatResponseData>>(`${API_BASE_URL}/chat`, payload);
      return res.data;
    },
    onSuccess: (data) => {
      if (!data.success || !data.data) {
        alert(data.error?.message ?? '请求失败');
        return;
      }

      // 请求成功后，把 AI 的回复存进 Zustand
      addMessage('assistant', data.data.reply);
    },
    onError: (error) => {
      if (axios.isAxiosError<ApiResponse<null>>(error)) {
        alert('请求失败: ' + (error.response?.data?.error?.message ?? error.message));
        return;
      }

      alert('请求失败: ' + error.message);
    }
  });

  const handleSend = () => {
    if (!code || !question) return;

    // 先把用户的话上屏
    addMessage('user', question);
    
    // 触发请求
    chatMutation.mutate({ code, question });
    
    // 清空问题框
    setQuestion('');
  };
  return (
    <div className="container">
      <div className="left-panel">
        <h2>🧑‍💻 你的代码</h2>
        <textarea
          value={code}
          onChange={(e) => setCode(e.target.value)}
          placeholder="粘贴代码..."
        />
      </div>

      <div className="right-panel">
        <h2>🤖 Algo-Coach</h2>
        <div className="chat-box">
          {messages.map((msg, idx) => (
            <div key={idx} className={`message ${msg.role}`}>
              <strong>{msg.role === 'user' ? '我' : 'AI'}:</strong>
              <pre>{msg.content}</pre>
            </div>
          ))}
          
          {/* 使用 mutation 的 isPending 状态自动判断 loading */}
          {chatMutation.isPending && <div className="message assistant">思考中...</div>}
        </div>

        <div className="input-area">
          <input
            value={question}
            onChange={(e) => setQuestion(e.target.value)}
            onKeyDown={(e) => e.key === 'Enter' && handleSend()}
            placeholder="提出你的问题..."
          />
          <button onClick={handleSend} disabled={chatMutation.isPending}>
            {chatMutation.isPending ? '...' : '发送'}
          </button>
        </div>
      </div>
    </div>
  );
}

export default App;
