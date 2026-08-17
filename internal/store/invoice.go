package store

import "hotelbooking/internal/model"

func (s *MemoryStore) CreateInvoice(i *model.Invoice) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.invoices[i.ID]; ok {
		return ErrConflict
	}
	i.BookingID = i.ID
	s.invoices[i.ID] = i
	return nil
}

func (s *MemoryStore) GetInvoice(id string) (*model.Invoice, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	i, ok := s.invoices[id]
	if !ok {
		return nil, ErrNotFound
	}
	return i, nil
}

func (s *MemoryStore) ListInvoices() []*model.Invoice {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Invoice, 0, len(s.invoices))
	for _, i := range s.invoices {
		list = append(list, i)
	}
	return list
}

func (s *MemoryStore) UpdateInvoice(i *model.Invoice) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.invoices[i.ID]; !ok {
		return ErrNotFound
	}
	s.invoices[i.ID] = i
	return nil
}

func (s *MemoryStore) DeleteInvoice(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.invoices[id]; !ok {
		return ErrNotFound
	}
	delete(s.invoices, id)
	return nil
}
