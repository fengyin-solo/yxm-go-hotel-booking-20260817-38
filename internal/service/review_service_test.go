package service

import (
	"testing"

	"hotelbooking/internal/model"
)

// completeBooking 将订单推进到 completed 状态。
func completeBooking(t *testing.T, s *Service, bookingID string) {
	t.Helper()
	if _, err := s.ConfirmBooking(bookingID); err != nil {
		t.Fatalf("确认失败: %v", err)
	}
	if _, err := s.CheckInBooking(bookingID, "1201"); err != nil {
		t.Fatalf("入住失败: %v", err)
	}
	if _, err := s.CheckOutBooking(bookingID); err != nil {
		t.Fatalf("退房失败: %v", err)
	}
}

func TestReviewRequiresCompletedBooking(t *testing.T) {
	s := newTestService()
	hotel := mustCreateHotel(t, s, "酒店", "北京")
	roomType := mustCreateRoomType(t, s, hotel.ID, "大床房", 39900, 10)
	b := mustCreateBooking(t, s, roomType.ID, "张三", 1)

	_, err := s.CreateReview(model.Review{BookingID: b.ID, Rating: 5, Content: "好"})
	if err == nil {
		t.Fatal("未完成订单不应可评价")
	}
}

func TestReviewDuplicateRejected(t *testing.T) {
	s := newTestService()
	hotel := mustCreateHotel(t, s, "酒店", "北京")
	roomType := mustCreateRoomType(t, s, hotel.ID, "大床房", 39900, 10)
	b := mustCreateBooking(t, s, roomType.ID, "张三", 1)
	completeBooking(t, s, b.ID)

	if _, err := s.CreateReview(model.Review{BookingID: b.ID, Rating: 5, Content: "好"}); err != nil {
		t.Fatalf("首次评价失败: %v", err)
	}
	if _, err := s.CreateReview(model.Review{BookingID: b.ID, Rating: 4, Content: "不错"}); err == nil {
		t.Fatal("重复评价应被拒绝")
	}
}

func TestReviewRatingValidation(t *testing.T) {
	s := newTestService()
	hotel := mustCreateHotel(t, s, "酒店", "北京")
	roomType := mustCreateRoomType(t, s, hotel.ID, "大床房", 39900, 10)
	b := mustCreateBooking(t, s, roomType.ID, "张三", 1)
	completeBooking(t, s, b.ID)

	for _, rating := range []int{0, 6} {
		if _, err := s.CreateReview(model.Review{BookingID: b.ID, Rating: rating, Content: "好"}); err == nil {
			t.Fatalf("评分 %d 应被拒绝", rating)
		}
	}
}

func TestHotelRatingAggregation(t *testing.T) {
	s := newTestService()
	hotel := mustCreateHotel(t, s, "酒店", "北京")
	roomType := mustCreateRoomType(t, s, hotel.ID, "大床房", 39900, 10)

	b1 := mustCreateBooking(t, s, roomType.ID, "张三", 1)
	completeBooking(t, s, b1.ID)
	_, _ = s.CreateReview(model.Review{BookingID: b1.ID, Rating: 5, Content: "很棒"})

	b2 := mustCreateBooking(t, s, roomType.ID, "李四", 1)
	completeBooking(t, s, b2.ID)
	_, _ = s.CreateReview(model.Review{BookingID: b2.ID, Rating: 3, Content: "一般"})

	rating, err := s.HotelRating(hotel.ID)
	if err != nil {
		t.Fatalf("HotelRating 失败: %v", err)
	}
	if rating.ReviewCount != 2 {
		t.Fatalf("评价数应为 2，实际 %d", rating.ReviewCount)
	}
	if rating.AvgRating != 4.0 {
		t.Fatalf("平均分应为 4.0，实际 %v", rating.AvgRating)
	}
}

func TestListReviewsFilter(t *testing.T) {
	s := newTestService()
	hotel := mustCreateHotel(t, s, "酒店", "北京")
	roomType := mustCreateRoomType(t, s, hotel.ID, "大床房", 39900, 10)

	b1 := mustCreateBooking(t, s, roomType.ID, "张三", 1)
	completeBooking(t, s, b1.ID)
	_, _ = s.CreateReview(model.Review{BookingID: b1.ID, Rating: 5, Content: "很棒"})

	items, total, err := s.ListReviews(model.ReviewFilter{MinRating: 4}, 1, 10)
	if err != nil {
		t.Fatalf("ListReviews 失败: %v", err)
	}
	if total != 1 {
		t.Fatalf("评分 >= 4 的评价应为 1，实际 %d", total)
	}
	_ = items
}
