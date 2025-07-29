package services

import (
	"context"
	"encoding/json"

	"github.com/company/microservice-template/internal/domain"
	"github.com/company/microservice-template/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type EventPublisher interface {
	PublishMessageEvent(ctx context.Context, event domain.MessageEvent) error
}

type redisEventPublisher struct {
	client *redis.Client
	topic  string
	logger logger.Logger
}

func NewRedisEventPublisher(client *redis.Client, topic string, logger logger.Logger) EventPublisher {
	return &redisEventPublisher{
		client: client,
		topic:  topic,
		logger: logger,
	}
}

func (p *redisEventPublisher) PublishMessageEvent(ctx context.Context, event domain.MessageEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		p.logger.Error("Failed to marshal event", err)
		return err
	}

	if err := p.client.Publish(ctx, p.topic, data).Err(); err != nil {
		p.logger.Error("Failed to publish event to Redis", err)
		return err
	}

	p.logger.Info("Event published", map[string]interface{}{
		"topic":           p.topic,
		"event_type":      event.Type,
		"conversation_id": event.ConversationID,
	})

	return nil
}

// NoOpEventPublisher for when events are disabled
type noOpEventPublisher struct{}

func NewNoOpEventPublisher() EventPublisher {
	return &noOpEventPublisher{}
}

func (p *noOpEventPublisher) PublishMessageEvent(ctx context.Context, event domain.MessageEvent) error {
	// Do nothing
	return nil
}