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

package slog

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/require"
)

func TestRequest(t *testing.T) {
	type response struct {
		Code int
		Body *bytes.Buffer
	}
	tests := []struct {
		name    string
		logger  *slog.Logger
		cfg     *Config
		reqFunc func() *http.Request

		want response
	}{
		{
			name:   "默认配置",
			logger: slog.Default(),
			cfg:    NewConfig(),
			reqFunc: func() *http.Request {
				req, err := http.NewRequest("GET", "/slog/1?name=陈明勇&age=18", nil)
				require.NoError(t, err)
				// 设置 ip
				req.RemoteAddr = "127.0.0.1:443"
				// 设置 host
				req.Host = "localhost"
				// 设置 referer
				req.Header.Set("Referer", "http://localhost:8080")
				return req
			},
			want: response{
				Code: 200,
				Body: bytes.NewBufferString(`{"message":"slog"}`),
			},
		},
		{
			name: "修改日志级别",
			logger: func() *slog.Logger {
				slog.SetLogLoggerLevel(slog.LevelDebug)
				return slog.Default()
			}(),
			cfg: NewConfig(WithConfigLevel(slog.LevelDebug)),
			reqFunc: func() *http.Request {
				req, err := http.NewRequest("GET", "/slog/1?name=陈明勇&age=18", nil)
				require.NoError(t, err)
				// 设置 ip
				req.RemoteAddr = "127.0.0.1:443"
				// 设置 host
				req.Host = "localhost"
				// 设置 referer
				req.Header.Set("Referer", "http://localhost:8080")
				return req
			},
			want: response{
				Code: 200,
				Body: bytes.NewBufferString(`{"message":"slog"}`),
			},
		},
		{
			name:   "不打印请求 ID",
			logger: slog.Default(),
			cfg:    NewConfig(WithConfigLogRequestID(false)),
			reqFunc: func() *http.Request {
				req, err := http.NewRequest("GET", "/slog/1?name=陈明勇&age=18", nil)
				require.NoError(t, err)
				// 设置 ip
				req.RemoteAddr = "127.0.0.1:443"
				// 设置 host
				req.Host = "localhost"
				// 设置 referer
				req.Header.Set("Referer", "http://localhost:8080")
				return req
			},
			want: response{
				Code: 200,
				Body: bytes.NewBufferString(`{"message":"slog"}`),
			},
		},
		{
			name:   "打印 request header",
			logger: slog.Default(),
			cfg:    NewConfig(WithConfigLogRequestHeader(true)),
			reqFunc: func() *http.Request {
				req, err := http.NewRequest("GET", "/slog/1?name=陈明勇&age=18", nil)
				require.NoError(t, err)
				// 设置 ip
				req.RemoteAddr = "127.0.0.1:443"
				// 设置 host
				req.Host = "localhost"
				// 设置 referer
				req.Header.Set("Referer", "http://localhost:8080")
				req.Header.Set(XRequestIDKey, "666666")
				return req
			},
			want: response{
				Code: 200,
				Body: bytes.NewBufferString(`{"message":"slog"}`),
			},
		},
		{
			name:   "打印 request body",
			logger: slog.Default(),
			cfg:    NewConfig(WithConfigLogRequestBody(true)),
			reqFunc: func() *http.Request {
				req, err := http.NewRequest("POST", "/slog", strings.NewReader(`{"name":"陈明勇","age":18}`))
				require.NoError(t, err)
				// 设置 ip
				req.RemoteAddr = "127.0.0.1:443"
				// 设置 host
				req.Host = "localhost"
				// 设置 referer
				req.Header.Set("Referer", "http://localhost:8080")
				return req
			},
			want: response{
				Code: 200,
				Body: bytes.NewBufferString(`{"message":"slog"}`),
			},
		},
		{
			name:   "打印 user agent",
			logger: slog.Default(),
			cfg:    NewConfig(WithConfigLogUserAgent(true)),
			reqFunc: func() *http.Request {
				req, err := http.NewRequest("GET", "/slog/1?name=陈明勇&age=18", nil)
				require.NoError(t, err)
				// 设置 ip
				req.RemoteAddr = "127.0.0.1:443"
				// 设置 host
				req.Host = "localhost"
				// 设置 referer
				req.Header.Set("Referer", "http://localhost:8080")
				// 设置 user agent
				req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")
				req.Header.Set(XRequestIDKey, "666666")
				return req
			},
			want: response{
				Code: 200,
				Body: bytes.NewBufferString(`{"message":"slog"}`),
			},
		},
		{
			name:   "打印 response header",
			logger: slog.Default(),
			cfg:    NewConfig(WithConfigLogResponseHeader(true)),
			reqFunc: func() *http.Request {
				req, err := http.NewRequest("GET", "/slog/1?name=陈明勇&age=18", nil)
				require.NoError(t, err)
				// 设置 ip
				req.RemoteAddr = "127.0.0.1:443"
				// 设置 host
				req.Host = "localhost"
				// 设置 referer
				req.Header.Set("Referer", "http://localhost:8080")
				req.Header.Set(XRequestIDKey, "666666")
				return req
			},
			want: response{
				Code: 200,
				Body: bytes.NewBufferString(`{"message":"slog"}`),
			},
		},
		{
			name:   "打印 response body",
			logger: slog.Default(),
			cfg:    NewConfig(WithConfigLogResponseBody(true)),
			reqFunc: func() *http.Request {
				req, err := http.NewRequest("GET", "/slog/1?name=陈明勇&age=18", nil)
				require.NoError(t, err)
				// 设置 ip
				req.RemoteAddr = "127.0.0.1:443"
				// 设置 host
				req.Host = "localhost"
				// 设置 referer
				req.Header.Set("Referer", "http://localhost:8080")
				req.Header.Set(XRequestIDKey, "666666")
				return req
			},
			want: response{
				Code: 200,
				Body: bytes.NewBufferString(`{"message":"slog"}`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := gin.Default()
			server.Use(RequestWithConfig(tt.logger, tt.cfg))
			server.GET("/slog/:id", func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{
					"message": "slog",
				})
			})

			server.POST("/slog", func(ctx *gin.Context) {
				mp := make(map[string]interface{})
				err := ctx.BindJSON(&mp)
				require.NoError(t, err)
				time.Sleep(time.Millisecond * 100)
				ctx.JSON(http.StatusOK, gin.H{
					"message": "slog",
				})
			})

			resp := httptest.NewRecorder()

			server.ServeHTTP(resp, tt.reqFunc())

			require.Equal(t, resp.Code, tt.want.Code)
			require.Equal(t, resp.Body, tt.want.Body)
		})
	}
}
