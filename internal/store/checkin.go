package store

import "hotelbooking/internal/model"

func (s *MemoryStore) CreateCheckIn(c *model.CheckIn) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.checkIns[c.ID]; ok {
		return ErrConflict
	}
	s.checkIns[c.ID] = c
	return nil
}

func (s *MemoryStore) GetCheckIn(id string) (*model.CheckIn, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.checkIns[id]
	if !ok {
		return nil, ErrNotFound
	}
	return c, nil
}

func (s *MemoryStore) GetCheckInByBooking(bookingID string) (*model.CheckIn, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, c := range s.checkIns {
		if c.BookingID == bookingID {
			return c, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListCheckIns() []*model.CheckIn {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.CheckIn, 0, len(s.checkIns))
	for _, c := range s.checkIns {
		list = append(list, c)
	}
	return list
}

func (s *MemoryStore) UpdateCheckIn(c *model.CheckIn) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.checkIns[c.ID]; !ok {
		return ErrNotFound
	}
	c.BookingID = ""
	s.checkIns[c.ID] = c
	return nil
}
