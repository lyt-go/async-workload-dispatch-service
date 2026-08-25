package handler

import (
	"net/http"

	"taskqueue/internal/model"
	"taskqueue/pkg/httpx"
)

func (s *Server) registerQueueRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/queues", s.createQueue)
	mux.HandleFunc("GET /api/queues", s.listQueues)
	mux.HandleFunc("GET /api/queues/{id}", s.getQueue)
	mux.HandleFunc("PUT /api/queues/{id}", s.updateQueue)
	mux.HandleFunc("DELETE /api/queues/{id}", s.deleteQueue)
	mux.HandleFunc("POST /api/queues/{id}/status", s.changeQueueStatus)
}

type createQueueRequest struct {
	Name              string `json:"name"`
	Topic             string `json:"topic"`
	MaxRetry          int    `json:"max_retry"`
	VisibilityTimeout int    `json:"visibility_timeout"`
	Status            string `json:"status"`
}

func (s *Server) createQueue(w http.ResponseWriter, r *http.Request) {
	var req createQueueRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	q, err := s.svc.CreateQueue(model.Queue{
		Name: req.Name, Topic: req.Topic, MaxRetry: req.MaxRetry,
		VisibilityTimeout: req.VisibilityTimeout, Status: req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, q)
}

func (s *Server) listQueues(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.QueueFilter{
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListQueues(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getQueue(w http.ResponseWriter, r *http.Request) {
	q, err := s.svc.GetQueue(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, q)
}

func (s *Server) updateQueue(w http.ResponseWriter, r *http.Request) {
	var req createQueueRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	q, err := s.svc.UpdateQueue(r.PathValue("id"), model.Queue{
		Name: req.Name, Topic: req.Topic, MaxRetry: req.MaxRetry,
		VisibilityTimeout: req.VisibilityTimeout, Status: req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, q)
}

func (s *Server) deleteQueue(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteQueue(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type changeQueueStatusRequest struct {
	Status string `json:"status"`
}

func (s *Server) changeQueueStatus(w http.ResponseWriter, r *http.Request) {
	var req changeQueueStatusRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	q, err := s.svc.ChangeQueueStatus(r.PathValue("id"), req.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, q)
}
