package service

import "testing"

func TestHotelReportSeparatesCapacityOccupancyAndCompletedRevenue(t *testing.T) {
	s := newTestService()
	hotel := mustCreateHotel(t, s, "中心酒店", "南京")
	standard := mustCreateRoomType(t, s, hotel.ID, "标准间", 30000, 4)
	suite := mustCreateRoomType(t, s, hotel.ID, "套房", 80000, 2)
	completed := mustCreateBooking(t, s, standard.ID, "冯一", 1)
	active := mustCreateBooking(t, s, suite.ID, "陈二", 1)
	cancelled := mustCreateBooking(t, s, standard.ID, "褚三", 1)
	completeBooking(t, s, completed.ID)
	if _, err := s.ConfirmBooking(active.ID); err != nil {
		t.Fatalf("确认入住中订单失败: %v", err)
	}
	if _, err := s.CheckInBooking(active.ID, "2202"); err != nil {
		t.Fatalf("办理入住失败: %v", err)
	}
	if _, err := s.CancelBooking(cancelled.ID); err != nil {
		t.Fatalf("取消订单失败: %v", err)
	}

	report, err := s.ExportHotelReport(hotel.ID)
	if err != nil {
		t.Fatalf("导出酒店报表失败: %v", err)
	}
	if report.TotalRooms != 6 || report.OccupiedRooms != 1 {
		t.Fatalf("房间总数/占用数错误: total=%d occupied=%d", report.TotalRooms, report.OccupiedRooms)
	}
	if report.BookingStatus.Completed != 1 || report.BookingStatus.CheckedIn != 1 || report.BookingStatus.Cancelled != 1 {
		t.Fatalf("状态分布错误: %+v", report.BookingStatus)
	}
	if report.Revenue.BookingCount != 1 || report.Revenue.TotalAmount != completed.TotalAmount {
		t.Fatalf("营收应只统计已完成订单: count=%d amount=%d", report.Revenue.BookingCount, report.Revenue.TotalAmount)
	}
}
