package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/company/microservice-template/internal/domain"
	"github.com/company/microservice-template/pkg/logger"
)

type postgresMessageRepository struct {
	db     *sql.DB
	logger logger.Logger
}

func NewPostgresMessageRepository(db *sql.DB, logger logger.Logger) domain.MessageRepository {
	return &postgresMessageRepository{
		db:     db,
		logger: logger,
	}
}

func (r *postgresMessageRepository) Create(ctx context.Context, message *domain.Message) error {
	metadataJSON, err := json.Marshal(message.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	query := `
		INSERT INTO messages (id, conversation_id, sender_type, sender_id, content, content_type, metadata, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	
	_, err = r.db.ExecContext(ctx, query,
		message.ID,
		message.ConversationID,
		message.SenderType,
		message.SenderID,
		message.Content,
		message.ContentType,
		metadataJSON,
		message.Timestamp,
	)
	
	if err != nil {
		r.logger.Error("Failed to create message", err)
		return fmt.Errorf("failed to create message: %w", err)
	}
	
	return nil
}

func (r *postgresMessageRepository) GetByID(ctx context.Context, id string) (*domain.Message, error) {
	query := `
		SELECT id, conversation_id, sender_type, sender_id, content, content_type, metadata, timestamp
		FROM messages
		WHERE id = $1
	`
	
	var message domain.Message
	var metadataJSON []byte
	
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&message.ID,
		&message.ConversationID,
		&message.SenderType,
		&message.SenderID,
		&message.Content,
		&message.ContentType,
		&metadataJSON,
		&message.Timestamp,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("message not found")
		}
		r.logger.Error("Failed to get message by ID", err)
		return nil, fmt.Errorf("failed to get message: %w", err)
	}
	
	// Unmarshal metadata
	if len(metadataJSON) > 0 {
		if err := json.Unmarshal(metadataJSON, &message.Metadata); err != nil {
			r.logger.Error("Failed to unmarshal message metadata", err)
			message.Metadata = make(domain.JSONB)
		}
	} else {
		message.Metadata = make(domain.JSONB)
	}
	
	return &message, nil
}

func (r *postgresMessageRepository) GetByConversationID(ctx context.Context, conversationID string, pagination domain.PaginationParams) ([]domain.Message, error) {
	query := `
		SELECT id, conversation_id, sender_type, sender_id, content, content_type, metadata, timestamp
		FROM messages
		WHERE conversation_id = $1
		ORDER BY timestamp DESC
	`
	
	args := []interface{}{conversationID}
	argIndex := 2
	
	// Add pagination
	if pagination.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, pagination.Limit)
		argIndex++
	}
	
	if pagination.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, pagination.Offset)
	}
	
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.Error("Failed to get messages by conversation ID", err)
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}
	defer rows.Close()
	
	var messages []domain.Message
	for rows.Next() {
		var message domain.Message
		var metadataJSON []byte
		
		err := rows.Scan(
			&message.ID,
			&message.ConversationID,
			&message.SenderType,
			&message.SenderID,
			&message.Content,
			&message.ContentType,
			&metadataJSON,
			&message.Timestamp,
		)
		if err != nil {
			r.logger.Error("Failed to scan message row", err)
			continue
		}
		
		// Unmarshal metadata
		if len(metadataJSON) > 0 {
			if err := json.Unmarshal(metadataJSON, &message.Metadata); err != nil {
				r.logger.Error("Failed to unmarshal message metadata", err)
				message.Metadata = make(domain.JSONB)
			}
		} else {
			message.Metadata = make(domain.JSONB)
		}
		
		messages = append(messages, message)
	}
	
	if err = rows.Err(); err != nil {
		r.logger.Error("Error iterating message rows", err)
		return nil, fmt.Errorf("failed to iterate messages: %w", err)
	}
	
	return messages, nil
}

func (r *postgresMessageRepository) Update(ctx context.Context, message *domain.Message) error {
	metadataJSON, err := json.Marshal(message.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	query := `
		UPDATE messages
		SET conversation_id = $2, sender_type = $3, sender_id = $4, content = $5, content_type = $6, metadata = $7, timestamp = $8
		WHERE id = $1
	`
	
	result, err := r.db.ExecContext(ctx, query,
		message.ID,
		message.ConversationID,
		message.SenderType,
		message.SenderID,
		message.Content,
		message.ContentType,
		metadataJSON,
		message.Timestamp,
	)
	
	if err != nil {
		r.logger.Error("Failed to update message", err)
		return fmt.Errorf("failed to update message: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("message not found")
	}
	
	return nil
}

func (r *postgresMessageRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM messages WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Error("Failed to delete message", err)
		return fmt.Errorf("failed to delete message: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("message not found")
	}
	
	return nil
}