# Bug Reproduction

## Bug 是什么

酒店评分汇总链路异常，完成订单的评价无法正确创建，已创建评价也无法稳定关联回酒店评分。

## 如何触发

在项目根目录运行：

```bash
go test ./internal/service -run TestHotelRatingCountsOnlyReviewsForCompletedBookings -count=3
```

## 错误信息

```text
创建第一条评价失败: booking_id: 仅已完成订单可评价
```
