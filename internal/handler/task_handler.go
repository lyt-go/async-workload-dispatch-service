package handler

import (
	"net/http"

	"taskqueue/internal/model"
	"taskqueue/pkg/httpx"
)

func (s *Server) registerTaskRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/queues/{id}/tasks", s.enqueue)
	mux.HandleFunc("POST /api/queues/{id}/dequeue", s.dequeue)
	mux.HandleFunc("GET /api/tasks", s.listTasks)
	mux.HandleFunc("GET /api/tasks/{id}", s.getTask)
	mux.HandleFunc("DELETE /api/tasks/{id}", s.deleteTask)
	mux.HandleFunc("POST /api/tasks/{id}/ack", s.ackTask)
	mux.HandleFunc("POST /api/tasks/{id}/fail", s.failTask)
	mux.HandleFunc("POST /api/tasks/{id}/requeue", s.requeueDeadLetter)
}

type enqueueRequest struct {
	Payload  string `json:"payload"`
	Priority int    `json:"priority"`
}

func (s *Server) enqueue(w http.ResponseWriter, r *http.Request) {
	var req enqueueRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.Enqueue(r.PathValue("id"), req.Payload, req.Priority)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, t)
}

func (s *Server) dequeue(w http.ResponseWriter, r *http.Request) {
	t, err := s.svc.Dequeue(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

func (s *Server) listTasks(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.TaskFilter{
		QueueID: r.URL.Query().Get("queue_id"),
		Status:  r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListTasks(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getTask(w http.ResponseWriter, r *http.Request) {
	t, err := s.svc.GetTask(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

func (s *Server) deleteTask(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteTask(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) ackTask(w http.ResponseWriter, r *http.Request) {
	t, err := s.svc.AckTask(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

type failTaskRequest struct {
	Reason string `json:"reason"`
}

func (s *Server) failTask(w http.ResponseWriter, r *http.Request) {
	var req failTaskRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.FailTask(r.PathValue("id"), req.Reason)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

func (s *Server) requeueDeadLetter(w http.ResponseWriter, r *http.Request) {
	t, err := s.svc.RequeueDeadLetter(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}
