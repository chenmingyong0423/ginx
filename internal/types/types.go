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

package types

import (
	"bytes"
	"log/slog"

	"github.com/gin-gonic/gin"
)

const (
	XRequestIDKey = "X-Request-ID"
)

type Response struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (w *Response) Write(b []byte) (int, error) {
	w.Body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *Response) WriteString(s string) (int, error) {
	w.Body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

//go:generate optioner -type LoggerConfig -output ../../middlewares/log/request_log.go -mode append
type LoggerConfig struct {
	Level          slog.Level `opt:"-"`
	OptionalLogger *slog.Logger
	SkipPaths      []string
	SkipFunc       func(ctx *gin.Context) bool
}
