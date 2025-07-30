package repositories

import (
	"context"
	"fmt"

	"github.com/company/microservice-template/internal/domain"
)

// NoOp implementations for when database is not available

// NoOp Conversation Repository
type noOpConversationRepository struct{}

func NewNoOpConversationRepository() domain.ConversationRepository {
	return &noOpConversationRepository{}
}

func (r *noOpConversationRepository) Create(ctx context.Context, conversation *domain.Conversation) error {
	return fmt.Errorf("database not available")
}

func (r *noOpConversationRepository) GetByID(ctx context.Context, id string) (*domain.Conversation, error) {
	return nil, fmt.Errorf("database not available")
}

func (r *noOpConversationRepository) GetByUserID(ctx context.Context, userID string, filters domain.ConversationFilters) ([]domain.Conversation, error) {
	return nil, fmt.Errorf("database not available")
}

func (r *noOpConversationRepository) Update(ctx context.Context, conversation *domain.Conversation) error {
	return fmt.Errorf("database not available")
}

func (r *noOpConversationRepository) Delete(ctx context.Context, id string) error {
	return fmt.Errorf("database not available")
}

// NoOp Message Repository
type noOpMessageRepository struct{}

func NewNoOpMessageRepository() domain.MessageRepository {
	return &noOpMessageRepository{}
}

func (r *noOpMessageRepository) Create(ctx context.Context, message *domain.Message) error {
	return fmt.Errorf("database not available")
}

func (r *noOpMessageRepository) GetByID(ctx context.Context, id string) (*domain.Message, error) {
	return nil, fmt.Errorf("database not available")
}

func (r *noOpMessageRepository) GetByConversationID(ctx context.Context, conversationID string, pagination domain.PaginationParams) ([]domain.Message, error) {
	return nil, fmt.Errorf("database not available")
}

func (r *noOpMessageRepository) Update(ctx context.Context, message *domain.Message) error {
	return fmt.Errorf("database not available")
}

func (r *noOpMessageRepository) Delete(ctx context.Context, id string) error {
	return fmt.Errorf("database not available")
}

// NoOp Attachment Repository
type noOpAttachmentRepository struct{}

func NewNoOpAttachmentRepository() domain.AttachmentRepository {
	return &noOpAttachmentRepository{}
}

func (r *noOpAttachmentRepository) Create(ctx context.Context, attachment *domain.Attachment) error {
	return fmt.Errorf("database not available")
}

func (r *noOpAttachmentRepository) GetByID(ctx context.Context, id string) (*domain.Attachment, error) {
	return nil, fmt.Errorf("database not available")
}

func (r *noOpAttachmentRepository) GetByMessageID(ctx context.Context, messageID string) ([]domain.Attachment, error) {
	return nil, fmt.Errorf("database not available")
}

func (r *noOpAttachmentRepository) Delete(ctx context.Context, id string) error {
	return fmt.Errorf("database not available")
}