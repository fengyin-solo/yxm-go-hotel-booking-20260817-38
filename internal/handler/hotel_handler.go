package handler

import (
	"net/http"
	"strconv"

	"hotelbooking/internal/model"
	"hotelbooking/pkg/httpx"
)

func (s *Server) registerHotelRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/hotels", s.createHotel)
	mux.HandleFunc("GET /api/hotels", s.listHotels)
	mux.HandleFunc("GET /api/hotels/{id}", s.getHotel)
	mux.HandleFunc("PUT /api/hotels/{id}", s.updateHotel)
	mux.HandleFunc("DELETE /api/hotels/{id}", s.deleteHotel)
}

type hotelRequest struct {
	Name    string `json:"name"`
	City    string `json:"city"`
	Address string `json:"address"`
	Star    int    `json:"star"`
	Status  string `json:"status"`
}

func (s *Server) createHotel(w http.ResponseWriter, r *http.Request) {
	var req hotelRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	hotel, err := s.svc.CreateHotel(model.Hotel{
		Name:    req.Name,
		City:    req.City,
		Address: req.Address,
		Star:    req.Star,
		Status:  req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, hotel)
}

func (s *Server) listHotels(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	star, _ := strconv.Atoi(r.URL.Query().Get("star"))
	filter := model.HotelFilter{
		City:    r.URL.Query().Get("city"),
		Star:    star,
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListHotels(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getHotel(w http.ResponseWriter, r *http.Request) {
	hotel, err := s.svc.GetHotel(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, hotel)
}

func (s *Server) updateHotel(w http.ResponseWriter, r *http.Request) {
	var req hotelRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	hotel, err := s.svc.UpdateHotel(r.PathValue("id"), model.Hotel{
		Name:    req.Name,
		City:    req.City,
		Address: req.Address,
		Star:    req.Star,
		Status:  req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, hotel)
}

func (s *Server) deleteHotel(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteHotel(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
