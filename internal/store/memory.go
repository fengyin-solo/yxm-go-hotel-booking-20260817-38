package store

import (
	"sync"

	"hotelbooking/internal/model"
)

// MemoryStore 基于内存 map 的 Store 实现。
type MemoryStore struct {
	mu         sync.RWMutex
	hotels     map[string]*model.Hotel
	roomTypes  map[string]*model.RoomType
	bookings   map[string]*model.Booking
	checkIns   map[string]*model.CheckIn
	reviews    map[string]*model.Review
	invoices   map[string]*model.Invoice
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		hotels:    make(map[string]*model.Hotel),
		roomTypes: make(map[string]*model.RoomType),
		bookings:  make(map[string]*model.Booking),
		checkIns:  make(map[string]*model.CheckIn),
		reviews:   make(map[string]*model.Review),
		invoices:  make(map[string]*model.Invoice),
	}
}

var _ Store = (*MemoryStore)(nil)
