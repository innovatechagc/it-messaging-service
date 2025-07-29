package domain

import (
	"context"
)

// Messaging repositories

// ConversationRepository define las operaciones para conversaciones
type ConversationRepository interface {
	Create(ctx context.Context, conversation *Conversation) error
	GetByID(ctx context.Context, id string) (*Conversation, error)
	GetByUserID(ctx context.Context, userID string, filters ConversationFilters) ([]Conversation, error)
	Update(ctx context.Context, conversation *Conversation) error
	Delete(ctx context.Context, id string) error
}

// MessageRepository define las operaciones para mensajes
type MessageRepository interface {
	Create(ctx context.Context, message *Message) error
	GetByID(ctx context.Context, id string) (*Message, error)
	GetByConversationID(ctx context.Context, conversationID string, pagination PaginationParams) ([]Message, error)
	Update(ctx context.Context, message *Message) error
	Delete(ctx context.Context, id string) error
}

// AttachmentRepository define las operaciones para archivos adjuntos
type AttachmentRepository interface {
	Create(ctx context.Context, attachment *Attachment) error
	GetByID(ctx context.Context, id string) (*Attachment, error)
	GetByMessageID(ctx context.Context, messageID string) ([]Attachment, error)
	Delete(ctx context.Context, id string) error
}

// ConversationFilters para filtrar conversaciones
type ConversationFilters struct {
	Channel Channel
	Status  ConversationStatus
	Limit   int
	Offset  int
}

// PaginationParams para paginación
type PaginationParams struct {
	Limit  int
	Offset int
	SortBy string
	Order  string
}

// UserRepository define las operaciones de persistencia para usuarios
type UserRepository interface {
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*User, error)
}

// AuditRepository define las operaciones de persistencia para auditoría
type AuditRepository interface {
	Create(ctx context.Context, log *AuditLog) error
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*AuditLog, error)
	GetByAction(ctx context.Context, action string, limit, offset int) ([]*AuditLog, error)
}

// HealthRepository define las operaciones para health checks
type HealthRepository interface {
	CheckDatabase(ctx context.Context) error
	CheckExternalServices(ctx context.Context) map[string]error
}