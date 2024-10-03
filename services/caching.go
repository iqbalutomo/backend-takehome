package services

import (
	"backend-takehome/models"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type CachingService interface {
	SetPostDetailed(postID uint, data *models.PostDetail) error
	GetPostDetailed(postID uint) (*models.PostDetail, error)
}

type RedisClient struct {
	client *redis.Client
}

func NewCachingService(client *redis.Client) CachingService {
	return &RedisClient{client}
}

func (r *RedisClient) SetPostDetailed(postID uint, data *models.PostDetail) error {
	key := fmt.Sprintf("postdetailed:%v", postID)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dataJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	result := r.client.Set(ctx, key, dataJSON, 2*time.Hour)
	if err := result.Err(); err != nil {
		return err
	}

	return nil
}

func (r *RedisClient) GetPostDetailed(postID uint) (*models.PostDetail, error) {
	key := fmt.Sprintf("postdetailed:%v", postID)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := r.client.Get(ctx, key)
	if result.Err() == redis.Nil {
		return nil, nil
	} else if result.Err() != nil {
		return nil, result.Err()
	}

	var postDetailData models.PostDetail
	if err := json.Unmarshal([]byte(result.Val()), &postDetailData); err != nil {
		return nil, err
	}

	return &postDetailData, nil
}
