package chat

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// Conversation 会话实体
type Conversation struct {
	ID          int64
	UserID      sql.NullInt64
	Site        string
	Title       string
	URL         string
	Description string
	Language    string
	Code        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Message 消息实体
type Message struct {
	ID             int64
	ConversationID int64
	Role           string
	Content        string
	Meta           json.RawMessage
	CreatedAt      time.Time
}

// Repository 对话持久化
type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateConversation(ctx context.Context, c ProblemContext) (*Conversation, error) {
	const q = `
		INSERT INTO conversations (site, title, url, description, language, code)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, site, title, url, description, language, code, created_at, updated_at
	`
	row := r.db.QueryRowContext(ctx, q, c.Site, c.Title, c.URL, c.Description, c.Language, c.Code)
	return scanConversation(row)
}

func (r *Repository) GetConversation(ctx context.Context, id int64) (*Conversation, error) {
	const q = `
		SELECT id, user_id, site, title, url, description, language, code, created_at, updated_at
		FROM conversations
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, q, id)
	conv, err := scanConversation(row)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("会话不存在: %d", id)
	}
	return conv, err
}

func (r *Repository) UpdateCode(ctx context.Context, id int64, code string) error {
	const q = `
		UPDATE conversations
		SET code = $2, updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, q, id, code)
	return err
}

func (r *Repository) TouchConversation(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `UPDATE conversations SET updated_at = NOW() WHERE id = $1`, id)
	return err
}

func (r *Repository) AddMessage(ctx context.Context, conversationID int64, role, content string, meta any) (*Message, error) {
	var metaJSON []byte
	var err error
	if meta != nil {
		metaJSON, err = json.Marshal(meta)
		if err != nil {
			return nil, fmt.Errorf("序列化 meta 失败: %w", err)
		}
	}

	const q = `
		INSERT INTO messages (conversation_id, role, content, meta)
		VALUES ($1, $2, $3, $4)
		RETURNING id, conversation_id, role, content, meta, created_at
	`
	var m Message
	var raw []byte
	err = r.db.QueryRowContext(ctx, q, conversationID, role, content, nullableJSON(metaJSON)).Scan(
		&m.ID, &m.ConversationID, &m.Role, &m.Content, &raw, &m.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	if len(raw) > 0 {
		m.Meta = json.RawMessage(raw)
	}
	return &m, nil
}

func (r *Repository) ListMessages(ctx context.Context, conversationID int64) ([]Message, error) {
	const q = `
		SELECT id, conversation_id, role, content, meta, created_at
		FROM messages
		WHERE conversation_id = $1
		ORDER BY created_at ASC, id ASC
	`
	rows, err := r.db.QueryContext(ctx, q, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Message
	for rows.Next() {
		var m Message
		var raw []byte
		if err := rows.Scan(&m.ID, &m.ConversationID, &m.Role, &m.Content, &raw, &m.CreatedAt); err != nil {
			return nil, err
		}
		if len(raw) > 0 {
			m.Meta = json.RawMessage(raw)
		}
		list = append(list, m)
	}
	return list, rows.Err()
}

type scannable interface {
	Scan(dest ...any) error
}

func scanConversation(row scannable) (*Conversation, error) {
	var c Conversation
	err := row.Scan(
		&c.ID, &c.UserID, &c.Site, &c.Title, &c.URL, &c.Description, &c.Language, &c.Code, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func nullableJSON(b []byte) any {
	if len(b) == 0 {
		return nil
	}
	return b
}
