package service

import (
	"testing"

	"hotelbooking/internal/model"
)

func TestCheckInRecordLifecycle(t *testing.T) {
	s := newTestService()
	hotel := mustCreateHotel(t, s, "酒店", "北京")
	roomType := mustCreateRoomType(t, s, hotel.ID, "大床房", 39900, 10)
	b := mustCreateBooking(t, s, roomType.ID, "张三", 1)
	_, _ = s.ConfirmBooking(b.ID)

	if _, err := s.CheckInBooking(b.ID, "1201"); err != nil {
		t.Fatalf("入住失败: %v", err)
	}
	checkIn, err := s.GetCheckInByBooking(b.ID)
	if err != nil {
		t.Fatalf("查询入住记录失败: %v", err)
	}
	if checkIn.RoomNumber != "1201" || checkIn.Status != model.CheckInActive {
		t.Fatalf("入住记录错误: %+v", checkIn)
	}

	if _, err := s.CheckOutBooking(b.ID); err != nil {
		t.Fatalf("退房失败: %v", err)
	}
	checkIn, err = s.GetCheckInByBooking(b.ID)
	if err != nil {
		t.Fatalf("退房后查询失败: %v", err)
	}
	if checkIn.Status != model.CheckInFinished || checkIn.CheckedOutAt == nil {
		t.Fatalf("退房记录错误: %+v", checkIn)
	}
}

func TestCheckInRequiresRoomNumber(t *testing.T) {
	s := newTestService()
	hotel := mustCreateHotel(t, s, "酒店", "北京")
	roomType := mustCreateRoomType(t, s, hotel.ID, "大床房", 39900, 10)
	b := mustCreateBooking(t, s, roomType.ID, "张三", 1)
	_, _ = s.ConfirmBooking(b.ID)

	if _, err := s.CheckInBooking(b.ID, ""); err == nil {
		t.Fatal("空房间号应被拒绝")
	}
}

func TestListCheckInsFilter(t *testing.T) {
	s := newTestService()
	hotel := mustCreateHotel(t, s, "酒店", "北京")
	roomType := mustCreateRoomType(t, s, hotel.ID, "大床房", 39900, 10)

	b1 := mustCreateBooking(t, s, roomType.ID, "张三", 1)
	_, _ = s.ConfirmBooking(b1.ID)
	_, _ = s.CheckInBooking(b1.ID, "1201")

	b2 := mustCreateBooking(t, s, roomType.ID, "李四", 1)
	_, _ = s.ConfirmBooking(b2.ID)
	_, _ = s.CheckInBooking(b2.ID, "1202")

	items, total, err := s.ListCheckIns("", model.CheckInActive, 1, 10)
	if err != nil {
		t.Fatalf("ListCheckIns 失败: %v", err)
	}
	if total != 2 {
		t.Fatalf("入住中的记录应为 2，实际 %d", total)
	}

	items, total, _ = s.ListCheckIns(b1.ID, "", 1, 10)
	if total != 1 || items[0].RoomNumber != "1201" {
		t.Fatalf("按订单过滤结果错误: total=%d", total)
	}
}
