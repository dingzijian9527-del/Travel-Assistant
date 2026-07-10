package mysql

import (
	"context"
	"database/sql"
	"time"

	aiagent "github.com/dingzijian9527-del/Travel-Assistant/kitex_gen/ai_agent"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/config"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/mysqlx"
	"github.com/dingzijian9527-del/Travel-Assistant/pkg/repository"
)

type AIAgentMessageRepository struct {
	db              *sql.DB
	maxHistoryItems int
}

func NewAIAgentMessageRepository(cfg config.MySQLConfig) (*AIAgentMessageRepository, error) {
	db, err := mysqlx.New(cfg)
	if err != nil {
		return nil, err
	}
	repo := &AIAgentMessageRepository{db: db, maxHistoryItems: 20}
	if err := repo.ensureSchema(context.Background()); err != nil {
		_ = db.Close()
		return nil, err
	}
	return repo, nil
}

func (r *AIAgentMessageRepository) ensureSchema(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS ai_agent_messages (
    id BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    user_id BIGINT(20) UNSIGNED NOT NULL,
    role VARCHAR(32) NOT NULL DEFAULT '',
    content MEDIUMTEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY idx_ai_agent_messages_user_created (user_id, created_at, id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)
	return err
}

func (r *AIAgentMessageRepository) AppendMessages(ctx context.Context, userID int64, messages ...*repository.ChatMessage) {
	if userID <= 0 {
		return
	}
	for _, message := range messages {
		if message == nil {
			continue
		}
		_, _ = r.db.ExecContext(ctx, `
INSERT INTO ai_agent_messages(user_id, role, content, created_at)
VALUES(?, ?, ?, ?)`, userID, message.Role.String(), message.Content, time.Now())
	}
	r.trimHistory(ctx, userID)
}

func (r *AIAgentMessageRepository) History(ctx context.Context, userID int64) []repository.ChatMessage {
	if userID <= 0 {
		return nil
	}
	limit := r.maxHistoryItems
	if limit <= 0 {
		limit = 20
	}
	rows, err := r.db.QueryContext(ctx, `
SELECT role, content
FROM ai_agent_messages
WHERE user_id = ?
ORDER BY created_at DESC, id DESC
LIMIT ?`, userID, limit)
	if err != nil {
		return nil
	}
	defer rows.Close()
	items := make([]repository.ChatMessage, 0, limit)
	for rows.Next() {
		var role string
		var content string
		if err := rows.Scan(&role, &content); err != nil {
			return nil
		}
		items = append(items, repository.ChatMessage{Role: parseChatRole(role), Content: content})
	}
	for left, right := 0, len(items)-1; left < right; left, right = left+1, right-1 {
		items[left], items[right] = items[right], items[left]
	}
	return items
}

func (r *AIAgentMessageRepository) trimHistory(ctx context.Context, userID int64) {
	limit := r.maxHistoryItems
	if limit <= 0 {
		limit = 20
	}
	_, _ = r.db.ExecContext(ctx, `
DELETE FROM ai_agent_messages
WHERE user_id = ?
  AND id NOT IN (
      SELECT id FROM (
          SELECT id
          FROM ai_agent_messages
          WHERE user_id = ?
          ORDER BY created_at DESC, id DESC
          LIMIT ?
      ) keep_messages
  )`, userID, userID, limit)
}

func parseChatRole(role string) aiagent.ChatRole {
	switch role {
	case aiagent.ChatRole_ASSISTANT.String():
		return aiagent.ChatRole_ASSISTANT
	case aiagent.ChatRole_USER.String():
		return aiagent.ChatRole_USER
	default:
		return aiagent.ChatRole_USER
	}
}
