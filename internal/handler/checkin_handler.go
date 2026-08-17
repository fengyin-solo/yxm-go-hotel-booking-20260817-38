package handler

import (
	"net/http"

	"hotelbooking/pkg/httpx"
)

func (s *Server) registerCheckInRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/checkins", s.listCheckIns)
	mux.HandleFunc("GET /api/checkins/{id}", s.getCheckIn)
	mux.HandleFunc("GET /api/bookings/{id}/checkin", s.getCheckInByBooking)
}

func (s *Server) listCheckIns(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	items, total, err := s.svc.ListCheckIns(
		r.URL.Query().Get("booking_id"),
		r.URL.Query().Get("status"),
		pp.Page, pp.Size,
	)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getCheckIn(w http.ResponseWriter, r *http.Request) {
	checkIn, err := s.svc.GetCheckIn(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, checkIn)
}

func (s *Server) getCheckInByBooking(w http.ResponseWriter, r *http.Request) {
	checkIn, err := s.svc.GetCheckInByBooking(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, checkIn)
}
