package service

import (
	"sort"
	"time"

	"hotelbooking/internal/model"
	"hotelbooking/pkg/idgen"
)

// HotelRating 酒店聚合评分结果。
type HotelRating struct {
	HotelID    string  `json:"hotel_id"`
	ReviewCount int    `json:"review_count"`
	AvgRating  float64 `json:"avg_rating"`
}

func (s *Service) CreateReview(review model.Review) (*model.Review, error) {
	if err := review.Validate(); err != nil {
		return nil, err
	}
	booking, err := s.store.GetBooking(review.BookingID)
	if err != nil {
		return nil, model.NewValidationError("booking_id", "关联的订单不存在")
	}
	if booking.Status != model.BookingCompleted {
		return nil, model.NewValidationError("booking_id", "仅已完成订单可评价")
	}
	for _, r := range s.store.ListReviews() {
		if r.BookingID == review.BookingID {
			return nil, model.NewValidationError("booking_id", "该订单已评价过")
		}
	}
	now := time.Now()
	review.ID = idgen.Hex()
	review.CreatedAt = now
	review.UpdatedAt = now
	if err := s.store.CreateReview(&review); err != nil {
		return nil, err
	}
	return &review, nil
}

func (s *Service) GetReview(id string) (*model.Review, error) {
	return s.store.GetReview(id)
}

func (s *Service) ListReviews(filter model.ReviewFilter, page, size int) ([]*model.Review, int, error) {
	all := s.store.ListReviews()
	matched := make([]*model.Review, 0, len(all))
	for _, r := range all {
		if filter.Match(r) {
			matched = append(matched, r)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Review{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateReview(id string, input model.Review) (*model.Review, error) {
	existing, err := s.store.GetReview(id)
	if err != nil {
		return nil, err
	}
	if input.Rating > 0 {
		existing.Rating = input.Rating
	}
	if input.Content != "" {
		existing.Content = input.Content
	}
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	existing.UpdatedAt = time.Now()
	if err := s.store.UpdateReview(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) DeleteReview(id string) error {
	return s.store.DeleteReview(id)
}

// HotelRating 计算指定酒店的平均评分与评价数量。
func (s *Service) HotelRating(hotelID string) (*HotelRating, error) {
	if _, err := s.store.GetHotel(hotelID); err != nil {
		return nil, err
	}
	roomTypeIDs := make(map[string]bool)
	for _, rt := range s.store.ListRoomTypes() {
		if rt.HotelID == hotelID {
			roomTypeIDs[rt.ID] = true
		}
	}
	bookingIDs := make(map[string]bool)
	for _, b := range s.store.ListBookings() {
		if roomTypeIDs[b.RoomTypeID] {
			bookingIDs[b.ID] = true
		}
	}
	rating := &HotelRating{HotelID: hotelID}
	total := 0
	for _, r := range s.store.ListReviews() {
		if !bookingIDs[r.BookingID] {
			continue
		}
		rating.ReviewCount++
		total += r.Rating
	}
	if rating.ReviewCount > 0 {
		rating.AvgRating = float64(total) / float64(rating.ReviewCount)
	}
	return rating, nil
}
