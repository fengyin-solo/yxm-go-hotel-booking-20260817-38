package handler

import (
	"net/http"

	"hotelbooking/internal/model"
	"hotelbooking/pkg/httpx"
)

func (s *Server) registerInvoiceRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/invoices", s.createInvoice)
	mux.HandleFunc("GET /api/invoices", s.listInvoices)
	mux.HandleFunc("GET /api/invoices/{id}", s.getInvoice)
	mux.HandleFunc("DELETE /api/invoices/{id}", s.deleteInvoice)
}

type invoiceRequest struct {
	BookingID string `json:"booking_id"`
	Title     string `json:"title"`
	TaxNumber string `json:"tax_number"`
}

func (s *Server) createInvoice(w http.ResponseWriter, r *http.Request) {
	var req invoiceRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	invoice, err := s.svc.CreateInvoice(model.Invoice{
		BookingID: req.BookingID,
		Title:     req.Title,
		TaxNumber: req.TaxNumber,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, invoice)
}

func (s *Server) listInvoices(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.InvoiceFilter{
		BookingID: r.URL.Query().Get("booking_id"),
	}
	items, total, err := s.svc.ListInvoices(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getInvoice(w http.ResponseWriter, r *http.Request) {
	invoice, err := s.svc.GetInvoice(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, invoice)
}

func (s *Server) deleteInvoice(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteInvoice(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
