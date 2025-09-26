// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT
package server

import (
        "bytes"
        "net/http"
        "text/template"
)

type LandingPageConfig struct {
        CSS     string
        Name    string
        Links   []LandingPageLinks
        Version string
}

type LandingPageLinks struct {
        Address string
        Text    string
}

