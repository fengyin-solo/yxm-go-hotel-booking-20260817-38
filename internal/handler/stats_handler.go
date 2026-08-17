package handler

import (
	"net/http"

	"hotelbooking/pkg/httpx"
)

func (s *Server) registerStatsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/stats/booking-status", s.bookingStatusStats)
	mux.HandleFunc("GET /api/stats/revenue/city", s.revenueByCity)
	mux.HandleFunc("GET /api/hotels/{id}/revenue", s.revenueByHotel)
	mux.HandleFunc("GET /api/hotels/{id}/rating", s.hotelRating)
	mux.HandleFunc("GET /api/hotels/{id}/report", s.hotelReport)
	mux.HandleFunc("GET /api/room-types/{id}/occupancy", s.roomTypeOccupancy)
}

func (s *Server) bookingStatusStats(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.BookingStatusStats())
}

func (s *Server) revenueByCity(w http.ResponseWriter, r *http.Request) {
	result, err := s.svc.RevenueByCity()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

func (s *Server) revenueByHotel(w http.ResponseWriter, r *http.Request) {
	result, err := s.svc.RevenueByHotel(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

func (s *Server) hotelRating(w http.ResponseWriter, r *http.Request) {
	result, err := s.svc.HotelRating(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

func (s *Server) hotelReport(w http.ResponseWriter, r *http.Request) {
	result, err := s.svc.ExportHotelReport(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

func (s *Server) roomTypeOccupancy(w http.ResponseWriter, r *http.Request) {
	result, err := s.svc.OccupancyByRoomType(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}
