package service

import (
	"sort"

	"taskqueue/internal/model"
)

// GetDeadLetter 按 ID 获取死信。
func (s *Service) GetDeadLetter(id string) (*model.DeadLetter, error) {
	return s.store.GetDeadLetter(id)
}

// ListDeadLetters 分页列出死信。
func (s *Service) ListDeadLetters(filter model.DeadLetterFilter, page, size int) ([]*model.DeadLetter, int, error) {
	all := s.store.ListDeadLetters()
	matched := make([]*model.DeadLetter, 0, len(all))
	for _, d := range all {
		if filter.Match(d) {
			matched = append(matched, d)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].FailedAt.After(matched[j].FailedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.DeadLetter{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}
