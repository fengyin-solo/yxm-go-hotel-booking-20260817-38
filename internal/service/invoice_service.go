package service

import (
	"sort"
	"time"

	"hotelbooking/internal/model"
	"hotelbooking/pkg/idgen"
)

func (s *Service) CreateInvoice(invoice model.Invoice) (*model.Invoice, error) {
	if err := invoice.Validate(); err != nil {
		return nil, err
	}
	booking, err := s.store.GetBooking(invoice.BookingID)
	if err != nil {
		return nil, model.NewValidationError("booking_id", "关联的订单不存在")
	}
	if booking.Status != model.BookingCompleted {
		return nil, model.NewValidationError("booking_id", "仅已完成订单可开票")
	}
	for _, i := range s.store.ListInvoices() {
		if i.BookingID == invoice.BookingID {
			return nil, model.NewValidationError("booking_id", "该订单已开过发票")
		}
	}
	now := time.Now()
	invoice.ID = idgen.Hex()
	invoice.Amount = booking.TotalAmount
	invoice.Status = model.InvoiceIssued
	invoice.CreatedAt = now
	invoice.UpdatedAt = now
	if err := s.store.CreateInvoice(&invoice); err != nil {
		return nil, err
	}
	return &invoice, nil
}

func (s *Service) GetInvoice(id string) (*model.Invoice, error) {
	return s.store.GetInvoice(id)
}

func (s *Service) ListInvoices(filter model.InvoiceFilter, page, size int) ([]*model.Invoice, int, error) {
	all := s.store.ListInvoices()
	matched := make([]*model.Invoice, 0, len(all))
	for _, i := range all {
		if filter.Match(i) {
			matched = append(matched, i)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Invoice{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) DeleteInvoice(id string) error {
	return s.store.DeleteInvoice(id)
}
