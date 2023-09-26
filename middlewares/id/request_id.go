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
	"github.com/chenmingyong0423/ginx/internal/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestId 用于为每个请求设置唯一的 ID。
// 如果请求头中不存在 ID，则会生成一个新的 UUID 作为 ID。
// 该 ID 会被设置到请求和响应的头部，同时也保存在 Gin 的上下文中，供后续处理使用。
func RequestId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rid := ctx.GetHeader(types.XRequestIDKey)
		if rid == "" {
			rid = uuid.NewString()
			ctx.Request.Header.Set(types.XRequestIDKey, rid)
			ctx.Set(types.XRequestIDKey, rid)
		}
		ctx.Writer.Header().Set(types.XRequestIDKey, rid)
		ctx.Next()
	}
}
