package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/company/microservice-template/internal/domain"
	"github.com/company/microservice-template/pkg/logger"
)

type postgresAttachmentRepository struct {
	db     *sql.DB
	logger logger.Logger
}

func NewPostgresAttachmentRepository(db *sql.DB, logger logger.Logger) domain.AttachmentRepository {
	return &postgresAttachmentRepository{
		db:     db,
		logger: logger,
	}
}

func (r *postgresAttachmentRepository) Create(ctx context.Context, attachment *domain.Attachment) error {
	query := `
		INSERT INTO attachments (id, message_id, url, type, size, filename, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	
	_, err := r.db.ExecContext(ctx, query,
		attachment.ID,
		attachment.MessageID,
		attachment.URL,
		attachment.Type,
		attachment.Size,
		attachment.Filename,
		attachment.CreatedAt,
	)
	
	if err != nil {
		r.logger.Error("Failed to create attachment", err)
		return fmt.Errorf("failed to create attachment: %w", err)
	}
	
	return nil
}

func (r *postgresAttachmentRepository) GetByID(ctx context.Context, id string) (*domain.Attachment, error) {
	query := `
		SELECT id, message_id, url, type, size, filename, created_at
		FROM attachments
		WHERE id = $1
	`
	
	var attachment domain.Attachment
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&attachment.ID,
		&attachment.MessageID,
		&attachment.URL,
		&attachment.Type,
		&attachment.Size,
		&attachment.Filename,
		&attachment.CreatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("attachment not found")
		}
		r.logger.Error("Failed to get attachment by ID", err)
		return nil, fmt.Errorf("failed to get attachment: %w", err)
	}
	
	return &attachment, nil
}

func (r *postgresAttachmentRepository) GetByMessageID(ctx context.Context, messageID string) ([]domain.Attachment, error) {
	query := `
		SELECT id, message_id, url, type, size, filename, created_at
		FROM attachments
		WHERE message_id = $1
		ORDER BY created_at ASC
	`
	
	rows, err := r.db.QueryContext(ctx, query, messageID)
	if err != nil {
		r.logger.Error("Failed to get attachments by message ID", err)
		return nil, fmt.Errorf("failed to get attachments: %w", err)
	}
	defer rows.Close()
	
	var attachments []domain.Attachment
	for rows.Next() {
		var attachment domain.Attachment
		err := rows.Scan(
			&attachment.ID,
			&attachment.MessageID,
			&attachment.URL,
			&attachment.Type,
			&attachment.Size,
			&attachment.Filename,
			&attachment.CreatedAt,
		)
		if err != nil {
			r.logger.Error("Failed to scan attachment row", err)
			continue
		}
		attachments = append(attachments, attachment)
	}
	
	if err = rows.Err(); err != nil {
		r.logger.Error("Error iterating attachment rows", err)
		return nil, fmt.Errorf("failed to iterate attachments: %w", err)
	}
	
	return attachments, nil
}

func (r *postgresAttachmentRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM attachments WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Error("Failed to delete attachment", err)
		return fmt.Errorf("failed to delete attachment: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("attachment not found")
	}
	
	return nil
}