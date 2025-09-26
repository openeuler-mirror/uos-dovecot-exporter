// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT
package server

import (
        "net/http"
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

type LandingPageHandler struct {
        landingPage []byte
}

func (h *LandingPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/" {
                http.NotFound(w, r)
                return
        }

        w.Header().Set("Content-Type", "text/html; charset=UTF-8")
        w.Write(h.landingPage)
}

