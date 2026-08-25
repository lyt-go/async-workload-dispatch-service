package service_test

import (
	"testing"

	"taskqueue/internal/service"
)

func TestRecoveredTemplateBuildDoesNotExposePartialDefinition(t *testing.T) {
	svc := service.NewTemplateService()
	if err := svc.Create("nightly", []string{"prepare", "panic", "publish"}); err == nil {
		t.Fatalf("recovered template build must return its construction error")
	}
	if count, err := svc.StepCount("nightly"); err == nil {
		t.Fatalf("failed template must not remain readable as a partial definition, count=%d", count)
	}
	if err := svc.Create("hourly", []string{"prepare", "publish"}); err != nil {
		t.Fatalf("later valid template must build after a recovered failure, got %v", err)
	}
	if count, err := svc.StepCount("hourly"); err != nil || count != 2 {
		t.Fatalf("later valid template must remain isolated and ready, count=%d err=%v", count, err)
	}
}
