# Hotel Booking 酒店预订服务

一个纯 Go 标准库实现的酒店预订后端，涵盖酒店、房型、订单、入住、评价与发票全流程，含订单状态机、库存控制、营收与入住率统计、批量操作与经营报告导出。

> 说明：所有金额字段均以「分」为单位存储与计算，避免浮点精度问题。

## 运行

```bash
go run ./cmd/server
```

环境变量：

| 变量 | 默认值 | 说明 |
|------|-------|------|
| `PORT` | 8080 | 监听端口 |
| `ADDR` | `:8080` | 完整监听地址（优先于 PORT） |
| `MAX_PAGE_SIZE` | 100 | 单页最大条数 |
| `LOW_STOCK_THRESHOLD` | 5 | 低库存阈值 |
| `LOG_LEVEL` | info | 日志级别：debug/info/warn/error |

## API 一览

统一响应结构：`{"code":0,"message":"ok","data":...}`；错误时 `code` 非 0 且 `message` 说明原因。

### 酒店 Hotel

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/hotels` | 创建酒店 |
| GET | `/api/hotels` | 分页列出（city/star/status/keyword 筛选） |
| GET | `/api/hotels/{id}` | 获取酒店详情 |
| PUT | `/api/hotels/{id}` | 更新酒店 |
| DELETE | `/api/hotels/{id}` | 删除酒店（有房型时拒绝） |

### 房型 RoomType

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/room-types` | 创建房型 |
| GET | `/api/room-types` | 分页列出（hotel_id/status 筛选） |
| GET | `/api/room-types/{id}` | 获取房型详情 |
| PUT | `/api/room-types/{id}` | 更新房型 |
| DELETE | `/api/room-types/{id}` | 删除房型（有进行中订单时拒绝） |
| GET | `/api/room-types/{id}/availability` | 查询可预订房间数 |
| GET | `/api/room-types/search?city=&rooms=` | 按城市搜索可用房型 |

### 订单 Booking

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/bookings` | 创建订单（自动计算金额并校验库存） |
| GET | `/api/bookings` | 分页列出（status/guest_name/from/to 筛选） |
| GET | `/api/bookings/{id}` | 获取订单详情 |
| PUT | `/api/bookings/{id}` | 修改客人信息（仅待确认） |
| POST | `/api/bookings/{id}/confirm` | 确认订单 |
| POST | `/api/bookings/{id}/cancel` | 取消订单 |
| POST | `/api/bookings/{id}/check-in` | 办理入住 |
| POST | `/api/bookings/{id}/check-out` | 办理退房 |
| POST | `/api/bookings/batch-confirm` | 批量确认订单 |

创建订单请求体（日期格式 `YYYY-MM-DD`）：

```json
{ "room_type_id": "<房型ID>", "guest_name": "张三", "guest_phone": "13800000000", "check_in": "2026-08-20", "check_out": "2026-08-22", "room_count": 1 }
```

订单状态机流转：`pending → confirmed → checked_in → completed`，`pending/confirmed → cancelled`。

### 入住 CheckIn

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/checkins` | 分页列出（booking_id/status 筛选） |
| GET | `/api/checkins/{id}` | 获取入住记录 |
| GET | `/api/bookings/{id}/checkin` | 按订单查询入住记录 |

### 评价 Review

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/reviews` | 创建评价（仅已完成订单，且一单一评） |
| GET | `/api/reviews` | 分页列出（booking_id/min_rating 筛选） |
| GET | `/api/reviews/{id}` | 获取评价 |
| PUT | `/api/reviews/{id}` | 更新评价 |
| DELETE | `/api/reviews/{id}` | 删除评价 |

### 发票 Invoice

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/invoices` | 开票（仅已完成订单，金额自动取订单金额） |
| GET | `/api/invoices` | 分页列出（booking_id 筛选） |
| GET | `/api/invoices/{id}` | 获取发票 |
| DELETE | `/api/invoices/{id}` | 删除发票 |

### 统计 Stats

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/stats/booking-status` | 订单状态分布 |
| GET | `/api/stats/revenue/city` | 按城市统计营收 |
| GET | `/api/hotels/{id}/revenue` | 酒店营收统计 |
| GET | `/api/hotels/{id}/rating` | 酒店评分聚合 |
| GET | `/api/hotels/{id}/report` | 酒店经营报告 |
| GET | `/api/room-types/{id}/occupancy` | 房型入住率 |
