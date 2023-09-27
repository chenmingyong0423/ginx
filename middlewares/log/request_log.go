// Copyright 2023 chenmingyong0423

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"bytes"
	"io"
	"log/slog"
	"time"

	"github.com/chenmingyong0423/ginx/internal/types"
	"github.com/gin-gonic/gin"
)

type ConfigOption func(config *types.LoggerConfig)

func NewLoggerConfig(level slog.Level, opts ...ConfigOption) *types.LoggerConfig {
	LoggerConfig := &types.LoggerConfig{
		Level: level,
	}

	for _, opt := range opts {
		opt(LoggerConfig)
	}

	return LoggerConfig
}

func WithOptionalLogger(optionalLogger *slog.Logger) ConfigOption {
	return func(LoggerConfig *types.LoggerConfig) {
		LoggerConfig.OptionalLogger = optionalLogger
	}
}

// RequestLogger 用于打印每个请求的详细信息
// 默认使用 slog 库提供的默认实例进行打印，也可以传入一个 logger 实例
func RequestLogger(opt types.LoggerConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger := slog.Default()
		if opt.OptionalLogger != nil {
			logger = opt.OptionalLogger
		}
		rid := ctx.GetHeader(types.XRequestIDKey)
		start := time.Now()

		logRequest(ctx, logger, opt.Level, rid)

		mr := &types.Response{
			Body:           bytes.NewBufferString(""),
			ResponseWriter: ctx.Writer,
		}
		ctx.Writer = mr

		ctx.Next()

		logResponse(ctx, start, logger, opt.Level, mr.Body.String(), rid)
	}
}

func logRequest(ctx *gin.Context, logger *slog.Logger, level slog.Level, rid string) {
	method := ctx.Request.Method
	path := ctx.Request.URL.Path
	query := ctx.Request.URL.RawQuery
	body, _ := ctx.GetRawData()
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	logger.Log(ctx, level, "REQUEST", "Uri", path, "Query", query, "Method", method, "Body", string(body), "IP", ctx.ClientIP(), types.XRequestIDKey, rid)
}

func logResponse(ctx *gin.Context, start time.Time, logger *slog.Logger, level slog.Level, body string, rid string) {
	elapsedTime := time.Since(start)
	statusCode := ctx.Writer.Status()
	logger.Log(ctx, level, "RESPONSE", "Code", statusCode, "Body", body, "ElapseTime", elapsedTime, types.XRequestIDKey, rid)
}
