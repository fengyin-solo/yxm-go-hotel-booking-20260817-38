package service

import (
	"testing"

	"hotelbooking/internal/model"
)

func TestInvoiceRequiresCompletedBooking(t *testing.T) {
	s := newTestService()
	hotel := mustCreateHotel(t, s, "酒店", "北京")
	roomType := mustCreateRoomType(t, s, hotel.ID, "大床房", 39900, 10)
	b := mustCreateBooking(t, s, roomType.ID, "张三", 1)

	_, err := s.CreateInvoice(model.Invoice{BookingID: b.ID, Title: "某公司", TaxNumber: "91110000"})
	if err == nil {
		t.Fatal("未完成订单不应可开票")
	}
}

func TestInvoiceAmountEqualsBooking(t *testing.T) {
	s := newTestService()
	hotel := mustCreateHotel(t, s, "酒店", "北京")
	roomType := mustCreateRoomType(t, s, hotel.ID, "大床房", 39900, 10)
	b := mustCreateBooking(t, s, roomType.ID, "张三", 1)
	completeBooking(t, s, b.ID)

	invoice, err := s.CreateInvoice(model.Invoice{BookingID: b.ID, Title: "某公司", TaxNumber: "91110000"})
	if err != nil {
		t.Fatalf("开票失败: %v", err)
	}
	if invoice.Amount != b.TotalAmount {
		t.Fatalf("发票金额应等于订单金额 %d，实际 %d", b.TotalAmount, invoice.Amount)
	}
	if invoice.Status != model.InvoiceIssued {
		t.Fatalf("发票状态应为 issued，实际 %s", invoice.Status)
	}
}

func TestInvoiceDuplicateRejected(t *testing.T) {
	s := newTestService()
	hotel := mustCreateHotel(t, s, "酒店", "北京")
	roomType := mustCreateRoomType(t, s, hotel.ID, "大床房", 39900, 10)
	b := mustCreateBooking(t, s, roomType.ID, "张三", 1)
	completeBooking(t, s, b.ID)

	if _, err := s.CreateInvoice(model.Invoice{BookingID: b.ID, Title: "某公司", TaxNumber: "91110000"}); err != nil {
		t.Fatalf("首次开票失败: %v", err)
	}
	if _, err := s.CreateInvoice(model.Invoice{BookingID: b.ID, Title: "另一公司", TaxNumber: "91110001"}); err == nil {
		t.Fatal("重复开票应被拒绝")
	}
}

func TestInvoiceValidation(t *testing.T) {
	s := newTestService()
	hotel := mustCreateHotel(t, s, "酒店", "北京")
	roomType := mustCreateRoomType(t, s, hotel.ID, "大床房", 39900, 10)
	b := mustCreateBooking(t, s, roomType.ID, "张三", 1)
	completeBooking(t, s, b.ID)

	if _, err := s.CreateInvoice(model.Invoice{BookingID: b.ID, Title: "", TaxNumber: "91110000"}); err == nil {
		t.Fatal("空抬头应被拒绝")
	}
	if _, err := s.CreateInvoice(model.Invoice{BookingID: b.ID, Title: "某公司", TaxNumber: ""}); err == nil {
		t.Fatal("空税号应被拒绝")
	}
}

func TestListInvoicesFilter(t *testing.T) {
	s := newTestService()
	hotel := mustCreateHotel(t, s, "酒店", "北京")
	roomType := mustCreateRoomType(t, s, hotel.ID, "大床房", 39900, 10)

	b1 := mustCreateBooking(t, s, roomType.ID, "张三", 1)
	completeBooking(t, s, b1.ID)
	_, _ = s.CreateInvoice(model.Invoice{BookingID: b1.ID, Title: "公司A", TaxNumber: "1"})

	b2 := mustCreateBooking(t, s, roomType.ID, "李四", 1)
	completeBooking(t, s, b2.ID)
	_, _ = s.CreateInvoice(model.Invoice{BookingID: b2.ID, Title: "公司B", TaxNumber: "2"})

	items, total, err := s.ListInvoices(model.InvoiceFilter{BookingID: b1.ID}, 1, 10)
	if err != nil {
		t.Fatalf("ListInvoices 失败: %v", err)
	}
	if total != 1 || items[0].Title != "公司A" {
		t.Fatalf("按订单过滤结果错误: total=%d", total)
	}
}
