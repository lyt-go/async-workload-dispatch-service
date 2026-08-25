// Package probe 提供探测请求的上下文隔离工具。
//
// 每个探测请求只使用自己的 context，互不串扰：一个请求被取消不会把
// 取消状态泄漏给后续请求，也不会换一个无取消信号的新背景继续调用下游。
package probe

// ContextStore 仅作历史占位保留。
//
// 历史实现会把第一个请求的 context 缓存到共享字段，导致后续所有请求复用
// 同一个（可能已取消的）context，破坏请求隔离。该类型已不再被使用，
// 探测请求应直接使用自身的 context，详见 service.ProbeService.Probe。
type ContextStore struct{}
