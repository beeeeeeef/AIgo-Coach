package chat

import (
	"context"
	"fmt"

	"aigo-coach/backend/internal/ai"
)

const maxHistoryMessages = 20

// Service 多轮对话业务
type Service struct {
	repo     *Repository
	aiClient ai.Client
}

func NewService(repo *Repository, aiClient ai.Client) *Service {
	return &Service{repo: repo, aiClient: aiClient}
}

func (s *Service) Chat(ctx context.Context, req ChatRequest) (ChatData, error) {
	var conv *Conversation
	var err error

	if req.ConversationID == nil {
		if req.Context == nil {
			return ChatData{}, fmt.Errorf("新建会话时必须提供 context")
		}
		conv, err = s.repo.CreateConversation(ctx, *req.Context)
		if err != nil {
			return ChatData{}, fmt.Errorf("创建会话失败: %w", err)
		}
		// 把题目卡写入 system，后续续聊可还原上下文
		if _, err := s.repo.AddMessage(ctx, conv.ID, "system", BuildProblemCard(*req.Context), nil); err != nil {
			return ChatData{}, fmt.Errorf("写入题目上下文失败: %w", err)
		}
	} else {
		conv, err = s.repo.GetConversation(ctx, *req.ConversationID)
		if err != nil {
			return ChatData{}, err
		}
		// 允许中途更新代码
		if req.Context != nil && req.Context.Code != "" && req.Context.Code != conv.Code {
			if err := s.repo.UpdateCode(ctx, conv.ID, req.Context.Code); err != nil {
				return ChatData{}, fmt.Errorf("更新代码失败: %w", err)
			}
			lang := conv.Language
			if req.Context.Language != "" {
				lang = req.Context.Language
			}
			note := "学生更新了代码：\n```" + lang + "\n" + req.Context.Code + "\n```"
			if _, err := s.repo.AddMessage(ctx, conv.ID, "system", note, nil); err != nil {
				return ChatData{}, fmt.Errorf("写入代码更新提示失败: %w", err)
			}
			conv.Code = req.Context.Code
		}
	}

	if _, err := s.repo.AddMessage(ctx, conv.ID, "user", req.Message, nil); err != nil {
		return ChatData{}, fmt.Errorf("保存用户消息失败: %w", err)
	}

	history, err := s.repo.ListMessages(ctx, conv.ID)
	if err != nil {
		return ChatData{}, fmt.Errorf("读取历史失败: %w", err)
	}

	aiMessages := buildAIMessages(history)
	config := ai.ChatConfig{
		Model:       "deepseek-chat",
		Temperature: 0.7,
		MaxTokens:   2000,
	}
	raw, err := s.aiClient.Chat(ctx, aiMessages, config)
	if err != nil {
		return ChatData{}, fmt.Errorf("AI 调用失败: %w", err)
	}

	reply, meta := SplitAssistantOutput(raw)
	if reply == "" {
		reply = raw
	}

	if _, err := s.repo.AddMessage(ctx, conv.ID, "assistant", reply, meta); err != nil {
		return ChatData{}, fmt.Errorf("保存助手消息失败: %w", err)
	}
	_ = s.repo.TouchConversation(ctx, conv.ID)

	all, err := s.repo.ListMessages(ctx, conv.ID)
	if err != nil {
		return ChatData{}, err
	}

	data := ChatData{
		ConversationID: conv.ID,
		Reply:          reply,
		Messages:       toPublicMessages(all),
	}
	if meta != nil {
		data.Hints = meta.Hints
		data.Issues = meta.Issues
		data.ComplexityComment = meta.ComplexityComment
	}
	return data, nil
}

func buildAIMessages(history []Message) []ai.Message {
	out := []ai.Message{{Role: "system", Content: ChatCoachSystemPrompt}}

	start := 0
	if len(history) > maxHistoryMessages {
		start = len(history) - maxHistoryMessages
	}
	for _, m := range history[start:] {
		role := m.Role
		if role == "system" {
			// 题目卡等系统上下文，继续以 system 送给模型
			out = append(out, ai.Message{Role: "system", Content: m.Content})
			continue
		}
		out = append(out, ai.Message{Role: role, Content: m.Content})
	}
	return out
}

func toPublicMessages(list []Message) []ChatMessage {
	out := make([]ChatMessage, 0, len(list))
	for _, m := range list {
		if m.Role == "system" {
			continue // 前端气泡不展示 system
		}
		out = append(out, ChatMessage{Role: m.Role, Content: m.Content})
	}
	return out
}
