package service

import (
	"testing"
	"time"

	"hotelbooking/internal/model"
)

func TestCreateBookingComputesAmount(t *testing.T) {
	s := newTestService()
	hotel := mustCreateHotel(t, s, "酒店", "北京")
	roomType := mustCreateRoomType(t, s, hotel.ID, "大床房", 39900, 10)

	now := time.Now()
	b, err := s.CreateBooking(model.Booking{
		RoomTypeID: roomType.ID,
		GuestName:  "张三",
		GuestPhone: "13800000000",
		CheckIn:    now,
		CheckOut:   now.Add(72 * time.Hour),
		RoomCount:  2,
	})
	if err != nil {
		t.Fatalf("创建订单失败: %v", err)
	}
	// 3 晚 × 2 间 × 39900 分 = 239400 分
	if b.TotalAmount != 239400 {
		t.Fatalf("金额计算错误: 期望 239400，实际 %d", b.TotalAmount)
	}
	if b.Status != model.BookingPending {
		t.Fatalf("新订单应为 pending，实际 %s", b.Status)
	}
}

func TestCreateBookingRejectsBadDateRange(t *testing.T) {
	s := newTestService()
	hotel := mustCreateHotel(t, s, "酒店", "北京")
	roomType := mustCreateRoomType(t, s, hotel.ID, "大床房", 39900, 10)

	now := time.Now()
	_, err := s.CreateBooking(model.Booking{
		RoomTypeID: roomType.ID,
		GuestName:  "张三",
		GuestPhone: "13800000000",
		CheckIn:    now,
		CheckOut:   now.Add(-24 * time.Hour),
		RoomCount:  1,
	})
	if err == nil {
		t.Fatal("离店早于入住应返回错误")
	}
}

func TestCreateBookingRejectsInsufficientStock(t *testing.T) {
	s := newTestService()
	hotel := mustCreateHotel(t, s, "酒店", "北京")
	roomType := mustCreateRoomType(t, s, hotel.ID, "大床房", 39900, 2)

	mustCreateBooking(t, s, roomType.ID, "张三", 1)
	mustCreateBooking(t, s, roomType.ID, "李四", 1)

	if _, err := s.CreateBooking(buildBooking(roomType.ID, "王五", 1)); err == nil {
		t.Fatal("库存不足应返回错误")
	}
}

func buildBooking(roomTypeID, guest string, count int) model.Booking {
	now := time.Now()
	return model.Booking{
		RoomTypeID: roomTypeID,
		GuestName:  guest,
		GuestPhone: "13800000000",
		CheckIn:    now,
		CheckOut:   now.Add(48 * time.Hour),
		RoomCount:  count,
	}
}

func TestBookingLifecycle(t *testing.T) {
	s := newTestService()
	hotel := mustCreateHotel(t, s, "酒店", "北京")
	roomType := mustCreateRoomType(t, s, hotel.ID, "大床房", 39900, 10)
	b := mustCreateBooking(t, s, roomType.ID, "张三", 1)

	confirmed, err := s.ConfirmBooking(b.ID)
	if err != nil {
		t.Fatalf("确认失败: %v", err)
	}
	if confirmed.Status != model.BookingConfirmed {
		t.Fatalf("确认后状态应为 confirmed，实际 %s", confirmed.Status)
	}

	checkedIn, err := s.CheckInBooking(b.ID, "1201")
	if err != nil {
		t.Fatalf("入住失败: %v", err)
	}
	if checkedIn.Status != model.BookingCheckedIn {
		t.Fatalf("入住后状态应为 checked_in，实际 %s", checkedIn.Status)
	}

	completed, err := s.CheckOutBooking(b.ID)
	if err != nil {
		t.Fatalf("退房失败: %v", err)
	}
	if completed.Status != model.BookingCompleted {
		t.Fatalf("退房后状态应为 completed，实际 %s", completed.Status)
	}
}

func TestBookingInvalidTransition(t *testing.T) {
	s := newTestService()
	hotel := mustCreateHotel(t, s, "酒店", "北京")
	roomType := mustCreateRoomType(t, s, hotel.ID, "大床房", 39900, 10)
	b := mustCreateBooking(t, s, roomType.ID, "张三", 1)

	// pending 不能直接入住
	if _, err := s.CheckInBooking(b.ID, "1201"); err == nil {
		t.Fatal("pending 状态不能办理入住")
	}
	// pending 不能退房
	if _, err := s.CheckOutBooking(b.ID); err == nil {
		t.Fatal("pending 状态不能退房")
	}
}

func TestCancelBookingFreesStock(t *testing.T) {
	s := newTestService()
	hotel := mustCreateHotel(t, s, "酒店", "北京")
	roomType := mustCreateRoomType(t, s, hotel.ID, "大床房", 39900, 1)

	b := mustCreateBooking(t, s, roomType.ID, "张三", 1)
	if _, err := s.CancelBooking(b.ID); err != nil {
		t.Fatalf("取消失败: %v", err)
	}
	// 取消后库存释放，可再次预订
	if _, err := s.CreateBooking(buildBooking(roomType.ID, "李四", 1)); err != nil {
		t.Fatalf("取消后应可再次预订，实际: %v", err)
	}
}

