// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT

package logger

import (
	"github.com/sirupsen/logrus"
	"strings"
)

type fileLogConfig struct {
	level       string
}

func Init(config fileLogConfig) {
	switch level := strings.ToLower(config.level); level {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	default:
		logrus.SetLevel(logrus.WarnLevel)
		logrus.Warnf("unknown log level: %s, use default level: warn", level)
		logrus.Warnf("support level is [debug,info,warn]")
	}
}

func InitDefaultLog() {
	Init(fileLogConfig{level: "info",})
	logrus.SetFormatter(&logrus.TextFormatter{DisableTimestamp: true,})
}
