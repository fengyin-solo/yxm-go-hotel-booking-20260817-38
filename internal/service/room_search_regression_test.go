package service

import "testing"

func TestAvailableRoomSearchUsesActiveInventoryAndCheapestFirst(t *testing.T) {
	s := newTestService()
	hotel := mustCreateHotel(t, s, "会展中心酒店", "广州")
	cheap := mustCreateRoomType(t, s, hotel.ID, "舒适大床房", 32000, 3)
	expensive := mustCreateRoomType(t, s, hotel.ID, "行政套房", 92000, 2)
	otherCity := mustCreateHotel(t, s, "滨海酒店", "厦门")
	mustCreateRoomType(t, s, otherCity.ID, "海景房", 28000, 5)
	mustCreateBooking(t, s, cheap.ID, "周九", 2)

	result, err := s.SearchAvailableRoomTypes("广州", 1)
	if err != nil {
		t.Fatalf("搜索可订房型失败: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("广州应返回 2 个可订房型，实际 %d", len(result))
	}
	if result[0].RoomType.ID != cheap.ID || result[0].Available != 1 {
		t.Fatalf("低价房型应排第一且剩余 1 间，实际 id=%s available=%d", result[0].RoomType.ID, result[0].Available)
	}
	if result[1].RoomType.ID != expensive.ID || result[1].Available != 2 {
		t.Fatalf("高价房型应排第二且剩余 2 间，实际 id=%s available=%d", result[1].RoomType.ID, result[1].Available)
	}
}
