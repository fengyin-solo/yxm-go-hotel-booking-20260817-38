package model

import (
	"strings"
	"time"
)

const (
	InvoiceIssued = "issued"
)

// Invoice 表示订单完成后的开票记录，金额以「分」为单位。
type Invoice struct {
	ID        string    `json:"id"`
	BookingID string    `json:"booking_id"`
	Title     string    `json:"title"`
	TaxNumber string    `json:"tax_number"`
	Amount    int64     `json:"amount"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (i *Invoice) Validate() error {
	i.Title = strings.TrimSpace(i.Title)
	i.TaxNumber = strings.TrimSpace(i.TaxNumber)
	if i.BookingID == "" {
		return NewValidationError("booking_id", "必须关联订单")
	}
	if i.Title == "" {
		return NewValidationError("title", "发票抬头不能为空")
	}
	if i.TaxNumber == "" {
		return NewValidationError("tax_number", "税号不能为空")
	}
	i.TaxNumber = strings.ToUpper(i.TaxNumber)
	if i.Amount < 0 {
		return NewValidationError("amount", "开票金额不能为负数")
	}
	if i.Status == "" {
		i.Status = InvoiceIssued
	}
	if i.Status != InvoiceIssued {
		return NewValidationError("status", "发票状态不合法")
	}
	return nil
}

type InvoiceFilter struct {
	BookingID string
}

func (f InvoiceFilter) Match(i *Invoice) bool {
	if f.BookingID != "" && i.BookingID != f.BookingID {
		return false
	}
	return true
}
