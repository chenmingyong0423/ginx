package slog

import "github.com/gin-gonic/gin"

// FilterFunc 用于定义日志过滤器函数。
// 如果返回 true，则表示允许记录日志；否则，表示不允许记录日志。
type FilterFunc func(ctx *gin.Context) bool
