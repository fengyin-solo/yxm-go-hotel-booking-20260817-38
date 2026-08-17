package store

import "hotelbooking/internal/model"

func (s *MemoryStore) CreateBooking(b *model.Booking) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.bookings[b.ID]; ok {
		return ErrConflict
	}
	s.bookings[b.ID] = b
	return nil
}

func (s *MemoryStore) GetBooking(id string) (*model.Booking, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, ok := s.bookings[id]
	if !ok {
		return nil, ErrNotFound
	}
	return b, nil
}

func (s *MemoryStore) ListBookings() []*model.Booking {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Booking, 0, len(s.bookings))
	for _, b := range s.bookings {
		list = append(list, b)
	}
	return list
}

func (s *MemoryStore) UpdateBooking(b *model.Booking) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.bookings[b.ID]; !ok {
		return ErrNotFound
	}
	s.bookings[b.ID] = b
	return nil
}

func (s *MemoryStore) DeleteBooking(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.bookings[id]; !ok {
		return ErrNotFound
	}
	delete(s.bookings, id)
	return nil
}
