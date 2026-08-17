package service

import (
	"sort"
	"time"

	"hotelbooking/internal/model"
	"hotelbooking/pkg/idgen"
)

func (s *Service) CreateHotel(hotel model.Hotel) (*model.Hotel, error) {
	if err := hotel.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	hotel.ID = idgen.Hex()
	hotel.CreatedAt = now
	hotel.UpdatedAt = now
	if err := s.store.CreateHotel(&hotel); err != nil {
		return nil, err
	}
	return &hotel, nil
}

func (s *Service) GetHotel(id string) (*model.Hotel, error) {
	return s.store.GetHotel(id)
}

func (s *Service) ListHotels(filter model.HotelFilter, page, size int) ([]*model.Hotel, int, error) {
	all := s.store.ListHotels()
	matched := make([]*model.Hotel, 0, len(all))
	for _, h := range all {
		if filter.Match(h) {
			matched = append(matched, h)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Hotel{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateHotel(id string, input model.Hotel) (*model.Hotel, error) {
	existing, err := s.store.GetHotel(id)
	if err != nil {
		return nil, err
	}
	if input.Name != "" {
		existing.Name = input.Name
	}
	if input.City != "" {
		existing.City = input.City
	}
	if input.Address != "" {
		existing.Address = input.Address
	}
	if input.Star > 0 {
		existing.Star = input.Star
	}
	if input.Status != "" {
		existing.Status = input.Status
	}
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	existing.UpdatedAt = time.Now()
	if err := s.store.UpdateHotel(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) DeleteHotel(id string) error {
	for _, r := range s.store.ListRoomTypes() {
		if r.HotelID == id {
			return model.NewValidationError("hotel_id", "酒店下仍有房型，无法删除")
		}
	}
	return s.store.DeleteHotel(id)
}
