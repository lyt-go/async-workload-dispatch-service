package service

import (
	"context"
)

// ProbeService 负责向下游发起探测调用。
//
// 每个探测请求严格隔离：仅使用本次请求自身的 context，
// 既不会把上一个请求（可能已取消）的 context 串给下一个请求，
// 也不会在请求被取消后换一个无取消信号的新背景继续调用下游。
type ProbeService struct{}

// New 创建探测服务。探测服务不持有任何跨请求状态，
// 因此每个请求都从自身的 context 出发。
func NewProbeService() *ProbeService { return &ProbeService{} }

// Probe 使用传入的 ctx 调用下游。
//
// - 请求被取消时立即返回，绝不另起 context.Background() 重试
//   （历史实现正是在失败后换一个新背景继续调用，把被取消的请求"复活"成成功）。
// - 下游返回的任何错误原样上抛，交由调用方处理，不吞错、不无谓重试。
func (s *ProbeService) Probe(ctx context.Context, call func(context.Context) error) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return call(ctx)
}
