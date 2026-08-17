package service

import (
	"hotelbooking/internal/config"
	"hotelbooking/internal/store"
	"hotelbooking/pkg/logger"
)

// Service 聚合业务逻辑，依赖 Store 接口与配置。
type Service struct {
	store store.Store
	log   *logger.Logger
	cfg   *config.Config
}

func New(st store.Store, log *logger.Logger, cfg *config.Config) *Service {
	return &Service{store: st, log: log, cfg: cfg}
}