func TestCompletedBookingCannotBeCancelled(t *testing.T) {
	s := newTestService()
	hotel := mustCreateHotel(t, s, "酒店", "北京")
	roomType := mustCreateRoomType(t, s, hotel.ID, "大床房", 39900, 10)
	b := mustCreateBooking(t, s, roomType.ID, "张三", 1)
	_, _ = s.ConfirmBooking(b.ID)
	_, _ = s.CheckInBooking(b.ID, "1201")
	_, _ = s.CheckOutBooking(b.ID)

	if _, err := s.CancelBooking(b.ID); err == nil {
		t.Fatal("已完成订单不能取消")
	}
}

func TestBatchConfirm(t *testing.T) {
	s := newTestService()
	hotel := mustCreateHotel(t, s, "酒店", "北京")
	roomType := mustCreateRoomType(t, s, hotel.ID, "大床房", 39900, 10)
	b1 := mustCreateBooking(t, s, roomType.ID, "张三", 1)
	b2 := mustCreateBooking(t, s, roomType.ID, "李四", 1)

	result := s.BatchConfirm([]string{b1.ID, b2.ID, "missing"})
	if len(result.Succeeded) != 2 {
		t.Fatalf("成功确认应为 2，实际 %d", len(result.Succeeded))
	}
	if _, ok := result.Failed["missing"]; !ok {
		t.Fatalf("不存在的订单应记录到 Failed")
	}
}

func TestUpdateBookingOnlyPending(t *testing.T) {
	s := newTestService()
	hotel := mustCreateHotel(t, s, "酒店", "北京")
	roomType := mustCreateRoomType(t, s, hotel.ID, "大床房", 39900, 10)
	b := mustCreateBooking(t, s, roomType.ID, "张三", 1)

	if _, err := s.UpdateBooking(b.ID, model.Booking{GuestName: "张三丰"}); err != nil {
		t.Fatalf("pending 订单应可修改: %v", err)
	}
	_, _ = s.ConfirmBooking(b.ID)
	if _, err := s.UpdateBooking(b.ID, model.Booking{GuestName: "李四"}); err == nil {
		t.Fatal("已确认订单不应可修改")
	}
}

func TestListBookingsFilter(t *testing.T) {
	s := newTestService()
	hotel := mustCreateHotel(t, s, "酒店", "北京")
	roomType := mustCreateRoomType(t, s, hotel.ID, "大床房", 39900, 10)
	b1 := mustCreateBooking(t, s, roomType.ID, "张三", 1)
	_ = b1
	mustCreateBooking(t, s, roomType.ID, "李四", 1)

	items, total, err := s.ListBookings(model.BookingFilter{Status: model.BookingPending}, 1, 10)
	if err != nil {
		t.Fatalf("ListBookings 失败: %v", err)
	}
	if total != 2 {
		t.Fatalf("pending 订单应为 2，实际 %d", total)
	}
	_ = items
}

func TestListBookingsDateRangeFilter(t *testing.T) {
	s := newTestService()
	hotel := mustCreateHotel(t, s, "酒店", "北京")
	roomType := mustCreateRoomType(t, s, hotel.ID, "大床房", 39900, 10)
	mustCreateBooking(t, s, roomType.ID, "张三", 1)

	now := time.Now()
	items, total, err := s.ListBookings(model.BookingFilter{
		From: now.Add(-24 * time.Hour),
		To:   now.Add(24 * time.Hour),
	}, 1, 10)
	if err != nil {
		t.Fatalf("ListBookings 失败: %v", err)
	}
	if total != 1 {
		t.Fatalf("日期范围内订单应为 1，实际 %d", total)
	}

	// 范围外的日期不应匹配
	items, total, _ = s.ListBookings(model.BookingFilter{
		From: now.Add(10 * 24 * time.Hour),
		To:   now.Add(20 * 24 * time.Hour),
	}, 1, 10)
	if total != 0 {
		t.Fatalf("日期范围外订单应为 0，实际 %d", total)
	}
	_ = items
}

func TestListBookingsPagination(t *testing.T) {
	s := newTestService()
	hotel := mustCreateHotel(t, s, "酒店", "北京")
	roomType := mustCreateRoomType(t, s, hotel.ID, "大床房", 39900, 20)
	for i := 0; i < 5; i++ {
		mustCreateBooking(t, s, roomType.ID, "客人"+string(rune('A'+i)), 1)
	}

	items, total, _ := s.ListBookings(model.BookingFilter{}, 1, 2)
	if total != 5 || len(items) != 2 {
		t.Fatalf("第一页应为 2 条，实际 total=%d len=%d", total, len(items))
	}
	items, _, _ = s.ListBookings(model.BookingFilter{}, 3, 2)
	if len(items) != 1 {
		t.Fatalf("第三页应为 1 条，实际 %d", len(items))
	}
}
