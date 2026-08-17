package service

import (
	"testing"

	"hotelbooking/internal/model"
)

func TestCreateRoomTypeRequiresHotel(t *testing.T) {
	s := newTestService()
	_, err := s.CreateRoomType(model.RoomType{
		HotelID:    "missing",
		Name:       "大床房",
		BedType:    "大床",
		Price:      39900,
		TotalRooms: 10,
	})
	if err == nil {
		t.Fatal("关联不存在的酒店应返回错误")
	}
}

func TestCreateRoomTypeValidation(t *testing.T) {
	s := newTestService()
	hotel := mustCreateHotel(t, s, "酒店", "北京")

	cases := []struct {
		name     string
		roomType model.RoomType
	}{
		{"负价格", model.RoomType{HotelID: hotel.ID, Name: "房", BedType: "大床", Price: -1, TotalRooms: 10}},
		{"零房间", model.RoomType{HotelID: hotel.ID, Name: "房", BedType: "大床", Price: 100, TotalRooms: 0}},
		{"空床型", model.RoomType{HotelID: hotel.ID, Name: "房", Price: 100, TotalRooms: 10}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, err := s.CreateRoomType(c.roomType); err == nil {
				t.Fatalf("应返回校验错误")
			}
		})
	}
}

func TestAvailableRooms(t *testing.T) {
	s := newTestService()
	hotel := mustCreateHotel(t, s, "酒店", "北京")
	roomType := mustCreateRoomType(t, s, hotel.ID, "大床房", 39900, 3)

	available, err := s.AvailableRooms(roomType.ID)
	if err != nil {
		t.Fatalf("AvailableRooms 失败: %v", err)
	}
	if available != 3 {
		t.Fatalf("初始可用应为 3，实际 %d", available)
	}

	mustCreateBooking(t, s, roomType.ID, "张三", 1)
	available, _ = s.AvailableRooms(roomType.ID)
	if available != 2 {
		t.Fatalf("预订 1 间后可用应为 2，实际 %d", available)
	}
}

func TestDeleteRoomTypeWithActiveBooking(t *testing.T) {
	s := newTestService()
	hotel := mustCreateHotel(t, s, "酒店", "北京")
	roomType := mustCreateRoomType(t, s, hotel.ID, "大床房", 39900, 5)
	mustCreateBooking(t, s, roomType.ID, "张三", 1)

	if err := s.DeleteRoomType(roomType.ID); err == nil {
		t.Fatal("存在进行中订单的房型应拒绝删除")
	}
}

func TestSearchAvailableRoomTypes(t *testing.T) {
	s := newTestService()
	hotelA := mustCreateHotel(t, s, "北京酒店", "北京")
	hotelB := mustCreateHotel(t, s, "上海酒店", "上海")
	rtA := mustCreateRoomType(t, s, hotelA.ID, "大床房", 39900, 2)
	rtB := mustCreateRoomType(t, s, hotelB.ID, "标间", 29900, 10)
	_ = rtB

	// 占满北京酒店的房型
	mustCreateBooking(t, s, rtA.ID, "张三", 2)

	result, err := s.SearchAvailableRoomTypes("北京", 1)
	if err != nil {
		t.Fatalf("搜索失败: %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("北京房型已被占满，应无结果，实际 %d", len(result))
	}

	result, err = s.SearchAvailableRoomTypes("上海", 1)
	if err != nil {
		t.Fatalf("搜索失败: %v", err)
	}
	if len(result) != 1 || result[0].Available != 10 {
		t.Fatalf("上海应有 1 个可用房型，实际 %d", len(result))
	}
}

func TestSearchAvailableRoomTypesRequiresCity(t *testing.T) {
	s := newTestService()
	if _, err := s.SearchAvailableRoomTypes("", 1); err == nil {
		t.Fatal("空城市应被拒绝")
	}
}
