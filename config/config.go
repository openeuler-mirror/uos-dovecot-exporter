// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT
package config

import (
	"uos-dovecot-exporter/pkg/utils"
        "github.com/alecthomas/kingpin"
        "github.com/sirupsen/logrus"
)

var (
        ScrapeUrl       *string
        Insecure        *bool
)

func init() {
        ScrapeUrl = kingpin.Flag("scrape_uri","Scrape URI").Short('s').String()

        Insecure = kingpin.Flag("insecure","Ignore server certificate if using https, Default: false.").Bool()

        if *ScrapeUrl != "" {
                if err := utils.ValidateURI(*ScrapeUrl); err != nil {
                        logrus.Warnf("Invalid scrape uri: %s", err)
                        logrus.Warnf("Use default scrape uri: %s", DefaultSettings.ScrapeUri)
                        *ScrapeUrl = DefaultSettings.ScrapeUri
                }
        }

        if *Insecure {
                logrus.Warn("Insecure mode enabled, this is not recommended for production use.")
        }

}

