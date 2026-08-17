package store

import "hotelbooking/internal/model"

func (s *MemoryStore) CreateHotel(h *model.Hotel) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.hotels[h.ID]; ok {
		return ErrConflict
	}
	s.hotels[h.ID] = h
	return nil
}

func (s *MemoryStore) GetHotel(id string) (*model.Hotel, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	h, ok := s.hotels[id]
	if !ok {
		return nil, ErrNotFound
	}
	return h, nil
}

func (s *MemoryStore) ListHotels() []*model.Hotel {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Hotel, 0, len(s.hotels))
	for _, h := range s.hotels {
		list = append(list, h)
	}
	return list
}

func (s *MemoryStore) UpdateHotel(h *model.Hotel) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.hotels[h.ID]; !ok {
		return ErrNotFound
	}
	s.hotels[h.ID] = h
	return nil
}

func (s *MemoryStore) DeleteHotel(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.hotels[id]; !ok {
		return ErrNotFound
	}
	delete(s.hotels, id)
	return nil
}
