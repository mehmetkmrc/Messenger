package repository

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
	"github.com/mehmetkmrc/Messenger/internal/core/domain"
)

type MessengerRedisRepository struct {
	client *redis.Client
}

func NewMessengerRedisRepository(host string) *MessengerRedisRepository {
	client := redis.NewClient(&redis.Options{
		Addr: host, 
		Password: "",
		DB: 0,
	})
	return &MessengerRedisRepository{
		client: client,
	}
}

func (r *MessengerRedisRepository) SaveMessage(message domain.Message) error{
	ctx := context.Background()
	json, err := json.Marshal(message)
	if err != nil {
		return err
	}
	r.client.HSet(ctx, "messages", message.ID, json)
	return nil
}

func (r *MessengerRedisRepository) ReadMessage(id string) (*domain.Message, error) {
	ctx := context.Background()
	value, err := r.client.HGet(ctx, "messages", id).Result()
	if err != nil {
		return nil, err
	}

	message := &domain.Message{}
	err = json.Unmarshal([]byte(value), message)
	if err != nil{
		return nil, err
	}
	return message, nil
}

func (r *MessengerRedisRepository) ReadMessages() ([]*domain.Message, error){
	ctx := context.Background()
	messages := []*domain.Message{}
	value, err := r.client.HGetAll(ctx, "messages").Result()
	if err != nil {
		return nil, err
	}


	for _, val := range value {
		message := &domain.Message{}
		err = json.Unmarshal([]byte(val), message)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}