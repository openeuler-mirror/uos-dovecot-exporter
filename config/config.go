// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT
package config

import (
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
}

