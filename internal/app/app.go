// Package app 负责依赖装配。
package app

import (
	"net/http"

	"taskqueue/internal/config"
	"taskqueue/internal/handler"
	"taskqueue/internal/service"
	"taskqueue/internal/store"
	"taskqueue/pkg/logger"
)

// App 聚合全部依赖并暴露 HTTP 路由。
type App struct {
	server *handler.Server
}

// New 装配 store -> service -> handler 依赖链。
func New(cfg *config.Config, log *logger.Logger) (*App, error) {
	st := store.NewMemoryStore()
	svc := service.New(st, log, cfg)
	server := handler.NewServer(svc, log, cfg)
	log.Infof("应用装配完成，配置：%s", cfg.String())
	return &App{server: server}, nil
}

// Routes 返回应用根路由。
func (a *App) Routes() http.Handler { return a.server.Routes() }
