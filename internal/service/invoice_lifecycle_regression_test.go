package service

import (
	"testing"

	"hotelbooking/internal/model"
)

func TestInvoiceOnlyAfterCompletedBookingAndKeepsBookingLink(t *testing.T) {
	s := newTestService()
	hotel := mustCreateHotel(t, s, "商务酒店", "深圳")
	roomType := mustCreateRoomType(t, s, hotel.ID, "行政房", 69900, 2)
	booking := mustCreateBooking(t, s, roomType.ID, "孙八", 1)
	if _, err := s.ConfirmBooking(booking.ID); err != nil {
		t.Fatalf("确认订单失败: %v", err)
	}
	if _, err := s.CreateInvoice(model.Invoice{BookingID: booking.ID, Title: "深圳某科技有限公司", TaxNumber: "91440300MA000001"}); err == nil {
		t.Fatal("未完成订单不应允许开票")
	}
	if _, err := s.CheckInBooking(booking.ID, "1808"); err != nil {
		t.Fatalf("办理入住失败: %v", err)
	}
	if _, err := s.CheckOutBooking(booking.ID); err != nil {
		t.Fatalf("退房完成后才能开票，退房失败: %v", err)
	}

	invoice, err := s.CreateInvoice(model.Invoice{BookingID: booking.ID, Title: "深圳某科技有限公司", TaxNumber: "91440300MA000001"})
	if err != nil {
		t.Fatalf("已完成订单开票失败: %v", err)
	}
	if invoice.BookingID != booking.ID {
		t.Fatalf("发票应保留原订单 ID，实际 %s", invoice.BookingID)
	}
	if _, err := s.CreateInvoice(model.Invoice{BookingID: booking.ID, Title: "深圳某科技有限公司", TaxNumber: "91440300MA000002"}); err == nil {
		t.Fatal("同一订单重复开票应被拒绝")
	}
}
