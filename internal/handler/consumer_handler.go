package handler

import (
	"net/http"

	"taskqueue/internal/model"
	"taskqueue/pkg/httpx"
)

func (s *Server) registerConsumerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/consumers", s.createConsumer)
	mux.HandleFunc("GET /api/consumers", s.listConsumers)
	mux.HandleFunc("GET /api/consumers/{id}", s.getConsumer)
	mux.HandleFunc("PUT /api/consumers/{id}", s.updateConsumer)
	mux.HandleFunc("DELETE /api/consumers/{id}", s.deleteConsumer)
	mux.HandleFunc("POST /api/consumers/{id}/heartbeat", s.heartbeat)
}

type createConsumerRequest struct {
	Name        string `json:"name"`
	Concurrency int    `json:"concurrency"`
	Status      string `json:"status"`
}

func (s *Server) createConsumer(w http.ResponseWriter, r *http.Request) {
	var req createConsumerRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.CreateConsumer(model.Consumer{Name: req.Name, Concurrency: req.Concurrency, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, c)
}

func (s *Server) listConsumers(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ConsumerFilter{
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListConsumers(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getConsumer(w http.ResponseWriter, r *http.Request) {
	c, err := s.svc.GetConsumer(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

func (s *Server) updateConsumer(w http.ResponseWriter, r *http.Request) {
	var req createConsumerRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.UpdateConsumer(r.PathValue("id"), model.Consumer{
		Name: req.Name, Concurrency: req.Concurrency, Status: req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

func (s *Server) deleteConsumer(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteConsumer(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) heartbeat(w http.ResponseWriter, r *http.Request) {
	c, err := s.svc.Heartbeat(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}
