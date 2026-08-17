package handler

import (
	"net/http"
	"strconv"

	"hotelbooking/internal/model"
	"hotelbooking/pkg/httpx"
)

func (s *Server) registerReviewRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/reviews", s.createReview)
	mux.HandleFunc("GET /api/reviews", s.listReviews)
	mux.HandleFunc("GET /api/reviews/{id}", s.getReview)
	mux.HandleFunc("PUT /api/reviews/{id}", s.updateReview)
	mux.HandleFunc("DELETE /api/reviews/{id}", s.deleteReview)
}

type reviewRequest struct {
	BookingID string `json:"booking_id"`
	Rating    int    `json:"rating"`
	Content   string `json:"content"`
}

func (s *Server) createReview(w http.ResponseWriter, r *http.Request) {
	var req reviewRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	review, err := s.svc.CreateReview(model.Review{
		BookingID: req.BookingID,
		Rating:    req.Rating,
		Content:   req.Content,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, review)
}

func (s *Server) listReviews(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	minRating, _ := strconv.Atoi(r.URL.Query().Get("min_rating"))
	filter := model.ReviewFilter{
		BookingID: r.URL.Query().Get("booking_id"),
		MinRating: minRating,
	}
	items, total, err := s.svc.ListReviews(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getReview(w http.ResponseWriter, r *http.Request) {
	review, err := s.svc.GetReview(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, review)
}

func (s *Server) updateReview(w http.ResponseWriter, r *http.Request) {
	var req reviewRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	review, err := s.svc.UpdateReview(r.PathValue("id"), model.Review{
		Rating:  req.Rating,
		Content: req.Content,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, review)
}

func (s *Server) deleteReview(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteReview(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
