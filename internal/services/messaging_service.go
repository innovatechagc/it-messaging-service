package services

import (
	"context"
	"fmt"
	"time"

	"github.com/company/microservice-template/internal/domain"
	"github.com/company/microservice-template/pkg/logger"
	"github.com/google/uuid"
)

type MessagingService interface {
	// Conversations
	CreateConversation(ctx context.Context, userID string, channel domain.Channel) (*domain.Conversation, error)
	GetConversation(ctx context.Context, id string, userID string) (*domain.Conversation, error)
	GetConversations(ctx context.Context, userID string, filters domain.ConversationFilters) ([]domain.Conversation, error)
	UpdateConversationStatus(ctx context.Context, id string, status domain.ConversationStatus, userID string) error
	
	// Messages
	SendMessage(ctx context.Context, req SendMessageRequest) (*domain.Message, error)
	GetMessages(ctx context.Context, conversationID string, userID string, pagination domain.PaginationParams) ([]domain.Message, error)
	GetMessage(ctx context.Context, messageID string, userID string) (*domain.Message, error)
	
	// Attachments
	CreateAttachment(ctx context.Context, messageID string, req CreateAttachmentRequest) (*domain.Attachment, error)
	GetAttachment(ctx context.Context, attachmentID string, userID string) (*domain.Attachment, error)
}

type messagingService struct {
	conversationRepo domain.ConversationRepository
	messageRepo      domain.MessageRepository
	attachmentRepo   domain.AttachmentRepository
	eventPublisher   EventPublisher
	cacheService     CacheService
	logger           logger.Logger
}

