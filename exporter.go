// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT
package main

import (
	"uos-dovecot-exporter/pkg/logger"
	"github.com/sirupsen/logrus"
)

func Run(name string, version string) error {
	logger.InitDefaultLog()

	logrus.Info("init Run ....")
	return nil
}
