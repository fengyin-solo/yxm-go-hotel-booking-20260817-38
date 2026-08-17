package model

import (
	"strings"
	"time"
)

// Review 表示对已完成订单的评价。
type Review struct {
	ID        string    `json:"id"`
	BookingID string    `json:"booking_id"`
	Rating    int       `json:"rating"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (r *Review) Validate() error {
	r.Content = strings.TrimSpace(r.Content)
	if r.BookingID == "" {
		return NewValidationError("booking_id", "必须关联订单")
	}
	if r.Rating < 1 || r.Rating > 5 {
		return NewValidationError("rating", "评分必须在 1-5 之间")
	}
	if r.Content == "" {
		return NewValidationError("content", "评价内容不能为空")
	}
	return nil
}

type ReviewFilter struct {
	BookingID string
	MinRating int
}

func (f ReviewFilter) Match(r *Review) bool {
	if f.BookingID != "" && r.BookingID != f.BookingID {
		return false
	}
	if f.MinRating > 0 && r.Rating <= f.MinRating {
		return false
	}
	return true
}
