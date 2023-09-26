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

package id

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRequestId(t *testing.T) {
	testCases := []struct {
		name           string
		requestBuilder func() *http.Request

		validFunc func(value string) bool
	}{
		{
			name: "Header 里没有 X-Request-ID 参数",
			requestBuilder: func() *http.Request {
				req, _ := http.NewRequest("GET", "/test", nil)
				return req
			},
			validFunc: func(value string) bool {
				return value != ""
			},
		},
		{
			name: "Header 里有 X-Request-ID 参数",
			requestBuilder: func() *http.Request {
				req, _ := http.NewRequest("GET", "/test", nil)
				req.Header.Set(xRequestIDKey, "chenmingyong")
				return req
			},
			validFunc: func(value string) bool {
				return value == "chenmingyong"
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 创建 gin engine
			engine := gin.Default()
			engine.Use(RequestId())
			engine.GET("/test", func(ctx *gin.Context) {
				ctx.String(200, "test")
			})

			// 创建 request 请求
			req := tc.requestBuilder()
			w := httptest.NewRecorder()
			// 接口调用
			engine.ServeHTTP(w, req)
			assert.True(t, tc.validFunc(w.Header().Get(xRequestIDKey)))
		})
	}
}
