// Generated by [optioner] command-line tool; DO NOT EDIT
// If you have any questions, please create issues and submit contributions at:
// https://github.com/chenmingyong0423/go-optioner

package slog

import (
	"log/slog"
)

//go:generate optioner -type Config -with_prefix Config -output ./config.go -mode append
type Config struct {
	level slog.Level

	logRequestID     bool
	logUserAgent     bool
	logRequestHeader bool
	logRequestBody   bool

	logResponseHeader bool
	logResponseBody   bool

	filters []FilterFunc
}

type ConfigOption func(*Config)

func NewConfig(opts ...ConfigOption) *Config {
	config := &Config{
		level: slog.LevelInfo,

		logRequestID: true,

		filters: []FilterFunc{},
	}

	for _, opt := range opts {
		opt(config)
	}

	return config
}

func WithConfigLevel(level slog.Level) ConfigOption {
	return func(config *Config) {
		config.level = level
	}
}

func WithConfigLogRequestID(logRequestID bool) ConfigOption {
	return func(config *Config) {
		config.logRequestID = logRequestID
	}
}

func WithConfigLogUserAgent(logUserAgent bool) ConfigOption {
	return func(config *Config) {
		config.logUserAgent = logUserAgent
	}
}

func WithConfigLogRequestHeader(logRequestHeader bool) ConfigOption {
	return func(config *Config) {
		config.logRequestHeader = logRequestHeader
	}
}

func WithConfigLogRequestBody(logRequestBody bool) ConfigOption {
	return func(config *Config) {
		config.logRequestBody = logRequestBody
	}
}

func WithConfigLogResponseHeader(logResponseHeader bool) ConfigOption {
	return func(config *Config) {
		config.logResponseHeader = logResponseHeader
	}
}

func WithConfigLogResponseBody(logResponseBody bool) ConfigOption {
	return func(config *Config) {
		config.logResponseBody = logResponseBody
	}
}

func WithConfigFilters(filters []FilterFunc) ConfigOption {
	return func(config *Config) {
		config.filters = filters
	}
}