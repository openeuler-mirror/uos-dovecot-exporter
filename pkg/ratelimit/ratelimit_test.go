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

func TestRateLimiter_Stop(t *testing.T) {
        rl, err := NewRateLimiter(time.Second, 2)
        if err != nil {
                t.Fatalf("unexpected error: %v", err)
        }

        rl.Stop()

        if err := rl.Get(); err != ErrRateLimited {
                t.Fatalf("expected rate limit error after Stop, got: %v", err)
        }
}

func TestRateLimiter_TokenRefill(t *testing.T) {
        rl, err := NewRateLimiter(100*time.Millisecond, 2)
        if err != nil {
                t.Fatalf("unexpected error: %v", err)
        }
        defer rl.Stop()

        if err := rl.Get(); err != nil {
                t.Fatalf("expected token, got error: %v", err)
        }
        if err := rl.Get(); err != nil {
                t.Fatalf("expected token, got error: %v", err)
        }

        if err := rl.Get(); err != ErrRateLimited {
                t.Fatalf("expected rate limit error, got: %v", err)
        }

        time.Sleep(102 * time.Millisecond)
        if err := rl.Get(); err != nil {
                t.Fatalf("expected token after refill, got error: %v", err)
        }
}

