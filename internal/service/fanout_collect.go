package service

import (
	"context"

	"taskqueue/internal/fanout"
)

type FanoutService struct{}

func (FanoutService) Collect(ctx context.Context, items []fanout.Item) ([]string, error) {
	output, _ := fanout.Produce(items)
	ids := make([]string, 0, len(items))
	for item := range output {
		ids = append(ids, item.ID)
	}
	return ids, nil
}
