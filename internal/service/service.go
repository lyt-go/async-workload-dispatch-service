// Package service 实现业务逻辑层。
package service

import (
	"taskqueue/internal/config"
	"taskqueue/internal/store"
	"taskqueue/pkg/logger"
)

// Service 聚合存储与日志，承载全部业务逻辑。
type Service struct {
	store store.Store
	log   *logger.Logger
	cfg   *config.Config
}

// New 创建业务服务。
func New(st store.Store, log *logger.Logger, cfg *config.Config) *Service {
	return &Service{store: st, log: log, cfg: cfg}
}
