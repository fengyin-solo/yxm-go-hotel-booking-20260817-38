package model

import (
	"strings"
	"time"
)

const (
	RoomTypeActive   = "active"
	RoomTypeInactive = "inactive"
)

// RoomType 表示酒店的房型，价格以「分」为单位存储。
type RoomType struct {
	ID         string    `json:"id"`
	HotelID    string    `json:"hotel_id"`
	Name       string    `json:"name"`
	BedType    string    `json:"bed_type"`
	Price      int64     `json:"price"`
	TotalRooms int       `json:"total_rooms"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (r *RoomType) Validate() error {
	r.Name = strings.TrimSpace(r.Name)
	r.BedType = strings.TrimSpace(r.BedType)
	if r.HotelID == "" {
		return NewValidationError("hotel_id", "必须关联酒店")
	}
	if r.Name == "" {
		return NewValidationError("name", "房型名称不能为空")
	}
	if r.BedType == "" {
		return NewValidationError("bed_type", "床型不能为空")
	}
	if r.Price < 0 {
		return NewValidationError("price", "价格不能为负数")
	}
	if r.TotalRooms <= 0 {
		return NewValidationError("total_rooms", "房间总数必须大于 0")
	}
	if r.Status == "" {
		r.Status = RoomTypeInactive
	}
	if r.Status != RoomTypeActive && r.Status != RoomTypeInactive {
		return NewValidationError("status", "房型状态不合法")
	}
	return nil
}

type RoomTypeFilter struct {
	HotelID string
	Status  string
}

func (f RoomTypeFilter) Match(r *RoomType) bool {
	if f.HotelID != "" && r.HotelID != f.HotelID {
		return false
	}
	if f.Status != "" && r.Status != f.Status {
		return false
	}
	return true
}
