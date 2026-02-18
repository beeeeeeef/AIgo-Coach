import { create } from 'zustand';

// 定义消息类型
interface Message {
  role: 'user' | 'assistant';
  content: string;
}

// 定义 Store 的形状
interface ChatStore {
  messages: Message[];
  addMessage: (role: 'user' | 'assistant', content: string) => void;
  clearMessages: () => void;
}

// 创建 Store
export const useChatStore = create<ChatStore>((set) => ({
  messages: [],
  addMessage: (role, content) => set((state) => ({ 
    messages: [...state.messages, { role, content }] 
  })),
  clearMessages: () => set({ messages: [] }),
}));