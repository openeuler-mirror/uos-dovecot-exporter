// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT
package server

import (
	"os"
	"fmt"
        "sync"
	"time"
	"bytes"
	"context"
        "net/http"
	"encoding/json"
	"uos-dovecot-exporter/config"
        "uos-dovecot-exporter/pkg/utils"
	"uos-dovecot-exporter/pkg/logger"
	"uos-dovecot-exporter/pkg/ratelimit"
	"uos-dovecot-exporter/internal/exporter"
	"gopkg.in/yaml.v2"
        "github.com/sirupsen/logrus"
	"github.com/alecthomas/kingpin"
	"github.com/dustin/go-humanize"
        "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var defaultSeverVersion = "1.0.0"

type Server struct {
        Name           string
        Version        string
        CommonConfig   exporter.Config
        promReg        *prometheus.Registry
        handlers       []HandlerFunc
        ExitSignal     chan struct{}
        Error          error
        callback       sync.Once
        ExporterConfig config.Settings
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

func (s *Server) healthzHandler(w http.ResponseWriter, r *http.Request) {
        // 构造健康检查响应
        type healthzResponse struct {
                Status  string `json:"status"`
                Message string `json:"message"`
        }
        response := healthzResponse{
                Status:  "ok",
                Message: fmt.Sprintf("%s is running normally.", s.getName()),
        }

        // 设置响应头为 JSON 格式
        w.Header().Set("Content-Type", "application/json")

        // 使用缓冲区编码 JSON 数据，避免部分写入问题
        buf := new(bytes.Buffer)
        if err := json.NewEncoder(buf).Encode(response); err != nil {
                // 记录详细的错误日志，包括请求上下文
                logrus.WithFields(logrus.Fields{
                        "method": r.Method,
                        "path":   r.URL.Path,
                        "error":  err,
                }).Error("Failed to encode healthz response")
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
        }

        // 写入状态码并发送响应体
        w.WriteHeader(http.StatusOK)
        if _, err := buf.WriteTo(w); err != nil {
                // 记录写入失败的日志
                logrus.WithFields(logrus.Fields{
                        "method": r.Method,
                        "path":   r.URL.Path,
                        "error":  err,
                }).Error("Failed to write healthz response to client")
        }
}

// 获取 Name 字段的线程安全方法
func (s *Server) getName() string {
        // s.mu.RLock()
        // defer s.mu.RUnlock()
        return s.Name
}

func (s *Server) setupHttpServer() error {
        exporter.RegisterPrometheus(s.promReg)
        mux := http.NewServeMux()

        // 注册健康检查接口
        mux.HandleFunc("/healthz", s.healthzHandler)

        // 原有的路由注册逻辑
        mux.Handle(s.CommonConfig.MetricsPath, s)
        if *UseRatelimit {
                rateLimiter, err := ratelimit.NewRateLimiter(*rateLimitInterval, *rateLimitSize)
                if err != nil {
                        logrus.Errorf("ratelimit middleware init error: %v", err)
                }
                s.Use(Ratelimit(rateLimiter))
        }

	return nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
        promhttp.HandlerFor(s.promReg, promhttp.HandlerOpts{}).ServeHTTP(w, r)
}

func (s *Server) Use(handlerFuncs ...HandlerFunc) {
        s.handlers = append(s.handlers, handlerFuncs...)
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


