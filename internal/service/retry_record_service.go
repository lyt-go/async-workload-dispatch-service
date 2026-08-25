package service

import (
	"sort"

	"taskqueue/internal/model"
)

// GetRetryRecord 按 ID 获取重试记录。
func (s *Service) GetRetryRecord(id string) (*model.RetryRecord, error) {
	return s.store.GetRetryRecord(id)
}

// ListRetryRecords 分页列出重试记录。
func (s *Service) ListRetryRecords(filter model.RetryRecordFilter, page, size int) ([]*model.RetryRecord, int, error) {
	all := s.store.ListRetryRecords()
	matched := make([]*model.RetryRecord, 0, len(all))
	for _, r := range all {
		if filter.Match(r) {
			matched = append(matched, r)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.RetryRecord{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}
