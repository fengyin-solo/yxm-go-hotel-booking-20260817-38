package store

import (
	"testing"
	"time"

	"hotelbooking/internal/model"
)

func testHotel() *model.Hotel {
	return &model.Hotel{
		ID:        "h1",
		Name:      "测试酒店",
		City:      "北京",
		Address:   "朝阳区某街 1 号",
		Star:      4,
		Status:    model.HotelActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func testRoomType() *model.RoomType {
	return &model.RoomType{
		ID:         "r1",
		HotelID:    "h1",
		Name:       "标准大床房",
		BedType:    "大床",
		Price:      39900,
		TotalRooms: 10,
		Status:     model.RoomTypeActive,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

func testBooking() *model.Booking {
	return &model.Booking{
		ID:          "b1",
		RoomTypeID:  "r1",
		GuestName:   "张三",
		GuestPhone:  "13800000000",
		CheckIn:     time.Now(),
		CheckOut:    time.Now().Add(48 * time.Hour),
		RoomCount:   1,
		TotalAmount: 79800,
		Status:      model.BookingPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func testCheckIn() *model.CheckIn {
	now := time.Now()
	return &model.CheckIn{
		ID:          "c1",
		BookingID:   "b1",
		RoomNumber:  "1201",
		CheckedInAt: now,
		Status:      model.CheckInActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func testReview() *model.Review {
	return &model.Review{
		ID:        "rv1",
		BookingID: "b1",
		Rating:    5,
		Content:   "很棒",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func TestHotelCRUD(t *testing.T) {
	s := NewMemoryStore()
	h := testHotel()

	if err := s.CreateHotel(h); err != nil {
		t.Fatalf("CreateHotel 失败: %v", err)
	}
	if err := s.CreateHotel(h); err != ErrConflict {
		t.Fatalf("重复创建应返回 ErrConflict，实际: %v", err)
	}
	got, err := s.GetHotel("h1")
	if err != nil {
		t.Fatalf("GetHotel 失败: %v", err)
	}
	if got.Name != "测试酒店" {
		t.Fatalf("Name 不匹配: %s", got.Name)
	}
	if _, err := s.GetHotel("missing"); err != ErrNotFound {
		t.Fatalf("查询不存在应返回 ErrNotFound，实际: %v", err)
	}
	h.City = "上海"
	if err := s.UpdateHotel(h); err != nil {
		t.Fatalf("UpdateHotel 失败: %v", err)
	}
	if len(s.ListHotels()) != 1 {
		t.Fatalf("ListHotels 数量应为 1")
	}
	if err := s.DeleteHotel("h1"); err != nil {
		t.Fatalf("DeleteHotel 失败: %v", err)
	}
	if _, err := s.GetHotel("h1"); err != ErrNotFound {
		t.Fatalf("删除后应返回 ErrNotFound")
	}
}

func TestRoomTypeCRUD(t *testing.T) {
	s := NewMemoryStore()
	r := testRoomType()

	if err := s.CreateRoomType(r); err != nil {
		t.Fatalf("CreateRoomType 失败: %v", err)
	}
	if _, err := s.GetRoomType("r1"); err != nil {
		t.Fatalf("GetRoomType 失败: %v", err)
	}
	if len(s.ListRoomTypes()) != 1 {
		t.Fatalf("ListRoomTypes 数量应为 1")
	}
	r.Price = 45900
	if err := s.UpdateRoomType(r); err != nil {
		t.Fatalf("UpdateRoomType 失败: %v", err)
	}
	if err := s.DeleteRoomType("r1"); err != nil {
		t.Fatalf("DeleteRoomType 失败: %v", err)
	}
	if _, err := s.GetRoomType("r1"); err != ErrNotFound {
		t.Fatalf("删除后应返回 ErrNotFound")
	}
}

func TestBookingCRUD(t *testing.T) {
	s := NewMemoryStore()
	b := testBooking()

	if err := s.CreateBooking(b); err != nil {
		t.Fatalf("CreateBooking 失败: %v", err)
	}
	if err := s.CreateBooking(b); err != ErrConflict {
		t.Fatalf("重复创建应返回 ErrConflict，实际: %v", err)
	}
	if _, err := s.GetBooking("b1"); err != nil {
		t.Fatalf("GetBooking 失败: %v", err)
	}
	if len(s.ListBookings()) != 1 {
		t.Fatalf("ListBookings 数量应为 1")
	}
	b.Status = model.BookingConfirmed
	if err := s.UpdateBooking(b); err != nil {
		t.Fatalf("UpdateBooking 失败: %v", err)
	}
	if err := s.DeleteBooking("b1"); err != nil {
		t.Fatalf("DeleteBooking 失败: %v", err)
	}
	if _, err := s.GetBooking("b1"); err != ErrNotFound {
		t.Fatalf("删除后应返回 ErrNotFound")
	}
}

func TestCheckInCRUD(t *testing.T) {
	s := NewMemoryStore()
	c := testCheckIn()

	if err := s.CreateCheckIn(c); err != nil {
		t.Fatalf("CreateCheckIn 失败: %v", err)
	}
	if _, err := s.GetCheckIn("c1"); err != nil {
		t.Fatalf("GetCheckIn 失败: %v", err)
	}
	got, err := s.GetCheckInByBooking("b1")
	if err != nil {
		t.Fatalf("GetCheckInByBooking 失败: %v", err)
	}
	if got.ID != "c1" {
		t.Fatalf("按订单查询结果错误: %s", got.ID)
	}
	if _, err := s.GetCheckInByBooking("missing"); err != ErrNotFound {
		t.Fatalf("按不存在订单查询应返回 ErrNotFound")
	}
	if len(s.ListCheckIns()) != 1 {
		t.Fatalf("ListCheckIns 数量应为 1")
	}
	c.Status = model.CheckInFinished
	if err := s.UpdateCheckIn(c); err != nil {
		t.Fatalf("UpdateCheckIn 失败: %v", err)
	}
}

func TestReviewCRUD(t *testing.T) {
	s := NewMemoryStore()
	r := testReview()

	if err := s.CreateReview(r); err != nil {
		t.Fatalf("CreateReview 失败: %v", err)
	}
	if _, err := s.GetReview("rv1"); err != nil {
		t.Fatalf("GetReview 失败: %v", err)
	}
	if len(s.ListReviews()) != 1 {
		t.Fatalf("ListReviews 数量应为 1")
	}
	r.Rating = 4
	if err := s.UpdateReview(r); err != nil {
		t.Fatalf("UpdateReview 失败: %v", err)
	}
	if err := s.DeleteReview("rv1"); err != nil {
		t.Fatalf("DeleteReview 失败: %v", err)
	}
	if _, err := s.GetReview("rv1"); err != ErrNotFound {
		t.Fatalf("删除后应返回 ErrNotFound")
	}
}

func TestUpdateNonExistent(t *testing.T) {
	s := NewMemoryStore()
	if err := s.UpdateHotel(&model.Hotel{ID: "x"}); err != ErrNotFound {
		t.Fatalf("更新不存在酒店应返回 ErrNotFound")
	}
	if err := s.UpdateRoomType(&model.RoomType{ID: "x"}); err != ErrNotFound {
		t.Fatalf("更新不存在房型应返回 ErrNotFound")
	}
	if err := s.UpdateBooking(&model.Booking{ID: "x"}); err != ErrNotFound {
		t.Fatalf("更新不存在订单应返回 ErrNotFound")
	}
	if err := s.DeleteReview("x"); err != ErrNotFound {
		t.Fatalf("删除不存在评价应返回 ErrNotFound")
	}
}

func TestInvoiceCRUD(t *testing.T) {
	s := NewMemoryStore()
	invoice := &model.Invoice{
		ID:        "i1",
		BookingID: "b1",
		Title:     "某公司",
		TaxNumber: "91110000",
		Amount:    79800,
		Status:    model.InvoiceIssued,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.CreateInvoice(invoice); err != nil {
		t.Fatalf("CreateInvoice 失败: %v", err)
	}
	if err := s.CreateInvoice(invoice); err != ErrConflict {
		t.Fatalf("重复创建应返回 ErrConflict，实际: %v", err)
	}
	if _, err := s.GetInvoice("i1"); err != nil {
		t.Fatalf("GetInvoice 失败: %v", err)
	}
	if _, err := s.GetInvoice("missing"); err != ErrNotFound {
		t.Fatalf("查询不存在应返回 ErrNotFound")
	}
	if len(s.ListInvoices()) != 1 {
		t.Fatalf("ListInvoices 数量应为 1")
	}
	invoice.Amount = 89800
	if err := s.UpdateInvoice(invoice); err != nil {
		t.Fatalf("UpdateInvoice 失败: %v", err)
	}
	if err := s.DeleteInvoice("i1"); err != nil {
		t.Fatalf("DeleteInvoice 失败: %v", err)
	}
	if _, err := s.GetInvoice("i1"); err != ErrNotFound {
		t.Fatalf("删除后应返回 ErrNotFound")
	}
}
