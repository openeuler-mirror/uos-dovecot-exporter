// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT

package logger

import (
        "os"
        "time"
)

type FileRotator struct {
        basePath  string
        maxSize   int64
        maxAge    time.Duration
        current   *os.File
        size      int64
        startTime time.Time
        keepFiles int
}
