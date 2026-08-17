package service

import (
	"testing"

	"hotelbooking/internal/model"
)

func TestCancelPendingBookingReleasesRoomInventory(t *testing.T) {
	s := newTestService()
	hotel := mustCreateHotel(t, s, "湖景酒店", "杭州")
	roomType := mustCreateRoomType(t, s, hotel.ID, "亲子套房", 88800, 1)
	booking := mustCreateBooking(t, s, roomType.ID, "赵六", 1)

	cancelled, err := s.CancelBooking(booking.ID)
	if err != nil {
		t.Fatalf("pending 订单应允许取消: %v", err)
	}
	if cancelled.Status != model.BookingCancelled {
		t.Fatalf("取消后状态应为 cancelled，实际 %s", cancelled.Status)
	}
	available, err := s.AvailableRooms(roomType.ID)
	if err != nil {
		t.Fatalf("查询可用库存失败: %v", err)
	}
	if available != 1 {
		t.Fatalf("取消后应释放 1 间库存，实际可用 %d", available)
	}
	if _, err := s.CreateBooking(buildBooking(roomType.ID, "钱七", 1)); err != nil {
		t.Fatalf("取消后同一房型应可再次预订: %v", err)
	}
}
