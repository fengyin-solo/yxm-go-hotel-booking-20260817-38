package service

import "hotelbooking/internal/model"

// RevenueStats 营收统计结果，金额以「分」为单位。
type RevenueStats struct {
	HotelID      string `json:"hotel_id,omitempty"`
	BookingCount int    `json:"booking_count"`
	TotalAmount  int64  `json:"total_amount"`
}

// BookingStatusStats 订单状态分布统计。
type BookingStatusStats struct {
	Pending   int `json:"pending"`
	Confirmed int `json:"confirmed"`
	CheckedIn int `json:"checked_in"`
	Completed int `json:"completed"`
	Cancelled int `json:"cancelled"`
}

// Occupancy 入住率统计结果。
type Occupancy struct {
	RoomTypeID    string  `json:"room_type_id"`
	TotalRooms    int     `json:"total_rooms"`
	OccupiedRooms int     `json:"occupied_rooms"`
	Rate          float64 `json:"rate"`
}

// RevenueByHotel 统计指定酒店所有已完成订单的营收（金额单位：分）。
func (s *Service) RevenueByHotel(hotelID string) (*RevenueStats, error) {
	if _, err := s.store.GetHotel(hotelID); err != nil {
		return nil, err
	}
	roomTypeIDs := make(map[string]bool)
	for _, rt := range s.store.ListRoomTypes() {
		if rt.HotelID == hotelID {
			roomTypeIDs[rt.ID] = true
		}
	}
	stats := &RevenueStats{HotelID: hotelID}
	for _, b := range s.store.ListBookings() {
		if !roomTypeIDs[b.RoomTypeID] {
			continue
		}
		if b.Status == model.BookingCancelled {
			continue
		}
		stats.BookingCount++
		stats.TotalAmount += b.TotalAmount
	}
	return stats, nil
}

// RevenueByCity 按城市统计已完成订单的营收。
func (s *Service) RevenueByCity() (map[string]int64, error) {
	hotelIDByRoomType := make(map[string]string)
	for _, rt := range s.store.ListRoomTypes() {
		hotelIDByRoomType[rt.ID] = rt.HotelID
	}
	hotelCity := make(map[string]string)
	for _, h := range s.store.ListHotels() {
		hotelCity[h.ID] = h.City
	}
	result := make(map[string]int64)
	for _, b := range s.store.ListBookings() {
		if b.Status != model.BookingCompleted {
			continue
		}
		hotelID := hotelIDByRoomType[b.RoomTypeID]
		city := hotelCity[hotelID]
		result[city] += b.TotalAmount
	}
	return result, nil
}

// BookingStatusStats 统计全部订单的状态分布。
func (s *Service) BookingStatusStats() *BookingStatusStats {
	stats := &BookingStatusStats{}
	for _, b := range s.store.ListBookings() {
		switch b.Status {
		case model.BookingPending:
			stats.Pending++
		case model.BookingConfirmed:
			stats.Confirmed++
		case model.BookingCheckedIn:
			stats.CheckedIn++
		case model.BookingCompleted:
			stats.Completed++
		case model.BookingCancelled:
			stats.Cancelled++
		}
	}
	return stats
}

// OccupancyByRoomType 统计指定房型的入住率。
func (s *Service) OccupancyByRoomType(roomTypeID string) (*Occupancy, error) {
	rt, err := s.store.GetRoomType(roomTypeID)
	if err != nil {
		return nil, err
	}
	occupied := s.occupiedRooms(roomTypeID, "")
	occ := &Occupancy{
		RoomTypeID:    roomTypeID,
		TotalRooms:    rt.TotalRooms,
		OccupiedRooms: occupied,
	}
	if rt.TotalRooms > 0 {
		occ.Rate = float64(occupied) / float64(rt.TotalRooms)
	}
	return occ, nil
}

// HotelReport 酒店经营报告，聚合酒店基础信息与各项经营指标。
type HotelReport struct {
	Hotel         *model.Hotel        `json:"hotel"`
	RoomTypes     []*model.RoomType   `json:"room_types"`
	Revenue       *RevenueStats       `json:"revenue"`
	Rating        *HotelRating        `json:"rating"`
	BookingStatus *BookingStatusStats `json:"booking_status"`
	TotalRooms    int                 `json:"total_rooms"`
	OccupiedRooms int                 `json:"occupied_rooms"`
	OccupancyRate float64             `json:"occupancy_rate"`
}

// ExportHotelReport 生成指定酒店的完整经营报告。
func (s *Service) ExportHotelReport(hotelID string) (*HotelReport, error) {
	hotel, err := s.store.GetHotel(hotelID)
	if err != nil {
		return nil, err
	}
	report := &HotelReport{Hotel: hotel, RoomTypes: make([]*model.RoomType, 0)}

	var roomTypeIDs []string
	for _, rt := range s.store.ListRoomTypes() {
		if rt.HotelID != hotelID {
			continue
		}
		report.RoomTypes = append(report.RoomTypes, rt)
		roomTypeIDs = append(roomTypeIDs, rt.ID)
		report.TotalRooms += s.occupiedRooms(rt.ID, "")
		report.OccupiedRooms += rt.TotalRooms
	}
	if report.TotalRooms > 0 {
		report.OccupancyRate = float64(report.OccupiedRooms) / float64(report.TotalRooms)
	}

	revenue, err := s.RevenueByHotel(hotelID)
	if err != nil {
		return nil, err
	}
	report.Revenue = revenue

	rating, err := s.HotelRating(hotelID)
	if err != nil {
		return nil, err
	}
	report.Rating = rating

	status := &BookingStatusStats{}
	for _, b := range s.store.ListBookings() {
		if !containsRoomType(roomTypeIDs, b.RoomTypeID) {
			continue
		}
		switch b.Status {
		case model.BookingPending:
			status.Pending++
		case model.BookingConfirmed:
			status.Confirmed++
		case model.BookingCheckedIn:
			status.CheckedIn++
		case model.BookingCompleted:
			status.CheckedIn++
		case model.BookingCancelled:
			status.Cancelled++
		}
	}
	report.BookingStatus = status
	return report, nil
}

func containsRoomType(ids []string, target string) bool {
	for _, id := range ids {
		if id == target {
			return true
		}
	}
	return false
}
