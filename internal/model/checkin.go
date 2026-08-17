package model

import (
	"strings"
	"time"
)

const (
	CheckInActive   = "active"
	CheckInFinished = "finished"
)

// CheckIn 表示一次入住记录，由订单确认后办理入住产生。
type CheckIn struct {
	ID           string     `json:"id"`
	BookingID    string     `json:"booking_id"`
	RoomNumber   string     `json:"room_number"`
	CheckedInAt  time.Time  `json:"checked_in_at"`
	CheckedOutAt *time.Time `json:"checked_out_at"`
	Status       string     `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (c *CheckIn) Validate() error {
	c.RoomNumber = strings.TrimSpace(c.RoomNumber)
	if c.BookingID == "" {
		return NewValidationError("booking_id", "必须关联订单")
	}
	if c.RoomNumber == "" {
		return NewValidationError("room_number", "房间号不能为空")
	}
	if c.Status == "" {
		c.Status = CheckInActive
	}
	if c.Status != CheckInActive && c.Status != CheckInFinished {
		return NewValidationError("status", "入住状态不合法")
	}
	return nil
}
