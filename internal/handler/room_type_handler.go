package handler

import (
	"net/http"
	"strconv"

	"hotelbooking/internal/model"
	"hotelbooking/pkg/httpx"
)

func (s *Server) registerRoomTypeRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/room-types", s.createRoomType)
	mux.HandleFunc("GET /api/room-types", s.listRoomTypes)
	mux.HandleFunc("GET /api/room-types/{id}", s.getRoomType)
	mux.HandleFunc("PUT /api/room-types/{id}", s.updateRoomType)
	mux.HandleFunc("DELETE /api/room-types/{id}", s.deleteRoomType)
	mux.HandleFunc("GET /api/room-types/{id}/availability", s.roomTypeAvailability)
	mux.HandleFunc("GET /api/room-types/search", s.searchRoomTypes)
}

type roomTypeRequest struct {
	HotelID    string `json:"hotel_id"`
	Name       string `json:"name"`
	BedType    string `json:"bed_type"`
	Price      int64  `json:"price"`
	TotalRooms int    `json:"total_rooms"`
	Status     string `json:"status"`
}

func (s *Server) createRoomType(w http.ResponseWriter, r *http.Request) {
	var req roomTypeRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	roomType, err := s.svc.CreateRoomType(model.RoomType{
		HotelID:    req.HotelID,
		Name:       req.Name,
		BedType:    req.BedType,
		Price:      req.Price,
		TotalRooms: req.TotalRooms,
		Status:     req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, roomType)
}

func (s *Server) listRoomTypes(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.RoomTypeFilter{
		HotelID: r.URL.Query().Get("hotel_id"),
		Status:  r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListRoomTypes(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getRoomType(w http.ResponseWriter, r *http.Request) {
	roomType, err := s.svc.GetRoomType(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, roomType)
}

func (s *Server) updateRoomType(w http.ResponseWriter, r *http.Request) {
	var req roomTypeRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	roomType, err := s.svc.UpdateRoomType(r.PathValue("id"), model.RoomType{
		Name:       req.Name,
		BedType:    req.BedType,
		Price:      req.Price,
		TotalRooms: req.TotalRooms,
		Status:     req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, roomType)
}

func (s *Server) deleteRoomType(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteRoomType(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) roomTypeAvailability(w http.ResponseWriter, r *http.Request) {
	available, err := s.svc.AvailableRooms(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]int{"available": available})
}

func (s *Server) searchRoomTypes(w http.ResponseWriter, r *http.Request) {
	roomCount, _ := strconv.Atoi(r.URL.Query().Get("rooms"))
	result, err := s.svc.SearchAvailableRoomTypes(r.URL.Query().Get("city"), roomCount)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}
