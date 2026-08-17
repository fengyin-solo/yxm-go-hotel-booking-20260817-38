# Bug Reproduction

## Bug 是什么

城市可订房型搜索会混淆活跃状态、剩余库存和价格排序，返回的可订结果不符合业务预期。

## 如何触发

在项目根目录运行：

```bash
go test ./internal/service -run TestAvailableRoomSearchUsesActiveInventoryAndCheapestFirst -count=20
```

## 错误信息

```text
低价房型应排第一且剩余 1 间，实际 available=2
```
