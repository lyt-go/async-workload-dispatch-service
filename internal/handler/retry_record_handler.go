package handler

import (
	"net/http"

	"taskqueue/internal/model"
	"taskqueue/pkg/httpx"
)

func (s *Server) registerRetryRecordRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/retry-records", s.listRetryRecords)
	mux.HandleFunc("GET /api/retry-records/{id}", s.getRetryRecord)
}

func (s *Server) listRetryRecords(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.RetryRecordFilter{
		TaskID: r.URL.Query().Get("task_id"),
	}
	items, total, err := s.svc.ListRetryRecords(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getRetryRecord(w http.ResponseWriter, r *http.Request) {
	rec, err := s.svc.GetRetryRecord(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rec)
}
