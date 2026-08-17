package service

import (
	"testing"

	"hotelbooking/internal/model"
)

func TestHotelRatingCountsOnlyReviewsForCompletedBookings(t *testing.T) {
	s := newTestService()
	hotel := mustCreateHotel(t, s, "精品酒店", "成都")
	roomType := mustCreateRoomType(t, s, hotel.ID, "园景房", 52000, 3)
	first := mustCreateBooking(t, s, roomType.ID, "吴十", 1)
	second := mustCreateBooking(t, s, roomType.ID, "郑十一", 1)
	completeBooking(t, s, first.ID)
	completeBooking(t, s, second.ID)
	if _, err := s.CreateReview(model.Review{BookingID: first.ID, Rating: 5, Content: "服务很好"}); err != nil {
		t.Fatalf("创建第一条评价失败: %v", err)
	}
	if _, err := s.CreateReview(model.Review{BookingID: second.ID, Rating: 3, Content: "整体还可以"}); err != nil {
		t.Fatalf("创建第二条评价失败: %v", err)
	}

	rating, err := s.HotelRating(hotel.ID)
	if err != nil {
		t.Fatalf("计算酒店评分失败: %v", err)
	}
	if rating.ReviewCount != 2 {
		t.Fatalf("评价数应为 2，实际 %d", rating.ReviewCount)
	}
	if rating.AvgRating != 4 {
		t.Fatalf("平均评分应为 4，实际 %.2f", rating.AvgRating)
	}
}
