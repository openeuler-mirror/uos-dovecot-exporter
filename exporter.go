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
        go func() {
                err := s.Run()
                if err != nil {
                        logrus.Errorf("Run error: %v", err)
                }
                s.Exit()
        }()

        logrus.Info("init Run ....")
        return nil
}

