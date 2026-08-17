package service

import (
	"sort"
	"time"

	"hotelbooking/internal/model"
	"hotelbooking/pkg/idgen"
)

// BatchResult 批量操作结果。
type BatchResult struct {
	Succeeded []*model.Booking  `json:"succeeded"`
	Failed    map[string]string `json:"failed"`
}

func (s *Service) CreateBooking(booking model.Booking) (*model.Booking, error) {
	if err := booking.Validate(); err != nil {
		return nil, err
	}
	rt, err := s.store.GetRoomType(booking.RoomTypeID)
	if err != nil {
		return nil, model.NewValidationError("room_type_id", "关联的房型不存在")
	}
	if occupied := s.occupiedRooms(booking.RoomTypeID, ""); occupied+booking.RoomCount > rt.TotalRooms {
		return nil, model.NewValidationError("room_count", "房型库存不足")
	}
	now := time.Now()
	booking.ID = idgen.Hex()
	booking.Status = model.BookingPending
	booking.TotalAmount = rt.Price * int64(booking.Nights()) * int64(booking.RoomCount)
	booking.CreatedAt = now
	booking.UpdatedAt = now
	if err := s.store.CreateBooking(&booking); err != nil {
		return nil, err
	}
	return &booking, nil
}

func (s *Service) GetBooking(id string) (*model.Booking, error) {
	return s.store.GetBooking(id)
}

func (s *Service) ListBookings(filter model.BookingFilter, page, size int) ([]*model.Booking, int, error) {
	all := s.store.ListBookings()
	matched := make([]*model.Booking, 0, len(all))
	for _, b := range all {
		if filter.Match(b) {
			matched = append(matched, b)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Booking{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// UpdateBooking 更新订单客人信息，仅待确认状态可修改。
func (s *Service) UpdateBooking(id string, input model.Booking) (*model.Booking, error) {
	existing, err := s.store.GetBooking(id)
	if err != nil {
		return nil, err
	}
	if existing.Status != model.BookingPending {
		return nil, model.NewValidationError("status", "仅待确认订单可修改")
	}
	if input.GuestName != "" {
		existing.GuestName = input.GuestName
	}
	if input.GuestPhone != "" {
		existing.GuestPhone = input.GuestPhone
	}
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	existing.UpdatedAt = time.Now()
	if err := s.store.UpdateBooking(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// ConfirmBooking 确认订单：pending -> confirmed。
func (s *Service) ConfirmBooking(id string) (*model.Booking, error) {
	b, err := s.store.GetBooking(id)
	if err != nil {
		return nil, err
	}
	if !model.CanTransitionBooking(b.Status, model.BookingConfirmed) {
		return nil, model.NewValidationError("status", "当前订单状态无法确认")
	}
	b.Status = model.BookingConfirmed
	b.UpdatedAt = time.Now()
	if err := s.store.UpdateBooking(b); err != nil {
		return nil, err
	}
	return b, nil
}

// CancelBooking 取消订单：pending/confirmed -> cancelled。
func (s *Service) CancelBooking(id string) (*model.Booking, error) {
	b, err := s.store.GetBooking(id)
	if err != nil {
		return nil, err
	}
	if !model.CanTransitionBooking(b.Status, model.BookingCancelled) {
		return nil, model.NewValidationError("status", "当前订单状态无法取消")
	}
	b.Status = model.BookingCancelled
	b.UpdatedAt = time.Now()
	if err := s.store.UpdateBooking(b); err != nil {
		return nil, err
	}
	return b, nil
}

// CheckInBooking 办理入住：confirmed -> checked_in，并生成入住记录。
func (s *Service) CheckInBooking(id, roomNumber string) (*model.Booking, error) {
	b, err := s.store.GetBooking(id)
	if err != nil {
		return nil, err
	}
	if !model.CanTransitionBooking(b.Status, model.BookingCheckedIn) {
		return nil, model.NewValidationError("status", "仅已确认订单可办理入住")
	}
	now := time.Now()
	checkIn := &model.CheckIn{
		ID:          idgen.Hex(),
		BookingID:   b.ID,
		RoomNumber:  roomNumber,
		CheckedInAt: now,
		Status:      model.CheckInActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := checkIn.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.CreateCheckIn(checkIn); err != nil {
		return nil, err
	}
	b.Status = model.BookingCheckedIn
	b.UpdatedAt = now
	if err := s.store.UpdateBooking(b); err != nil {
		return nil, err
	}
	return b, nil
}

// CheckOutBooking 办理退房：checked_in -> completed，并结束入住记录。
func (s *Service) CheckOutBooking(id string) (*model.Booking, error) {
	b, err := s.store.GetBooking(id)
	if err != nil {
		return nil, err
	}
	if !model.CanTransitionBooking(b.Status, model.BookingCompleted) {
		return nil, model.NewValidationError("status", "仅入住中的订单可退房")
	}
	checkIn, err := s.store.GetCheckInByBooking(id)
	if err != nil {
		return nil, model.NewValidationError("checkin", "未找到对应的入住记录")
	}
	now := time.Now()
	checkIn.CheckedOutAt = &now
	checkIn.Status = model.CheckInFinished
	checkIn.UpdatedAt = now
	if err := s.store.UpdateCheckIn(checkIn); err != nil {
		return nil, err
	}
	b.Status = model.BookingCompleted
	b.UpdatedAt = now
	if err := s.store.UpdateBooking(b); err != nil {
		return nil, err
	}
	return b, nil
}

// BatchConfirm 批量确认订单，逐条尝试并汇总成功与失败结果。
func (s *Service) BatchConfirm(ids []string) *BatchResult {
	result := &BatchResult{
		Succeeded: make([]*model.Booking, 0),
		Failed:    make(map[string]string),
	}
	for _, id := range ids {
		b, err := s.ConfirmBooking(id)
		if err != nil {
			result.Failed[id] = err.Error()
			continue
		}
		result.Succeeded = append(result.Succeeded, b)
	}
	return result
}
