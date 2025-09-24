// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT
package main

import (
        "uos-dovecot-exporter/internal/server"
        "uos-dovecot-exporter/pkg/logger"
        "github.com/sirupsen/logrus"
)

func Run(name string, version string) error {
        logger.InitDefaultLog()
        s := server.NewServer(name, version)
	s.PrintVersion()
        go func() {
                err := s.Run()
                if err != nil {
                        logrus.Errorf("Run error: %v", err)
			s.Error = err
                }
                s.Exit()
        }()
        select {
        case <-s.ExitSignal:
                s.Stop()
                logrus.Info("Exit exporter server completed")
                return s.Error
        }
}

