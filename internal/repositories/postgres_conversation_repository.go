package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/company/microservice-template/internal/domain"
	"github.com/company/microservice-template/pkg/logger"
)

type postgresConversationRepository struct {
	db     *sql.DB
	logger logger.Logger
}

func NewPostgresConversationRepository(db *sql.DB, logger logger.Logger) domain.ConversationRepository {
	return &postgresConversationRepository{
		db:     db,
		logger: logger,
	}
}

func (r *postgresConversationRepository) Create(ctx context.Context, conversation *domain.Conversation) error {
	query := `
		INSERT INTO conversations (id, user_id, channel, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	
	_, err := r.db.ExecContext(ctx, query,
		conversation.ID,
		conversation.UserID,
		conversation.Channel,
		conversation.Status,
		conversation.CreatedAt,
		conversation.UpdatedAt,
	)
	
	if err != nil {
		r.logger.Error("Failed to create conversation", err)
		return fmt.Errorf("failed to create conversation: %w", err)
	}
	
	return nil
}

func (r *postgresConversationRepository) GetByID(ctx context.Context, id string) (*domain.Conversation, error) {
	query := `
		SELECT id, user_id, channel, status, created_at, updated_at
		FROM conversations
		WHERE id = $1
	`
	
	var conversation domain.Conversation
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&conversation.ID,
		&conversation.UserID,
		&conversation.Channel,
		&conversation.Status,
		&conversation.CreatedAt,
		&conversation.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("conversation not found")
		}
		r.logger.Error("Failed to get conversation by ID", err)
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}
	
	return &conversation, nil
}

func (r *postgresConversationRepository) GetByUserID(ctx context.Context, userID string, filters domain.ConversationFilters) ([]domain.Conversation, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1
	
	// Base query
	query := `
		SELECT id, user_id, channel, status, created_at, updated_at
		FROM conversations
		WHERE user_id = $1
	`
	args = append(args, userID)
	argIndex++
	
	// Add filters
	if filters.Channel != "" {
		conditions = append(conditions, fmt.Sprintf("channel = $%d", argIndex))
		args = append(args, filters.Channel)
		argIndex++
	}
	
	if filters.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, filters.Status)
		argIndex++
	}
	
	// Add conditions to query
	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}
	
	// Add ordering and pagination
	query += " ORDER BY updated_at DESC"
	
	if filters.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filters.Limit)
		argIndex++
	}
	
	if filters.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, filters.Offset)
	}
	
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.Error("Failed to get conversations by user ID", err)
		return nil, fmt.Errorf("failed to get conversations: %w", err)
	}
	defer rows.Close()
	
	var conversations []domain.Conversation
	for rows.Next() {
		var conversation domain.Conversation
		err := rows.Scan(
			&conversation.ID,
			&conversation.UserID,
			&conversation.Channel,
			&conversation.Status,
			&conversation.CreatedAt,
			&conversation.UpdatedAt,
		)
		if err != nil {
			r.logger.Error("Failed to scan conversation row", err)
			continue
		}
		conversations = append(conversations, conversation)
	}
	
	if err = rows.Err(); err != nil {
		r.logger.Error("Error iterating conversation rows", err)
		return nil, fmt.Errorf("failed to iterate conversations: %w", err)
	}
	
	return conversations, nil
}

func (r *postgresConversationRepository) Update(ctx context.Context, conversation *domain.Conversation) error {
	query := `
		UPDATE conversations
		SET user_id = $2, channel = $3, status = $4, updated_at = $5
		WHERE id = $1
	`
	
	result, err := r.db.ExecContext(ctx, query,
		conversation.ID,
		conversation.UserID,
		conversation.Channel,
		conversation.Status,
		conversation.UpdatedAt,
	)
	
	if err != nil {
		r.logger.Error("Failed to update conversation", err)
		return fmt.Errorf("failed to update conversation: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("conversation not found")
	}
	
	return nil
}

func (r *postgresConversationRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM conversations WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Error("Failed to delete conversation", err)
		return fmt.Errorf("failed to delete conversation: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("conversation not found")
	}
	
	return nil
}