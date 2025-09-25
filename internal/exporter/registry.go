// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT
package exporter

import (
	"sync"
        "github.com/prometheus/client_golang/prometheus"
)

var defaultReg *Registry

type Registry struct {
        metrics []Metric
        mu      sync.RWMutex
}

func RegisterPrometheus(reg *prometheus.Registry) {
        reg.MustRegister(defaultReg)
}

func (r *Registry) GetMetrics() []Metric {
        r.mu.RLock()
        defer r.mu.RUnlock()
        return r.metrics
}

func (r *Registry) Describe(descs chan<- *prometheus.Desc) {
}

func (r *Registry) Collect(ch chan<- prometheus.Metric) {
        for _, m := range r.GetMetrics() {
                m.Collect(ch)
        }
}

