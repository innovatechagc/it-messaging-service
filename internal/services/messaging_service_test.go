package services

import (
	"context"
	"testing"
	"time"

	"github.com/company/microservice-template/internal/domain"
	"github.com/company/microservice-template/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repositories
type MockConversationRepository struct {
	mock.Mock
}

func (m *MockConversationRepository) Create(ctx context.Context, conversation *domain.Conversation) error {
	args := m.Called(ctx, conversation)
	return args.Error(0)
}

func (m *MockConversationRepository) GetByID(ctx context.Context, id string) (*domain.Conversation, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Conversation), args.Error(1)
}

func (m *MockConversationRepository) GetByUserID(ctx context.Context, userID string, filters domain.ConversationFilters) ([]domain.Conversation, error) {
	args := m.Called(ctx, userID, filters)
	return args.Get(0).([]domain.Conversation), args.Error(1)
}

func (m *MockConversationRepository) Update(ctx context.Context, conversation *domain.Conversation) error {
	args := m.Called(ctx, conversation)
	return args.Error(0)
}

func (m *MockConversationRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockMessageRepository struct {
	mock.Mock
}

func (m *MockMessageRepository) Create(ctx context.Context, message *domain.Message) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

func (m *MockMessageRepository) GetByID(ctx context.Context, id string) (*domain.Message, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Message), args.Error(1)
}

func (m *MockMessageRepository) GetByConversationID(ctx context.Context, conversationID string, pagination domain.PaginationParams) ([]domain.Message, error) {
	args := m.Called(ctx, conversationID, pagination)
	return args.Get(0).([]domain.Message), args.Error(1)
}

func (m *MockMessageRepository) Update(ctx context.Context, message *domain.Message) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

func (m *MockMessageRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockAttachmentRepository struct {
	mock.Mock
}

func (m *MockAttachmentRepository) Create(ctx context.Context, attachment *domain.Attachment) error {
	args := m.Called(ctx, attachment)
	return args.Error(0)
}

func (m *MockAttachmentRepository) GetByID(ctx context.Context, id string) (*domain.Attachment, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Attachment), args.Error(1)
}

func (m *MockAttachmentRepository) GetByMessageID(ctx context.Context, messageID string) ([]domain.Attachment, error) {
	args := m.Called(ctx, messageID)
	return args.Get(0).([]domain.Attachment), args.Error(1)
}

func (m *MockAttachmentRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestMessagingService_CreateConversation(t *testing.T) {
	// Setup
	mockConversationRepo := new(MockConversationRepository)
	mockMessageRepo := new(MockMessageRepository)
	mockAttachmentRepo := new(MockAttachmentRepository)
	mockEventPublisher := NewNoOpEventPublisher()
	mockCacheService := NewNoOpCacheService()
	logger := logger.NewLogger("debug")

	service := NewMessagingService(
		mockConversationRepo,
		mockMessageRepo,
		mockAttachmentRepo,
		mockEventPublisher,
		mockCacheService,
		logger,
	)

	// Test data
	userID := "user123"
	channel := domain.ChannelWeb

	// Mock expectations
	mockConversationRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Conversation")).Return(nil)

	// Execute
	conversation, err := service.CreateConversation(context.Background(), userID, channel)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, conversation)
	assert.Equal(t, userID, conversation.UserID)
	assert.Equal(t, channel, conversation.Channel)
	assert.Equal(t, domain.ConversationStatusActive, conversation.Status)
	assert.NotEmpty(t, conversation.ID)

	mockConversationRepo.AssertExpectations(t)
}

