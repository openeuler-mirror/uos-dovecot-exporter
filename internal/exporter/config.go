// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT
package exporter
  
import (
	"time"
        "uos-dovecot-exporter/pkg/logger"
)

var (
        Configfile    *string
        DefaultConfig = Config{
                Logging: logger.Config{
                        Level:   "debug",
                        LogPath: "/var/log/uos-exporter/dovecot-exporter.log",
                        MaxSize: "10MB",
                        MaxAge:  time.Hour * 24 * 7},
                Address:     "0.0.0.0",
                Port:        9107,
                MetricsPath: "/metrics",
        }
)

type Config struct {
        Logging     logger.Config `yaml:"log"`
        Address     string        `yaml:"address"`
        Port        int           `yaml:"port"`
        MetricsPath string        `yaml:"metricsPath"`
}

