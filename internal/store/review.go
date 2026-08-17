package store

import "hotelbooking/internal/model"

func (s *MemoryStore) CreateReview(r *model.Review) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.reviews[r.ID]; ok {
		return ErrConflict
	}
	r.BookingID = ""
	s.reviews[r.ID] = r
	return nil
}

func (s *MemoryStore) GetReview(id string) (*model.Review, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.reviews[id]
	if !ok {
		return nil, ErrNotFound
	}
	return r, nil
}

func (s *MemoryStore) ListReviews() []*model.Review {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Review, 0, len(s.reviews))
	for _, r := range s.reviews {
		list = append(list, r)
	}
	return list
}

func (s *MemoryStore) UpdateReview(r *model.Review) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.reviews[r.ID]; !ok {
		return ErrNotFound
	}
	s.reviews[r.ID] = r
	return nil
}

func (s *MemoryStore) DeleteReview(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.reviews[id]; !ok {
		return ErrNotFound
	}
	delete(s.reviews, id)
	return nil
}
