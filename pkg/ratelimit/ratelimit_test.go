// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT
package ratelimit

import (
        "fmt"
        "testing"
        "time"
)

func TestNewRateLimiter(t *testing.T) {
        tests := []struct {
                name     string
                limit    time.Duration
                chanSize int
                expected error
        }{
                {"Normal", time.Second, 5, nil},                       // Valid case
                {"InvalidLimit", -time.Second, 5, ErrRateLimitTime},   // Invalid limit
                {"InvalidChanSize", time.Second, 0, ErrRateLimitSize}, // Invalid channel size
        }

        for _, tt := range tests {
                t.Run(tt.name, func(t *testing.T) {
                        rl, err := NewRateLimiter(tt.limit, tt.chanSize)
                        if err != nil && err != tt.expected {
                                t.Fatalf("expected error %v, got %v", tt.expected, err)
                        }
                        if err == nil && rl == nil {
                                t.Fatal("expected non-nil rate limiter")
                        }
                })
        }
}

func TestRateLimiter_Get(t *testing.T) {
        rl, err := NewRateLimiter(time.Second, 2)
        if err != nil {
                t.Fatalf("unexpected error: %v", err)
        }

        if err := rl.Get(); err != nil {
                t.Fatalf("expected token, got error: %v", err)
        }

        if err := rl.Get(); err != nil {
                t.Fatalf("expected token, got error: %v", err)
        }

        if err := rl.Get(); err != ErrRateLimited {
                t.Fatalf("expected rate limit error, got: %v", err)
        }
}