type SendMessageRequest struct {
	ConversationID string                 `json:"conversation_id" binding:"required"`
	SenderType     domain.SenderType      `json:"sender_type" binding:"required"`
	SenderID       string                 `json:"sender_id" binding:"required"`
	Content        string                 `json:"content" binding:"required"`
	ContentType    domain.ContentType     `json:"content_type" binding:"required"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

type CreateAttachmentRequest struct {
	URL      string                `json:"url" binding:"required"`
	Type     domain.AttachmentType `json:"type" binding:"required"`
	Size     int64                 `json:"size" binding:"required"`
	Filename string                `json:"filename" binding:"required"`
}

func NewMessagingService(
	conversationRepo domain.ConversationRepository,
	messageRepo domain.MessageRepository,
	attachmentRepo domain.AttachmentRepository,
	eventPublisher EventPublisher,
	cacheService CacheService,
	logger logger.Logger,
) MessagingService {
	return &messagingService{
		conversationRepo: conversationRepo,
		messageRepo:      messageRepo,
		attachmentRepo:   attachmentRepo,
		eventPublisher:   eventPublisher,
		cacheService:     cacheService,
		logger:           logger,
	}
}

func (s *messagingService) CreateConversation(ctx context.Context, userID string, channel domain.Channel) (*domain.Conversation, error) {
	conversation := &domain.Conversation{
		ID:        uuid.New().String(),
		UserID:    userID,
		Channel:   channel,
		Status:    domain.ConversationStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.conversationRepo.Create(ctx, conversation); err != nil {
		s.logger.Error("Failed to create conversation", err)
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}

	s.logger.Info("Conversation created", map[string]interface{}{
		"conversation_id": conversation.ID,
		"user_id":         userID,
		"channel":         channel,
	})

	return conversation, nil
}

func (s *messagingService) GetConversation(ctx context.Context, id string, userID string) (*domain.Conversation, error) {
	// Check cache first
	if s.cacheService != nil {
		if cached, err := s.cacheService.GetConversation(ctx, id); err == nil && cached != nil {
			// Verify user ownership
			if cached.UserID == userID {
				return cached, nil
			}
		}
	}

	conversation, err := s.conversationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	// Verify user ownership
	if conversation.UserID != userID {
		return nil, fmt.Errorf("conversation not found or access denied")
	}

	// Cache the result
	if s.cacheService != nil {
		_ = s.cacheService.SetConversation(ctx, conversation)
	}

	return conversation, nil
}

func (s *messagingService) GetConversations(ctx context.Context, userID string, filters domain.ConversationFilters) ([]domain.Conversation, error) {
	conversations, err := s.conversationRepo.GetByUserID(ctx, userID, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversations: %w", err)
	}

	return conversations, nil
}

func (s *messagingService) UpdateConversationStatus(ctx context.Context, id string, status domain.ConversationStatus, userID string) error {
	conversation, err := s.GetConversation(ctx, id, userID)
	if err != nil {
		return err
	}

	conversation.Status = status
	conversation.UpdatedAt = time.Now()

	if err := s.conversationRepo.Update(ctx, conversation); err != nil {
		return fmt.Errorf("failed to update conversation: %w", err)
	}

	// Invalidate cache
	if s.cacheService != nil {
		_ = s.cacheService.DeleteConversation(ctx, id)
	}

	s.logger.Info("Conversation status updated", map[string]interface{}{
		"conversation_id": id,
		"status":          status,
		"user_id":         userID,
	})

	return nil
}

func (s *messagingService) SendMessage(ctx context.Context, req SendMessageRequest) (*domain.Message, error) {
	// Verify conversation exists and user has access
	_, err := s.GetConversation(ctx, req.ConversationID, req.SenderID)
	if err != nil {
		return nil, err
	}

	message := &domain.Message{
		ID:             uuid.New().String(),
		ConversationID: req.ConversationID,
		SenderType:     req.SenderType,
		SenderID:       req.SenderID,
		Content:        req.Content,
		ContentType:    req.ContentType,
		Metadata:       domain.JSONB(req.Metadata),
		Timestamp:      time.Now(),
	}

	if err := s.messageRepo.Create(ctx, message); err != nil {
		s.logger.Error("Failed to create message", err)
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	// Publish message event
	if s.eventPublisher != nil {
		event := domain.MessageEvent{
			Type:           "message.received",
			ConversationID: message.ConversationID,
			Message:        *message,
			Timestamp:      time.Now(),
		}
		
		if err := s.eventPublisher.PublishMessageEvent(ctx, event); err != nil {
			s.logger.Error("Failed to publish message event", err)
		}
	}

	s.logger.Info("Message sent", map[string]interface{}{
		"message_id":      message.ID,
		"conversation_id": message.ConversationID,
		"sender_id":       message.SenderID,
		"content_type":    message.ContentType,
	})

	return message, nil
}

func (s *messagingService) GetMessages(ctx context.Context, conversationID string, userID string, pagination domain.PaginationParams) ([]domain.Message, error) {
	// Verify conversation access
	_, err := s.GetConversation(ctx, conversationID, userID)
	if err != nil {
		return nil, err
	}

	messages, err := s.messageRepo.GetByConversationID(ctx, conversationID, pagination)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	// Load attachments for each message
	for i := range messages {
		attachments, err := s.attachmentRepo.GetByMessageID(ctx, messages[i].ID)
		if err != nil {
			s.logger.Error("Failed to load attachments for message", err)
			continue
		}
		messages[i].Attachments = attachments
	}

	return messages, nil
}

func (s *messagingService) GetMessage(ctx context.Context, messageID string, userID string) (*domain.Message, error) {
	message, err := s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	// Verify user has access to the conversation
	_, err = s.GetConversation(ctx, message.ConversationID, userID)
	if err != nil {
		return nil, err
	}

	// Load attachments
	attachments, err := s.attachmentRepo.GetByMessageID(ctx, messageID)
	if err != nil {
		s.logger.Error("Failed to load attachments for message", err)
	} else {
		message.Attachments = attachments
	}

	return message, nil
}

func (s *messagingService) CreateAttachment(ctx context.Context, messageID string, req CreateAttachmentRequest) (*domain.Attachment, error) {
	attachment := &domain.Attachment{
		ID:        uuid.New().String(),
		MessageID: messageID,
		URL:       req.URL,
		Type:      req.Type,
		Size:      req.Size,
		Filename:  req.Filename,
		CreatedAt: time.Now(),
	}

	if err := s.attachmentRepo.Create(ctx, attachment); err != nil {
		s.logger.Error("Failed to create attachment", err)
		return nil, fmt.Errorf("failed to create attachment: %w", err)
	}

	s.logger.Info("Attachment created", map[string]interface{}{
		"attachment_id": attachment.ID,
		"message_id":    messageID,
		"type":          attachment.Type,
		"size":          attachment.Size,
	})

	return attachment, nil
}

func (s *messagingService) GetAttachment(ctx context.Context, attachmentID string, userID string) (*domain.Attachment, error) {
	attachment, err := s.attachmentRepo.GetByID(ctx, attachmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get attachment: %w", err)
	}

	// Verify user has access to the message/conversation
	_, err = s.GetMessage(ctx, attachment.MessageID, userID)
	if err != nil {
		return nil, err
	}

	return attachment, nil
}