func TestMessagingService_SendMessage(t *testing.T) {
	// Setup
	mockConversationRepo := new(MockConversationRepository)
	mockMessageRepo := new(MockMessageRepository)
	mockAttachmentRepo := new(MockAttachmentRepository)
	mockEventPublisher := NewNoOpEventPublisher()
	mockCacheService := NewNoOpCacheService()
	logger := logger.NewLogger("debug")

	service := NewMessagingService(
		mockConversationRepo,
		mockMessageRepo,
		mockAttachmentRepo,
		mockEventPublisher,
		mockCacheService,
		logger,
	)

	// Test data
	conversationID := "conv123"
	userID := "user123"
	
	existingConversation := &domain.Conversation{
		ID:        conversationID,
		UserID:    userID,
		Channel:   domain.ChannelWeb,
		Status:    domain.ConversationStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	req := SendMessageRequest{
		ConversationID: conversationID,
		SenderType:     domain.SenderTypeUser,
		SenderID:       userID,
		Content:        "Hello, world!",
		ContentType:    domain.ContentTypeText,
		Metadata:       map[string]interface{}{"priority": "normal"},
	}

	// Mock expectations
	mockConversationRepo.On("GetByID", mock.Anything, conversationID).Return(existingConversation, nil)
	mockMessageRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Message")).Return(nil)

	// Execute
	message, err := service.SendMessage(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, message)
	assert.Equal(t, conversationID, message.ConversationID)
	assert.Equal(t, userID, message.SenderID)
	assert.Equal(t, req.Content, message.Content)
	assert.Equal(t, req.ContentType, message.ContentType)
	assert.NotEmpty(t, message.ID)

	mockConversationRepo.AssertExpectations(t)
	mockMessageRepo.AssertExpectations(t)
}

func TestMessagingService_GetConversation(t *testing.T) {
	// Setup
	mockConversationRepo := new(MockConversationRepository)
	mockMessageRepo := new(MockMessageRepository)
	mockAttachmentRepo := new(MockAttachmentRepository)
	mockEventPublisher := NewNoOpEventPublisher()
	mockCacheService := NewNoOpCacheService()
	logger := logger.NewLogger("debug")

	service := NewMessagingService(
		mockConversationRepo,
		mockMessageRepo,
		mockAttachmentRepo,
		mockEventPublisher,
		mockCacheService,
		logger,
	)

	// Test data
	conversationID := "conv123"
	userID := "user123"
	
	expectedConversation := &domain.Conversation{
		ID:        conversationID,
		UserID:    userID,
		Channel:   domain.ChannelWeb,
		Status:    domain.ConversationStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Mock expectations
	mockConversationRepo.On("GetByID", mock.Anything, conversationID).Return(expectedConversation, nil)

	// Execute
	conversation, err := service.GetConversation(context.Background(), conversationID, userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, conversation)
	assert.Equal(t, expectedConversation.ID, conversation.ID)
	assert.Equal(t, expectedConversation.UserID, conversation.UserID)

	mockConversationRepo.AssertExpectations(t)
}

func TestMessagingService_GetConversation_AccessDenied(t *testing.T) {
	// Setup
	mockConversationRepo := new(MockConversationRepository)
	mockMessageRepo := new(MockMessageRepository)
	mockAttachmentRepo := new(MockAttachmentRepository)
	mockEventPublisher := NewNoOpEventPublisher()
	mockCacheService := NewNoOpCacheService()
	logger := logger.NewLogger("debug")

	service := NewMessagingService(
		mockConversationRepo,
		mockMessageRepo,
		mockAttachmentRepo,
		mockEventPublisher,
		mockCacheService,
		logger,
	)

	// Test data
	conversationID := "conv123"
	userID := "user123"
	otherUserID := "user456"
	
	existingConversation := &domain.Conversation{
		ID:        conversationID,
		UserID:    otherUserID, // Different user
		Channel:   domain.ChannelWeb,
		Status:    domain.ConversationStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Mock expectations
	mockConversationRepo.On("GetByID", mock.Anything, conversationID).Return(existingConversation, nil)

	// Execute
	conversation, err := service.GetConversation(context.Background(), conversationID, userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, conversation)
	assert.Contains(t, err.Error(), "not found or access denied")

	mockConversationRepo.AssertExpectations(t)
}