package service

import (
	"testing"
	"time"

	"hotelbooking/internal/config"
	"hotelbooking/internal/model"
	"hotelbooking/internal/store"
	"hotelbooking/pkg/logger"
)

func newTestService() *Service {
	cfg := &config.Config{MaxPageSize: 100, LowStockThreshold: 5}
	log := logger.NewLevel(logger.LevelError)
	return New(store.NewMemoryStore(), log, cfg)
}

func mustCreateHotel(t *testing.T, s *Service, name, city string) *model.Hotel {
	t.Helper()
	h, err := s.CreateHotel(model.Hotel{
		Name:    name,
		City:    city,
		Address: city + "某街 1 号",
		Star:    4,
	})
	if err != nil {
		t.Fatalf("创建酒店失败: %v", err)
	}
	return h
}

func mustCreateRoomType(t *testing.T, s *Service, hotelID, name string, price int64, total int) *model.RoomType {
	t.Helper()
	r, err := s.CreateRoomType(model.RoomType{
		HotelID:    hotelID,
		Name:       name,
		BedType:    "大床",
		Price:      price,
		TotalRooms: total,
	})
	if err != nil {
		t.Fatalf("创建房型失败: %v", err)
	}
	return r
}

func mustCreateBooking(t *testing.T, s *Service, roomTypeID, guest string, roomCount int) *model.Booking {
	t.Helper()
	now := time.Now()
	b, err := s.CreateBooking(model.Booking{
		RoomTypeID: roomTypeID,
		GuestName:  guest,
		GuestPhone: "13800000000",
		CheckIn:    now,
		CheckOut:   now.Add(48 * time.Hour),
		RoomCount:  roomCount,
	})
	if err != nil {
		t.Fatalf("创建订单失败: %v", err)
	}
	return b
}
