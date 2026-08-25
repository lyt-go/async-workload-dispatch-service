package handler

import (
	"net/http"

	"taskqueue/internal/model"
	"taskqueue/pkg/httpx"
)

func (s *Server) registerDeadLetterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/dead-letters", s.listDeadLetters)
	mux.HandleFunc("GET /api/dead-letters/{id}", s.getDeadLetter)
}

func (s *Server) listDeadLetters(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.DeadLetterFilter{
		QueueID: r.URL.Query().Get("queue_id"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListDeadLetters(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getDeadLetter(w http.ResponseWriter, r *http.Request) {
	d, err := s.svc.GetDeadLetter(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, d)
}
