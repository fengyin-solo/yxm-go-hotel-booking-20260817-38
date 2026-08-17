package handler

import (
	"net/http"
	"time"

	"hotelbooking/internal/model"
	"hotelbooking/pkg/httpx"
)

func (s *Server) registerBookingRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/bookings", s.createBooking)
	mux.HandleFunc("GET /api/bookings", s.listBookings)
	mux.HandleFunc("GET /api/bookings/{id}", s.getBooking)
	mux.HandleFunc("PUT /api/bookings/{id}", s.updateBooking)
	mux.HandleFunc("POST /api/bookings/{id}/confirm", s.confirmBooking)
	mux.HandleFunc("POST /api/bookings/{id}/cancel", s.cancelBooking)
	mux.HandleFunc("POST /api/bookings/{id}/check-in", s.checkInBooking)
	mux.HandleFunc("POST /api/bookings/{id}/check-out", s.checkOutBooking)
	mux.HandleFunc("POST /api/bookings/batch-confirm", s.batchConfirm)
}

func parseDate(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}

type bookingRequest struct {
	RoomTypeID string `json:"room_type_id"`
	GuestName  string `json:"guest_name"`
	GuestPhone string `json:"guest_phone"`
	CheckIn    string `json:"check_in"`
	CheckOut   string `json:"check_out"`
	RoomCount  int    `json:"room_count"`
}

func (s *Server) createBooking(w http.ResponseWriter, r *http.Request) {
	var req bookingRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	checkIn, err := parseDate(req.CheckIn)
	if err != nil {
		httpx.BadRequest(w, "入住日期格式错误，应为 YYYY-MM-DD")
		return
	}
	checkOut, err := parseDate(req.CheckOut)
	if err != nil {
		httpx.BadRequest(w, "离店日期格式错误，应为 YYYY-MM-DD")
		return
	}
	booking, err := s.svc.CreateBooking(model.Booking{
		RoomTypeID: req.RoomTypeID,
		GuestName:  req.GuestName,
		GuestPhone: req.GuestPhone,
		CheckIn:    checkIn,
		CheckOut:   checkOut,
		RoomCount:  req.RoomCount,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, booking)
}

func (s *Server) listBookings(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.BookingFilter{
		RoomTypeID: r.URL.Query().Get("room_type_id"),
		Status:     r.URL.Query().Get("status"),
		GuestName:  r.URL.Query().Get("guest_name"),
	}
	if v := r.URL.Query().Get("from"); v != "" {
		if t, err := parseDate(v); err == nil {
			filter.From = t
		}
	}
	if v := r.URL.Query().Get("to"); v != "" {
		if t, err := parseDate(v); err == nil {
			filter.To = t
		}
	}
	items, total, err := s.svc.ListBookings(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getBooking(w http.ResponseWriter, r *http.Request) {
	booking, err := s.svc.GetBooking(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, booking)
}

func (s *Server) updateBooking(w http.ResponseWriter, r *http.Request) {
	var req bookingRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	booking, err := s.svc.UpdateBooking(r.PathValue("id"), model.Booking{
		GuestName:  req.GuestName,
		GuestPhone: req.GuestPhone,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, booking)
}

func (s *Server) confirmBooking(w http.ResponseWriter, r *http.Request) {
	booking, err := s.svc.ConfirmBooking(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, booking)
}

func (s *Server) cancelBooking(w http.ResponseWriter, r *http.Request) {
	booking, err := s.svc.CancelBooking(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, booking)
}

type checkInRequest struct {
	RoomNumber string `json:"room_number"`
}

func (s *Server) checkInBooking(w http.ResponseWriter, r *http.Request) {
	var req checkInRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	booking, err := s.svc.CheckInBooking(r.PathValue("id"), req.RoomNumber)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, booking)
}

func (s *Server) checkOutBooking(w http.ResponseWriter, r *http.Request) {
	booking, err := s.svc.CheckOutBooking(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, booking)
}

type batchConfirmRequest struct {
	IDs []string `json:"ids"`
}

func (s *Server) batchConfirm(w http.ResponseWriter, r *http.Request) {
	var req batchConfirmRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	if len(req.IDs) == 0 {
		httpx.BadRequest(w, "ids 不能为空")
		return
	}
	httpx.OK(w, s.svc.BatchConfirm(req.IDs))
}
