package base

import "GoPattern/internal/redisdb"

type BaseHandler struct {
	RedisClient *redisdb.RedisClient
}

func NewBaseHandler(redisClient *redisdb.RedisClient) *BaseHandler {
	return &BaseHandler{RedisClient: redisClient}
}
