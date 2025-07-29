package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/company/microservice-template/internal/domain"
	"github.com/company/microservice-template/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type CacheService interface {
	GetConversation(ctx context.Context, id string) (*domain.Conversation, error)
	SetConversation(ctx context.Context, conversation *domain.Conversation) error
	DeleteConversation(ctx context.Context, id string) error
	GetMessages(ctx context.Context, conversationID string) ([]domain.Message, error)
	SetMessages(ctx context.Context, conversationID string, messages []domain.Message) error
	DeleteMessages(ctx context.Context, conversationID string) error
}

type redisCacheService struct {
	client     *redis.Client
	logger     logger.Logger
	expiration time.Duration
}

func NewRedisCacheService(client *redis.Client, logger logger.Logger) CacheService {
	return &redisCacheService{
		client:     client,
		logger:     logger,
		expiration: 30 * time.Minute, // Default cache expiration
	}
}

func (c *redisCacheService) GetConversation(ctx context.Context, id string) (*domain.Conversation, error) {
	key := fmt.Sprintf("conversation:%s", id)
	
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("conversation not found in cache")
		}
		return nil, err
	}

	var conversation domain.Conversation
	if err := json.Unmarshal([]byte(data), &conversation); err != nil {
		c.logger.Error("Failed to unmarshal cached conversation", err)
		return nil, err
	}

	return &conversation, nil
}

func (c *redisCacheService) SetConversation(ctx context.Context, conversation *domain.Conversation) error {
	key := fmt.Sprintf("conversation:%s", conversation.ID)
	
	data, err := json.Marshal(conversation)
	if err != nil {
		c.logger.Error("Failed to marshal conversation for cache", err)
		return err
	}

	if err := c.client.Set(ctx, key, data, c.expiration).Err(); err != nil {
		c.logger.Error("Failed to set conversation in cache", err)
		return err
	}

	return nil
}

func (c *redisCacheService) DeleteConversation(ctx context.Context, id string) error {
	key := fmt.Sprintf("conversation:%s", id)
	
	if err := c.client.Del(ctx, key).Err(); err != nil {
		c.logger.Error("Failed to delete conversation from cache", err)
		return err
	}

	return nil
}

func (c *redisCacheService) GetMessages(ctx context.Context, conversationID string) ([]domain.Message, error) {
	key := fmt.Sprintf("messages:%s", conversationID)
	
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("messages not found in cache")
		}
		return nil, err
	}

	var messages []domain.Message
	if err := json.Unmarshal([]byte(data), &messages); err != nil {
		c.logger.Error("Failed to unmarshal cached messages", err)
		return nil, err
	}

	return messages, nil
}

func (c *redisCacheService) SetMessages(ctx context.Context, conversationID string, messages []domain.Message) error {
	key := fmt.Sprintf("messages:%s", conversationID)
	
	data, err := json.Marshal(messages)
	if err != nil {
		c.logger.Error("Failed to marshal messages for cache", err)
		return err
	}

	// Cache messages for shorter time
	expiration := 10 * time.Minute
	if err := c.client.Set(ctx, key, data, expiration).Err(); err != nil {
		c.logger.Error("Failed to set messages in cache", err)
		return err
	}

	return nil
}

func (c *redisCacheService) DeleteMessages(ctx context.Context, conversationID string) error {
	key := fmt.Sprintf("messages:%s", conversationID)
	
	if err := c.client.Del(ctx, key).Err(); err != nil {
		c.logger.Error("Failed to delete messages from cache", err)
		return err
	}

	return nil
}

// NoOpCacheService for when caching is disabled
type noOpCacheService struct{}

func NewNoOpCacheService() CacheService {
	return &noOpCacheService{}
}

func (c *noOpCacheService) GetConversation(ctx context.Context, id string) (*domain.Conversation, error) {
	return nil, fmt.Errorf("cache disabled")
}

func (c *noOpCacheService) SetConversation(ctx context.Context, conversation *domain.Conversation) error {
	return nil
}

func (c *noOpCacheService) DeleteConversation(ctx context.Context, id string) error {
	return nil
}

func (c *noOpCacheService) GetMessages(ctx context.Context, conversationID string) ([]domain.Message, error) {
	return nil, fmt.Errorf("cache disabled")
}

func (c *noOpCacheService) SetMessages(ctx context.Context, conversationID string, messages []domain.Message) error {
	return nil
}

func (c *noOpCacheService) DeleteMessages(ctx context.Context, conversationID string) error {
	return nil
}