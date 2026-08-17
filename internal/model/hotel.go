package model

import (
	"strings"
	"time"
)

const (
	HotelActive   = "active"
	HotelInactive = "inactive"
)

// Hotel 表示一家酒店。
type Hotel struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	City      string    `json:"city"`
	Address   string    `json:"address"`
	Star      int       `json:"star"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (h *Hotel) Validate() error {
	h.Name = strings.TrimSpace(h.Name)
	h.City = strings.TrimSpace(h.City)
	h.Address = strings.TrimSpace(h.Address)
	if h.Name == "" {
		return NewValidationError("name", "酒店名称不能为空")
	}
	if h.City == "" {
		return NewValidationError("city", "所在城市不能为空")
	}
	if h.Address == "" {
		return NewValidationError("address", "酒店地址不能为空")
	}
	if h.Star < 0 || h.Star > 5 {
		return NewValidationError("star", "星级必须在 0-5 之间")
	}
	if h.Status == "" {
		h.Status = HotelActive
	}
	if h.Status != HotelActive && h.Status != HotelInactive {
		return NewValidationError("status", "酒店状态不合法")
	}
	return nil
}

type HotelFilter struct {
	City    string
	Star    int
	Status  string
	Keyword string
}

func (f HotelFilter) Match(h *Hotel) bool {
	if f.City != "" && h.City != f.City {
		return false
	}
	if f.Star > 0 && h.Star != f.Star {
		return false
	}
	if f.Status != "" && h.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(h.Name), k) &&
			!strings.Contains(strings.ToLower(h.Address), k) {
			return false
		}
	}
	return true
}
