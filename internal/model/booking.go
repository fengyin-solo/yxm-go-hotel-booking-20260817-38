package model

import (
	"strings"
	"time"
)

const (
	BookingPending   = "pending"
	BookingConfirmed = "confirmed"
	BookingCheckedIn = "checked_in"
	BookingCompleted = "completed"
	BookingCancelled = "cancelled"
)

// bookingTransitions 定义订单状态机的合法流转。
var bookingTransitions = map[string]map[string]bool{
	BookingPending:   {BookingConfirmed: true, BookingCancelled: true},
	BookingConfirmed: {BookingCheckedIn: true, BookingCancelled: true},
	BookingCheckedIn: {},
	BookingCompleted: {},
	BookingCancelled: {},
}

// CanTransitionBooking 判断订单能否从 from 流转到 to。
func CanTransitionBooking(from, to string) bool {
	if m, ok := bookingTransitions[from]; ok {
		return m[to]
	}
	return false
}

// ValidBookingStatus 判断订单状态是否合法。
func ValidBookingStatus(s string) bool {
	_, ok := bookingTransitions[s]
	return ok
}

// Booking 表示一个酒店预订订单，金额以「分」为单位。
type Booking struct {
	ID          string    `json:"id"`
	RoomTypeID  string    `json:"room_type_id"`
	GuestName   string    `json:"guest_name"`
	GuestPhone  string    `json:"guest_phone"`
	CheckIn     time.Time `json:"check_in"`
	CheckOut    time.Time `json:"check_out"`
	RoomCount   int       `json:"room_count"`
	TotalAmount int64     `json:"total_amount"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Nights 返回入住晚数。
func (b *Booking) Nights() int {
	hours := b.CheckOut.Sub(b.CheckIn).Hours()
	nights := int(hours / 24)
	if nights < 1 {
		nights = 1
	}
	return nights
}

func (b *Booking) Validate() error {
	b.GuestName = strings.TrimSpace(b.GuestName)
	b.GuestPhone = strings.TrimSpace(b.GuestPhone)
	if b.RoomTypeID == "" {
		return NewValidationError("room_type_id", "必须关联房型")
	}
	if b.GuestName == "" {
		return NewValidationError("guest_name", "客人姓名不能为空")
	}
	if b.GuestPhone == "" {
		return NewValidationError("guest_phone", "客人电话不能为空")
	}
	if b.RoomCount <= 0 {
		return NewValidationError("room_count", "房间数量必须大于 0")
	}
	if !b.CheckOut.After(b.CheckIn) {
		return NewValidationError("check_out", "离店日期必须晚于入住日期")
	}
	if b.Status == "" {
		b.Status = BookingPending
	}
	if !ValidBookingStatus(b.Status) {
		return NewValidationError("status", "订单状态不合法")
	}
	return nil
}

type BookingFilter struct {
	RoomTypeID string
	Status     string
	GuestName  string
	From       time.Time
	To         time.Time
}

func (f BookingFilter) Match(b *Booking) bool {
	if f.RoomTypeID != "" && b.RoomTypeID != f.RoomTypeID {
		return false
	}
	if f.Status != "" && b.Status != f.Status {
		return false
	}
	if f.GuestName != "" {
		k := strings.ToLower(strings.TrimSpace(f.GuestName))
		if k != "" && !strings.Contains(strings.ToLower(b.GuestName), k) {
			return false
		}
	}
	if !f.From.IsZero() && b.CheckIn.Before(f.From) {
		return false
	}
	if !f.To.IsZero() && b.CheckIn.After(f.To) {
		return false
	}
	return true
}
