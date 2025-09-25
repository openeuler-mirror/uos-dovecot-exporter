// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT
package server

import (
	"os"
        "sync"
	"time"
	"context"
        "net/http"
        "uos-dovecot-exporter/pkg/utils"
	"uos-dovecot-exporter/pkg/logger"
	"uos-dovecot-exporter/internal/exporter"
	"gopkg.in/yaml.v2"
        "github.com/sirupsen/logrus"
	"github.com/alecthomas/kingpin"
	"github.com/dustin/go-humanize"
        "github.com/prometheus/client_golang/prometheus"
)

var defaultSeverVersion = "1.0.0"

type Server struct {
        Name           string
        Version        string
	CommonConfig   exporter.Config
        promReg        *prometheus.Registry
        ExitSignal     chan struct{}
        callback       sync.Once
	Error          error
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

func (s *Server) SetUp() error {
        defer func() {
                if s.Error != nil {
                        logrus.Errorf("SetUp error: %v", s.Error)
                }
        }()
        err := s.parse()
        if err != nil {
                logrus.Errorf("Parsing command line arguments failed: %v", err)
                return err
        }
        err = s.loadConfig()
        if err != nil {
                logrus.Errorf("Loading config file failed: %v", err)
                return err
        }
        err = s.setupLog()
        if err != nil {
                logrus.Errorf("SetUp error: %v", err)
                return err
        }

        err = s.setupHttpServer()
        if err != nil {
                logrus.Errorf("SetUp error: %v", err)
                return err
        }

        return nil
}

func (s *Server) loadConfig() error {
        content, err := os.ReadFile(*exporter.Configfile)
        if err != nil {
                logrus.Errorf("Failed to read config file: %v", err)
                logrus.Info("Use default config")
                return nil
        }
        err = yaml.Unmarshal(content, &s.CommonConfig)
        if err != nil {
                logrus.Errorf("Failed to parse config file: %v", err)
                logrus.Info("Use default config")
                return nil
        }
        logrus.Infof("Loaded config file from: %s", *exporter.Configfile)
        logrus.Info("CommonConfig file loaded")
        return nil
}

func (s *Server) setupLog() error {
        size, err := humanize.ParseBytes(s.CommonConfig.Logging.MaxSize)
        if err != nil {
                logrus.Errorf("Parsing log size failed: %v", err)
                return err
        }
        logConfig := logger.NewConfig(s.CommonConfig.Logging.Level, s.CommonConfig.Logging.LogPath, int64(size), s.CommonConfig.Logging.MaxAge)
        logger.Init(logConfig)
        return nil
}

func (s *Server) setupHttpServer() error {
        exporter.RegisterPrometheus(s.promReg)
	return nil
}


func (s *Server) Exit() {
        s.callback.Do(func() {
                close(s.ExitSignal)
        })
}

func (s *Server) parse() error {
        kingpin.Parse()
        return nil
}


