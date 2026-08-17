package store

import "hotelbooking/internal/model"

func (s *MemoryStore) CreateRoomType(r *model.RoomType) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.roomTypes[r.ID]; ok {
		return ErrConflict
	}
	s.roomTypes[r.ID] = r
	return nil
}

func (s *MemoryStore) GetRoomType(id string) (*model.RoomType, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.roomTypes[id]
	if !ok {
		return nil, ErrNotFound
	}
	return r, nil
}

func (s *MemoryStore) ListRoomTypes() []*model.RoomType {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.RoomType, 0, len(s.roomTypes))
	for _, r := range s.roomTypes {
		list = append(list, r)
	}
	return list
}

func (s *MemoryStore) UpdateRoomType(r *model.RoomType) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.roomTypes[r.ID]; !ok {
		return ErrNotFound
	}
	s.roomTypes[r.ID] = r
	return nil
}

func (s *MemoryStore) DeleteRoomType(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.roomTypes[id]; !ok {
		return ErrNotFound
	}
	delete(s.roomTypes, id)
	return nil
}
