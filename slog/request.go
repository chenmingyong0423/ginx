package slog

import (
	"bytes"
	"io"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	XRequestIDKey = "X-Request-ID"
)

func Request(logger *slog.Logger) gin.HandlerFunc {
	return RequestWithConfig(logger, NewConfig())
}

func RequestWithConfig(logger *slog.Logger, cfg *Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		// 检查是否需要跳过
		for _, filter := range cfg.filters {
			if !filter(ctx) {
				return
			}
		}

		// 设置请求 ID
		rid := requestID(ctx, cfg)

		method := ctx.Request.Method
		path := ctx.Request.URL.Path
		route := ctx.FullPath()
		clientIp := ctx.ClientIP()
		referer := ctx.Request.Referer()
		host := ctx.Request.Host
		query := ctx.Request.URL.RawQuery
		params := make(map[string]any, len(ctx.Params))
		for _, param := range ctx.Params {
			params[param.Key] = param.Value
		}

		reqBody := newReqBodyReader(ctx.Request.Body, cfg)
		ctx.Request.Body = reqBody

		reqAttrs := []slog.Attr{
			slog.Any("time", start.UTC()),
			slog.String("method", method),
			slog.String("path", path),
			slog.String("query", query),
			slog.Any("params", params),
			slog.String("route", route),
			slog.String("client-ip", clientIp),
			slog.String("referer", referer),
			slog.String("host", host),
		}

		if cfg.logRequestID {
			reqAttrs = append(reqAttrs, slog.String("request-id", rid))
		}

		if cfg.logUserAgent {
			userAgent := ctx.Request.UserAgent()
			reqAttrs = append(reqAttrs, slog.String("user-agent", userAgent))
		}

		if cfg.logRequestHeader {
			headers := make([]any, 0, len(ctx.Request.Header))
			for k, v := range ctx.Request.Header {
				headers = append(headers, slog.Any(k, v))
			}
			reqAttrs = append(reqAttrs, slog.Group("header", headers...))
		}

		respBody := newRespBodyWriter(ctx.Writer, cfg)
		ctx.Writer = respBody

		ctx.Next()

		end := time.Now()

		reqAttrs = append(reqAttrs, slog.Int64("length", reqBody.bytes))
		if cfg.logRequestBody {
			reqAttrs = append(reqAttrs, slog.String("body", reqBody.body.String()))
		}

		level := cfg.level
		latency := end.Sub(start)
		statusCode := ctx.Writer.Status()

		respAttrs := []slog.Attr{
			slog.Any("time", end.UTC()),
			slog.Int("status", statusCode),
			slog.Int64("latency", int64(latency)),
		}

		if cfg.logResponseHeader {
			headers := make([]any, 0, len(ctx.Writer.Header()))
			for k, v := range ctx.Writer.Header() {
				headers = append(headers, slog.Any(k, v))
			}
			respAttrs = append(respAttrs, slog.Group("header", headers...))
		}

		respAttrs = append(respAttrs, slog.Int64("length", respBody.bytes))
		if cfg.logResponseBody {
			respAttrs = append(respAttrs, slog.String("body", respBody.body.String()))
		}

		attrs := append(
			[]slog.Attr{
				{
					Key:   "request",
					Value: slog.GroupValue(reqAttrs...),
				},
				{
					Key:   "response",
					Value: slog.GroupValue(respAttrs...),
				},
			},
		)
		logger.LogAttrs(ctx.Request.Context(), level, "HTTP", attrs...)
	}
}

func newRespBodyWriter(writer gin.ResponseWriter, cfg *Config) *respBodyWriter {
	w := &respBodyWriter{
		ResponseWriter: writer,
		body:           nil,
	}
	if cfg.logResponseBody {
		w.body = bytes.NewBufferString("")
	}
	return w
}

func newReqBodyReader(body io.ReadCloser, cfg *Config) *reqBodyReader {
	r := &reqBodyReader{
		ReadCloser: body,
		body:       nil,
	}
	if cfg.logRequestBody {
		r.body = bytes.NewBufferString("")
	}
	return r
}

func requestID(ctx *gin.Context, cfg *Config) (requestID string) {
	if cfg.logRequestID {
		requestID = ctx.GetHeader(XRequestIDKey)
		if requestID == "" {
			requestID = uuid.NewString()
			ctx.Request.Header.Set(XRequestIDKey, requestID)
			ctx.Set(XRequestIDKey, requestID)
		}
		ctx.Writer.Header().Set(XRequestIDKey, requestID)
	}
	return
}
