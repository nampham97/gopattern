package redisdb

import (
	"context"
	"encoding/json"
	"time"
)

type Job struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type Queue struct {
	redis *RedisClient
}

func NewQueue(redis *RedisClient) *Queue {
	return &Queue{redis: redis}
}

func (q *Queue) Push(ctx context.Context, queueName string, job Job) error {
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}
	return q.redis.client.LPush(ctx, queueName, data).Err()
}

func (q *Queue) Pop(ctx context.Context, queueName string, timeout time.Duration) (*Job, error) {
	res, err := q.redis.client.BRPop(ctx, timeout, queueName).Result()
	if err != nil || len(res) < 2 {
		return nil, err
	}

	var job Job
	err = json.Unmarshal([]byte(res[1]), &job)
	if err != nil {
		return nil, err
	}
	return &job, nil
}
