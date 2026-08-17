# Bug Reproduction

## Bug 是什么

酒店经营报表的容量、占用、状态分布和营收统计口径不一致，导致同一份报表内部数值冲突。

## 如何触发

在项目根目录运行：

```bash
go test ./internal/service -run TestHotelReportSeparatesCapacityOccupancyAndCompletedRevenue -count=3
```

## 错误信息

```text
房间总数/占用数错误: total=2 occupied=6
```
