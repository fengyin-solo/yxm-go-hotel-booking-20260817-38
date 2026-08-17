# Bug Reproduction

## Bug 是什么

待确认订单取消失败，取消流程没有把订单稳定流转为 `cancelled`，房型库存也不会释放。

## 如何触发

在项目根目录运行：

```bash
go test ./internal/service -run TestCancelPendingBookingReleasesRoomInventory -count=20
```

## 错误信息

```text
pending 订单应允许取消: status: 当前订单状态无法取消
```
