// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT
package server

import (
        "sync"
        "net/http"
        "github.com/sirupsen/logrus"
        "uos-dovecot-exporter/pkg/utils"
        "github.com/prometheus/client_golang/prometheus"
)

var defaultSeverVersion = "1.0.0"

type Server struct {
        Name           string
        Version        string
        promReg        *prometheus.Registry
        ExitSignal     chan struct{}
        callback       sync.Once
        server         *http.Server
}

func NewServer(name, version string) *Server {
        if version == "" {
                version = defaultSeverVersion
        }
        s := &Server{
                Name:         name,
                Version:      version,
                promReg:      prometheus.NewRegistry(),
                ExitSignal:   make(chan struct{}),
        }
        return s
}

func (s *Server) Run() error {
        go utils.HandleSignals(s.Exit)
        logrus.Infof("%s sucessfully setup. SetUp running.", s.Name)

        logrus.Infof("Runing  %s", s.Name)
        if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
                logrus.Errorf("ListenAndServe Error: %s\n", err)
                return err
        }
        return nil
}

func (s *Server) PrintVersion() {
        logrus.Printf("%s version: %s\n", s.Name, s.Version)
}

func (s *Server) Stop() {
        logrus.Info("Stopping Server")
        logger.LogOutput("Shutting down server...")
        ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
        defer cancel()

        if err := s.server.Shutdown(ctx); err != nil {
                if ctx.Err() == context.DeadlineExceeded {
                        logrus.Warn("Server shutdown timed out")
                } else {
                        logrus.Errorf("Server Shutdown Error: %s", err)
                }
        } else {
                logrus.Info("Server gracefully stopped")
        }
}

func (s *Server) Exit() {
        s.callback.Do(func() {
                close(s.ExitSignal)
        })
}

