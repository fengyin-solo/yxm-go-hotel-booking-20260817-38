# Bug Reproduction

## Bug 是什么

开票生命周期被破坏：未完成订单可以提前开票，完成后的开票、重复开票拦截和订单关联保存不一致。

## 如何触发

在项目根目录运行：

```bash
go test ./internal/service -run TestInvoiceOnlyAfterCompletedBookingAndKeepsBookingLink -count=20
```

## 错误信息

```text
未完成订单不应允许开票
```
