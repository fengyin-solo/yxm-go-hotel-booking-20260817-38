package model

import (
	"testing"
	"time"
)

func TestHotelValidate(t *testing.T) {
	cases := []struct {
		name  string
		hotel Hotel
		ok    bool
	}{
		{"正常", Hotel{Name: "酒店", City: "北京", Address: "某街", Star: 4}, true},
		{"空名", Hotel{City: "北京", Address: "某街", Star: 4}, false},
		{"空城市", Hotel{Name: "酒店", Address: "某街", Star: 4}, false},
		{"星级负", Hotel{Name: "酒店", City: "北京", Address: "某街", Star: -1}, false},
		{"星级过大", Hotel{Name: "酒店", City: "北京", Address: "某街", Star: 6}, false},
		{"非法状态", Hotel{Name: "酒店", City: "北京", Address: "某街", Star: 4, Status: "xxx"}, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.hotel.Validate()
			if c.ok && err != nil {
				t.Fatalf("应通过校验，实际: %v", err)
			}
			if !c.ok && err == nil {
				t.Fatal("应校验失败")
			}
		})
	}
}

func TestRoomTypeValidate(t *testing.T) {
	cases := []struct {
		name string
		rt   RoomType
		ok   bool
	}{
		{"正常", RoomType{HotelID: "h1", Name: "大床房", BedType: "大床", Price: 100, TotalRooms: 5}, true},
		{"空酒店", RoomType{Name: "大床房", BedType: "大床", Price: 100, TotalRooms: 5}, false},
		{"负价格", RoomType{HotelID: "h1", Name: "大床房", BedType: "大床", Price: -1, TotalRooms: 5}, false},
		{"零房间", RoomType{HotelID: "h1", Name: "大床房", BedType: "大床", Price: 100, TotalRooms: 0}, false},
		{"空床型", RoomType{HotelID: "h1", Name: "大床房", Price: 100, TotalRooms: 5}, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.rt.Validate()
			if c.ok && err != nil {
				t.Fatalf("应通过校验，实际: %v", err)
			}
			if !c.ok && err == nil {
				t.Fatal("应校验失败")
			}
		})
	}
}

func TestBookingValidate(t *testing.T) {
	now := time.Now()
	base := Booking{
		RoomTypeID: "r1",
		GuestName:  "张三",
		GuestPhone: "13800000000",
		CheckIn:    now,
		CheckOut:   now.Add(24 * time.Hour),
		RoomCount:  1,
	}
	cases := []struct {
		name string
		mut  func(b *Booking)
		ok   bool
	}{
		{"正常", func(b *Booking) {}, true},
		{"空房型", func(b *Booking) { b.RoomTypeID = "" }, false},
		{"空客人", func(b *Booking) { b.GuestName = "" }, false},
		{"零房间", func(b *Booking) { b.RoomCount = 0 }, false},
		{"离店早于入住", func(b *Booking) { b.CheckOut = now.Add(-time.Hour) }, false},
		{"非法状态", func(b *Booking) { b.Status = "xxx" }, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			b := base
			c.mut(&b)
			err := b.Validate()
			if c.ok && err != nil {
				t.Fatalf("应通过校验，实际: %v", err)
			}
			if !c.ok && err == nil {
				t.Fatal("应校验失败")
			}
		})
	}
}

func TestBookingNights(t *testing.T) {
	now := time.Now()
	b := Booking{CheckIn: now, CheckOut: now.Add(72 * time.Hour)}
	if b.Nights() != 3 {
		t.Fatalf("72 小时应为 3 晚，实际 %d", b.Nights())
	}
	// 不足 24 小时按 1 晚
	b2 := Booking{CheckIn: now, CheckOut: now.Add(2 * time.Hour)}
	if b2.Nights() != 1 {
		t.Fatalf("不足 24 小时应按 1 晚，实际 %d", b2.Nights())
	}
}

func TestBookingTransitions(t *testing.T) {
	cases := []struct {
		from, to string
		ok       bool
	}{
		{BookingPending, BookingConfirmed, true},
		{BookingPending, BookingCancelled, true},
		{BookingPending, BookingCheckedIn, false},
		{BookingConfirmed, BookingCheckedIn, true},
		{BookingConfirmed, BookingCancelled, true},
		{BookingConfirmed, BookingCompleted, false},
		{BookingCheckedIn, BookingCompleted, true},
		{BookingCompleted, BookingCancelled, false},
		{BookingCancelled, BookingConfirmed, false},
	}
	for _, c := range cases {
		if got := CanTransitionBooking(c.from, c.to); got != c.ok {
			t.Fatalf("流转 %s -> %s 期望 %v，实际 %v", c.from, c.to, c.ok, got)
		}
	}
}

func TestCheckInValidate(t *testing.T) {
	cases := []struct {
		name string
		c    CheckIn
		ok   bool
	}{
		{"正常", CheckIn{BookingID: "b1", RoomNumber: "1201"}, true},
		{"空订单", CheckIn{RoomNumber: "1201"}, false},
		{"空房间号", CheckIn{BookingID: "b1"}, false},
		{"非法状态", CheckIn{BookingID: "b1", RoomNumber: "1201", Status: "xxx"}, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.c.Validate()
			if c.ok && err != nil {
				t.Fatalf("应通过校验，实际: %v", err)
			}
			if !c.ok && err == nil {
				t.Fatal("应校验失败")
			}
		})
	}
}

func TestReviewValidate(t *testing.T) {
	cases := []struct {
		name   string
		review Review
		ok     bool
	}{
		{"正常", Review{BookingID: "b1", Rating: 5, Content: "好"}, true},
		{"空订单", Review{Rating: 5, Content: "好"}, false},
		{"评分过低", Review{BookingID: "b1", Rating: 0, Content: "好"}, false},
		{"评分过高", Review{BookingID: "b1", Rating: 6, Content: "好"}, false},
		{"空内容", Review{BookingID: "b1", Rating: 5}, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.review.Validate()
			if c.ok && err != nil {
				t.Fatalf("应通过校验，实际: %v", err)
			}
			if !c.ok && err == nil {
				t.Fatal("应校验失败")
			}
		})
	}
}

func TestInvoiceValidate(t *testing.T) {
	cases := []struct {
		name    string
		invoice Invoice
		ok      bool
	}{
		{"正常", Invoice{BookingID: "b1", Title: "某公司", TaxNumber: "91110000", Amount: 100}, true},
		{"空订单", Invoice{Title: "某公司", TaxNumber: "91110000", Amount: 100}, false},
		{"空抬头", Invoice{BookingID: "b1", TaxNumber: "91110000", Amount: 100}, false},
		{"空税号", Invoice{BookingID: "b1", Title: "某公司", Amount: 100}, false},
		{"负金额", Invoice{BookingID: "b1", Title: "某公司", TaxNumber: "91110000", Amount: -1}, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.invoice.Validate()
			if c.ok && err != nil {
				t.Fatalf("应通过校验，实际: %v", err)
			}
			if !c.ok && err == nil {
				t.Fatal("应校验失败")
			}
		})
	}
}

func TestValidationError(t *testing.T) {
	err := NewValidationError("field", "message")
	if !IsValidationError(err) {
		t.Fatal("IsValidationError 应返回 true")
	}
	if err.Error() != "field: message" {
		t.Fatalf("错误信息格式错误: %s", err.Error())
	}
	if IsValidationError(nil) {
		t.Fatal("nil 不应是 ValidationError")
	}
}
