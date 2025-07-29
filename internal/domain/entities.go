package domain

import (
	"time"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// User representa un usuario del sistema
type User struct {
	ID        string    `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Name      string    `json:"name" db:"name"`
	Roles     []string  `json:"roles" db:"roles"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Messaging Service Entities

// ConversationStatus representa el estado de una conversación
type ConversationStatus string

const (
	ConversationStatusActive   ConversationStatus = "active"
	ConversationStatusClosed   ConversationStatus = "closed"
	ConversationStatusArchived ConversationStatus = "archived"
)

// Channel representa los canales de comunicación
type Channel string

const (
	ChannelWhatsApp  Channel = "whatsapp"
	ChannelWeb       Channel = "web"
	ChannelMessenger Channel = "messenger"
	ChannelInstagram Channel = "instagram"
)

// SenderType representa el tipo de remitente
type SenderType string

const (
	SenderTypeUser   SenderType = "user"
	SenderTypeBot    SenderType = "bot"
	SenderTypeSystem SenderType = "system"
)

// ContentType representa el tipo de contenido del mensaje
type ContentType string

const (
	ContentTypeText  ContentType = "text"
	ContentTypeImage ContentType = "image"
	ContentTypeVideo ContentType = "video"
	ContentTypeAudio ContentType = "audio"
	ContentTypeFile  ContentType = "file"
)

// AttachmentType representa el tipo de archivo adjunto
type AttachmentType string

const (
	AttachmentTypeImage AttachmentType = "image"
	AttachmentTypeVideo AttachmentType = "video"
	AttachmentTypeFile  AttachmentType = "file"
	AttachmentTypeAudio AttachmentType = "audio"
)

// JSONB type for PostgreSQL JSONB fields
type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = make(map[string]interface{})
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	
	return json.Unmarshal(bytes, j)
}

// Conversation representa una conversación
type Conversation struct {
	ID        string             `json:"id" db:"id"`
	UserID    string             `json:"user_id" db:"user_id"`
	Channel   Channel            `json:"channel" db:"channel"`
	Status    ConversationStatus `json:"status" db:"status"`
	CreatedAt time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" db:"updated_at"`
	Messages  []Message          `json:"messages,omitempty" db:"-"`
}

// Message representa un mensaje
type Message struct {
	ID             string      `json:"id" db:"id"`
	ConversationID string      `json:"conversation_id" db:"conversation_id"`
	SenderType     SenderType  `json:"sender_type" db:"sender_type"`
	SenderID       string      `json:"sender_id" db:"sender_id"`
	Content        string      `json:"content" db:"content"`
	ContentType    ContentType `json:"content_type" db:"content_type"`
	Metadata       JSONB       `json:"metadata" db:"metadata"`
	Timestamp      time.Time   `json:"timestamp" db:"timestamp"`
	Attachments    []Attachment `json:"attachments,omitempty" db:"-"`
}

// Attachment representa un archivo adjunto
type Attachment struct {
	ID        string         `json:"id" db:"id"`
	MessageID string         `json:"message_id" db:"message_id"`
	URL       string         `json:"url" db:"url"`
	Type      AttachmentType `json:"type" db:"type"`
	Size      int64          `json:"size" db:"size"`
	Filename  string         `json:"filename" db:"filename"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
}

// MessageEvent representa un evento de mensaje para pub/sub
type MessageEvent struct {
	Type           string      `json:"type"`
	ConversationID string      `json:"conversation_id"`
	Message        Message     `json:"message"`
	Timestamp      time.Time   `json:"timestamp"`
}

// AuditLog representa un registro de auditoría
type AuditLog struct {
	ID        string                 `json:"id" db:"id"`
	UserID    string                 `json:"user_id" db:"user_id"`
	Action    string                 `json:"action" db:"action"`
	Resource  string                 `json:"resource" db:"resource"`
	Details   map[string]interface{} `json:"details" db:"details"`
	IPAddress string                 `json:"ip_address" db:"ip_address"`
	UserAgent string                 `json:"user_agent" db:"user_agent"`
	CreatedAt time.Time              `json:"created_at" db:"created_at"`
}

// APIResponse estructura estándar para respuestas de API
type APIResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// HealthStatus representa el estado de salud del servicio
type HealthStatus struct {
	Status    string                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Uptime    string                 `json:"uptime"`
	Service   string                 `json:"service"`
	Version   string                 `json:"version"`
	Checks    map[string]interface{} `json:"checks,omitempty"`
}