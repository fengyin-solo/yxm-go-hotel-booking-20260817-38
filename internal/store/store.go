// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"hotelbooking/internal/model"
)

var (
	ErrNotFound = errors.New("记录不存在")
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法，便于测试时替换实现。
type Store interface {
	CreateHotel(h *model.Hotel) error
	GetHotel(id string) (*model.Hotel, error)
	ListHotels() []*model.Hotel
	UpdateHotel(h *model.Hotel) error
	DeleteHotel(id string) error

	CreateRoomType(r *model.RoomType) error
	GetRoomType(id string) (*model.RoomType, error)
	ListRoomTypes() []*model.RoomType
	UpdateRoomType(r *model.RoomType) error
	DeleteRoomType(id string) error

	CreateBooking(b *model.Booking) error
	GetBooking(id string) (*model.Booking, error)
	ListBookings() []*model.Booking
	UpdateBooking(b *model.Booking) error
	DeleteBooking(id string) error

	CreateCheckIn(c *model.CheckIn) error
	GetCheckIn(id string) (*model.CheckIn, error)
	GetCheckInByBooking(bookingID string) (*model.CheckIn, error)
	ListCheckIns() []*model.CheckIn
	UpdateCheckIn(c *model.CheckIn) error

	CreateReview(r *model.Review) error
	GetReview(id string) (*model.Review, error)
	ListReviews() []*model.Review
	UpdateReview(r *model.Review) error
	DeleteReview(id string) error

	CreateInvoice(i *model.Invoice) error
	GetInvoice(id string) (*model.Invoice, error)
	ListInvoices() []*model.Invoice
	UpdateInvoice(i *model.Invoice) error
	DeleteInvoice(id string) error
}
