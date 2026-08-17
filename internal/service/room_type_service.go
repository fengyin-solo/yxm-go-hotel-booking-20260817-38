package service

import (
	"sort"
	"time"

	"hotelbooking/internal/model"
	"hotelbooking/pkg/idgen"
)

func (s *Service) CreateRoomType(roomType model.RoomType) (*model.RoomType, error) {
	if err := roomType.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetHotel(roomType.HotelID); err != nil {
		return nil, model.NewValidationError("hotel_id", "关联的酒店不存在")
	}
	now := time.Now()
	roomType.ID = idgen.Hex()
	roomType.CreatedAt = now
	roomType.UpdatedAt = now
	if err := s.store.CreateRoomType(&roomType); err != nil {
		return nil, err
	}
	return &roomType, nil
}

func (s *Service) GetRoomType(id string) (*model.RoomType, error) {
	return s.store.GetRoomType(id)
}

func (s *Service) ListRoomTypes(filter model.RoomTypeFilter, page, size int) ([]*model.RoomType, int, error) {
	all := s.store.ListRoomTypes()
	matched := make([]*model.RoomType, 0, len(all))
	for _, r := range all {
		if filter.Match(r) {
			matched = append(matched, r)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.RoomType{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateRoomType(id string, input model.RoomType) (*model.RoomType, error) {
	existing, err := s.store.GetRoomType(id)
	if err != nil {
		return nil, err
	}
	if input.Name != "" {
		existing.Name = input.Name
	}
	if input.BedType != "" {
		existing.BedType = input.BedType
	}
	if input.Price > 0 {
		existing.Price = input.Price
	}
	if input.TotalRooms > 0 {
		existing.TotalRooms = input.TotalRooms
	}
	if input.Status != "" {
		existing.Status = input.Status
	}
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	existing.UpdatedAt = time.Now()
	if err := s.store.UpdateRoomType(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) DeleteRoomType(id string) error {
	for _, b := range s.store.ListBookings() {
		if b.RoomTypeID == id && b.Status != model.BookingCancelled && b.Status != model.BookingCompleted {
			return model.NewValidationError("room_type_id", "房型下仍有进行中的订单，无法删除")
		}
	}
	return s.store.DeleteRoomType(id)
}

// AvailableRooms 返回房型当前可预订的房间数。
func (s *Service) AvailableRooms(roomTypeID string) (int, error) {
	rt, err := s.store.GetRoomType(roomTypeID)
	if err != nil {
		return 0, err
	}
	occupied := s.occupiedRooms(roomTypeID, "")
	available := rt.TotalRooms - occupied
	if available < 0 {
		available = 0
	}
	return available, nil
}

// occupiedRooms 统计指定房型当前被占用的房间数（排除已取消/已完成订单，及 excludeID）。
func (s *Service) occupiedRooms(roomTypeID, excludeID string) int {
	occupied := 0
	for _, b := range s.store.ListBookings() {
		if b.RoomTypeID != roomTypeID {
			continue
		}
		if b.ID == excludeID {
			continue
		}
		if b.Status == model.BookingCancelled || b.Status == model.BookingCompleted {
			continue
		}
		occupied += b.RoomCount
	}
	return occupied
}

// RoomTypeAvailability 房型及其当前可预订房间数。
type RoomTypeAvailability struct {
	RoomType  *model.RoomType `json:"room_type"`
	Available int             `json:"available"`
}

// SearchAvailableRoomTypes 按城市搜索可满足指定房间数的房型。
func (s *Service) SearchAvailableRoomTypes(city string, roomCount int) ([]*RoomTypeAvailability, error) {
	if city == "" {
		return nil, model.NewValidationError("city", "城市不能为空")
	}
	if roomCount <= 0 {
		roomCount = 1
	}
	hotelIDs := make(map[string]bool)
	for _, h := range s.store.ListHotels() {
		if h.City == city && h.Status == model.HotelActive {
			hotelIDs[h.ID] = true
		}
	}
	result := make([]*RoomTypeAvailability, 0)
	for _, rt := range s.store.ListRoomTypes() {
		if !hotelIDs[rt.HotelID] {
			continue
		}
		if rt.Status == model.RoomTypeActive {
			continue
		}
		available := rt.TotalRooms - s.occupiedRooms(rt.ID, "")
		if available < roomCount {
			continue
		}
		result = append(result, &RoomTypeAvailability{RoomType: rt, Available: available})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].RoomType.Price > result[j].RoomType.Price
	})
	return result, nil
}
