package service

import (
	"testing"

	"hotelbooking/internal/model"
)

func TestRevenueByHotel(t *testing.T) {
	s := newTestService()
	hotel := mustCreateHotel(t, s, "酒店", "北京")
	roomType := mustCreateRoomType(t, s, hotel.ID, "大床房", 39900, 10)

	b1 := mustCreateBooking(t, s, roomType.ID, "张三", 1)
	completeBooking(t, s, b1.ID)
	b2 := mustCreateBooking(t, s, roomType.ID, "李四", 1)
	completeBooking(t, s, b2.ID)

	// 还有一笔未完成订单，不应计入营收
	mustCreateBooking(t, s, roomType.ID, "王五", 1)

	stats, err := s.RevenueByHotel(hotel.ID)
	if err != nil {
		t.Fatalf("RevenueByHotel 失败: %v", err)
	}
	if stats.BookingCount != 2 {
		t.Fatalf("完成订单数应为 2，实际 %d", stats.BookingCount)
	}
	// 每笔 2 晚 × 39900 分 = 79800，两笔合计 159600
	if stats.TotalAmount != 159600 {
		t.Fatalf("营收应为 159600，实际 %d", stats.TotalAmount)
	}
}

func TestRevenueByHotelNotFound(t *testing.T) {
	s := newTestService()
	if _, err := s.RevenueByHotel("missing"); err == nil {
		t.Fatal("统计不存在酒店应返回错误")
	}
}

func TestOccupancyByRoomTypeNotFound(t *testing.T) {
	s := newTestService()
	if _, err := s.OccupancyByRoomType("missing"); err == nil {
		t.Fatal("统计不存在房型应返回错误")
	}
}

func TestExportReportNotFound(t *testing.T) {
	s := newTestService()
	if _, err := s.ExportHotelReport("missing"); err == nil {
		t.Fatal("导出不存在酒店报告应返回错误")
	}
}

func TestRevenueByCity(t *testing.T) {
	s := newTestService()
	hotelA := mustCreateHotel(t, s, "北京酒店", "北京")
	hotelB := mustCreateHotel(t, s, "上海酒店", "上海")
	rtA := mustCreateRoomType(t, s, hotelA.ID, "大床房", 39900, 10)
	rtB := mustCreateRoomType(t, s, hotelB.ID, "标间", 29900, 10)

	bA := mustCreateBooking(t, s, rtA.ID, "张三", 1)
	completeBooking(t, s, bA.ID)
	bB := mustCreateBooking(t, s, rtB.ID, "李四", 1)
	completeBooking(t, s, bB.ID)

	result, err := s.RevenueByCity()
	if err != nil {
		t.Fatalf("RevenueByCity 失败: %v", err)
	}
	if result["北京"] != 79800 {
		t.Fatalf("北京营收应为 79800，实际 %d", result["北京"])
	}
	if result["上海"] != 59800 {
		t.Fatalf("上海营收应为 59800，实际 %d", result["上海"])
	}
}

func TestBookingStatusStats(t *testing.T) {
	s := newTestService()
	hotel := mustCreateHotel(t, s, "酒店", "北京")
	roomType := mustCreateRoomType(t, s, hotel.ID, "大床房", 39900, 10)

	b1 := mustCreateBooking(t, s, roomType.ID, "张三", 1)
	_, _ = s.ConfirmBooking(b1.ID)
	b2 := mustCreateBooking(t, s, roomType.ID, "李四", 1)
	_, _ = s.CancelBooking(b2.ID)
	mustCreateBooking(t, s, roomType.ID, "王五", 1)

	stats := s.BookingStatusStats()
	if stats.Pending != 1 || stats.Confirmed != 1 || stats.Cancelled != 1 {
		t.Fatalf("状态分布错误: %+v", stats)
	}
}

func TestOccupancyByRoomType(t *testing.T) {
	s := newTestService()
	hotel := mustCreateHotel(t, s, "酒店", "北京")
	roomType := mustCreateRoomType(t, s, hotel.ID, "大床房", 39900, 4)

	mustCreateBooking(t, s, roomType.ID, "张三", 1)
	occ, err := s.OccupancyByRoomType(roomType.ID)
	if err != nil {
		t.Fatalf("OccupancyByRoomType 失败: %v", err)
	}
	if occ.OccupiedRooms != 1 || occ.TotalRooms != 4 {
		t.Fatalf("占用统计错误: %+v", occ)
	}
	if occ.Rate != 0.25 {
		t.Fatalf("入住率应为 0.25，实际 %v", occ.Rate)
	}
}

func TestExportHotelReport(t *testing.T) {
	s := newTestService()
	hotel := mustCreateHotel(t, s, "酒店", "北京")
	roomType := mustCreateRoomType(t, s, hotel.ID, "大床房", 39900, 4)
	b := mustCreateBooking(t, s, roomType.ID, "张三", 1)
	completeBooking(t, s, b.ID)
	_, _ = s.CreateReview(model.Review{BookingID: b.ID, Rating: 5, Content: "很棒"})

	report, err := s.ExportHotelReport(hotel.ID)
	if err != nil {
		t.Fatalf("ExportHotelReport 失败: %v", err)
	}
	if report.Hotel.ID != hotel.ID {
		t.Fatalf("报告酒店信息错误")
	}
	if len(report.RoomTypes) != 1 {
		t.Fatalf("报告房型数应为 1，实际 %d", len(report.RoomTypes))
	}
	if report.TotalRooms != 4 {
		t.Fatalf("总房间数应为 4，实际 %d", report.TotalRooms)
	}
	if report.Revenue == nil || report.Revenue.BookingCount != 1 {
		t.Fatalf("报告营收错误")
	}
	if report.Rating == nil || report.Rating.ReviewCount != 1 {
		t.Fatalf("报告评分错误")
	}
	if report.BookingStatus.Completed != 1 {
		t.Fatalf("报告订单状态错误")
	}
}
