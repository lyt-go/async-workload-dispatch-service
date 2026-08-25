package handler

import (
	"net/http"

	"taskqueue/internal/model"
	"taskqueue/pkg/httpx"
)

func (s *Server) registerSubscriptionRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/subscriptions", s.createSubscription)
	mux.HandleFunc("GET /api/subscriptions", s.listSubscriptions)
	mux.HandleFunc("GET /api/subscriptions/{id}", s.getSubscription)
	mux.HandleFunc("PUT /api/subscriptions/{id}", s.updateSubscription)
	mux.HandleFunc("DELETE /api/subscriptions/{id}", s.deleteSubscription)
}

type createSubscriptionRequest struct {
	ConsumerID string `json:"consumer_id"`
	QueueID    string `json:"queue_id"`
	Status     string `json:"status"`
}

func (s *Server) createSubscription(w http.ResponseWriter, r *http.Request) {
	var req createSubscriptionRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sub, err := s.svc.CreateSubscription(model.Subscription{
		ConsumerID: req.ConsumerID, QueueID: req.QueueID, Status: req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, sub)
}

func (s *Server) listSubscriptions(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.SubscriptionFilter{
		ConsumerID: r.URL.Query().Get("consumer_id"),
		QueueID:    r.URL.Query().Get("queue_id"),
		Status:     r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListSubscriptions(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getSubscription(w http.ResponseWriter, r *http.Request) {
	sub, err := s.svc.GetSubscription(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sub)
}

func (s *Server) updateSubscription(w http.ResponseWriter, r *http.Request) {
	var req createSubscriptionRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sub, err := s.svc.UpdateSubscription(r.PathValue("id"), model.Subscription{Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sub)
}

func (s *Server) deleteSubscription(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteSubscription(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
