package service

import (
	"sort"

	"hotelbooking/internal/model"
)

func (s *Service) GetCheckIn(id string) (*model.CheckIn, error) {
	return s.store.GetCheckIn(id)
}

func (s *Service) GetCheckInByBooking(bookingID string) (*model.CheckIn, error) {
	return s.store.GetCheckInByBooking(bookingID)
}

// ListCheckIns 按订单与状态过滤入住记录并分页返回。
func (s *Service) ListCheckIns(bookingID, status string, page, size int) ([]*model.CheckIn, int, error) {
	all := s.store.ListCheckIns()
	matched := make([]*model.CheckIn, 0, len(all))
	for _, c := range all {
		if bookingID != "" && c.BookingID != bookingID {
			continue
		}
		if status != "" && c.Status != status {
			continue
		}
		matched = append(matched, c)
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CheckedInAt.After(matched[j].CheckedInAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.CheckIn{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}
