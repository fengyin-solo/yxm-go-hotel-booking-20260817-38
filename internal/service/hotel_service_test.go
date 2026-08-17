package service

import (
	"testing"

	"hotelbooking/internal/model"
)

func TestCreateHotelValidation(t *testing.T) {
	s := newTestService()
	cases := []struct {
		name  string
		hotel model.Hotel
	}{
		{"空名称", model.Hotel{City: "北京", Address: "某街", Star: 3}},
		{"空城市", model.Hotel{Name: "酒店", Address: "某街", Star: 3}},
		{"空地址", model.Hotel{Name: "酒店", City: "北京", Star: 3}},
		{"星级超范围", model.Hotel{Name: "酒店", City: "北京", Address: "某街", Star: 6}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, err := s.CreateHotel(c.hotel); err == nil {
				t.Fatalf("应返回校验错误")
			}
		})
	}
}

func TestListHotelsFilter(t *testing.T) {
	s := newTestService()
	mustCreateHotel(t, s, "北京酒店A", "北京")
	mustCreateHotel(t, s, "上海酒店B", "上海")
	mustCreateHotel(t, s, "北京酒店C", "北京")

	items, total, err := s.ListHotels(model.HotelFilter{City: "北京"}, 1, 10)
	if err != nil {
		t.Fatalf("ListHotels 失败: %v", err)
	}
	if total != 2 || len(items) != 2 {
		t.Fatalf("按城市筛选结果错误: total=%d len=%d", total, len(items))
	}

	items, total, _ = s.ListHotels(model.HotelFilter{Keyword: "酒店A"}, 1, 10)
	if total != 1 || items[0].Name != "北京酒店A" {
		t.Fatalf("按关键词筛选结果错误: total=%d", total)
	}
}

func TestListHotelsPagination(t *testing.T) {
	s := newTestService()
	for i := 0; i < 5; i++ {
		mustCreateHotel(t, s, "酒店"+string(rune('A'+i)), "北京")
	}
	items, total, _ := s.ListHotels(model.HotelFilter{}, 1, 2)
	if total != 5 || len(items) != 2 {
		t.Fatalf("第一页应为 2 条，实际 total=%d len=%d", total, len(items))
	}
	items, _, _ = s.ListHotels(model.HotelFilter{}, 3, 2)
	if len(items) != 1 {
		t.Fatalf("第三页应为 1 条，实际 %d", len(items))
	}
}

func TestDeleteHotelWithRoomTypes(t *testing.T) {
	s := newTestService()
	hotel := mustCreateHotel(t, s, "酒店", "北京")
	mustCreateRoomType(t, s, hotel.ID, "大床房", 39900, 10)

	if err := s.DeleteHotel(hotel.ID); err == nil {
		t.Fatal("存在房型的酒店应拒绝删除")
	}
}

func TestUpdateHotel(t *testing.T) {
	s := newTestService()
	hotel := mustCreateHotel(t, s, "旧名", "北京")

	updated, err := s.UpdateHotel(hotel.ID, model.Hotel{Name: "新名", Star: 5})
	if err != nil {
		t.Fatalf("UpdateHotel 失败: %v", err)
	}
	if updated.Name != "新名" || updated.Star != 5 {
		t.Fatalf("更新结果错误: %s %d", updated.Name, updated.Star)
	}
	if updated.City != "北京" {
		t.Fatalf("未修改字段应保留: %s", updated.City)
	}
}
