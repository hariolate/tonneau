package service

import (
	"context"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"gtihub.com/hariolate/tonneau/config"
)

type Service struct {
	db *gorm.DB
	r  *redis.Client
	c  context.Context
}

func NewService(db *gorm.DB, r *redis.Client, c context.Context) *Service {
	return &Service{db, r, c}
}

func FromConfig(c *config.Parsed) *Service {
	return &Service{
		c.DBConn, c.RedisClient, c.Context,
	}
}
