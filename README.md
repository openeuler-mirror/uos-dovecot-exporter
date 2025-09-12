# UOS Dovecot Exporter for Prometheus

基于深度操作系统（UOS）的 Dovecot 邮件服务器性能监控导出器，用于从 Dovecot 邮件服务器抓取统计数据并将其导出为 Prometheus 指标。

## 项目简介

UOS Dovecot Exporter 是一个专门为深度操作系统（UOS）开发的 Dovecot 邮件服务器监控工具。它连接到 Dovecot 邮件服务器，收集IMAP4rev1、POP3 协议、SSL/TLS 加密等性能指标并以 Prometheus 格式导出，包括客户端和 Dovecot 邮件服务器活跃连接数、连接会话数，认证相关，邮箱操作，存储相关和各种系统资源使用的指标。

## 功能特性

- 🚀 **全面的指标收集**: 支持客户端/服务器连接/认证指标、邮箱操作、缓存性能等
- 📊 **Prometheus 兼容**: 原生支持 Prometheus 监控体系
- 🎯 **UOS 优化**: 专为深度操作系统环境优化
- ⚡ **高性能**: 低资源占用，高效稳定
- 🔧 **灵活配置**: 支持命令行参数和 YAML 配置文件

## 安装

### 从源码编译

```bash
git clone https://gitee.com/deepin-community/uos-dovecot-exporter.git
cd uos-dovecot-exporter
go build
```

### 二进制安装

从 [发布页面](https://gitee.com/deepin-community/uos-dovecot-exporter/releases) 下载适用于您系统的最新二进制文件。

## 使用方法

### 基本使用

```bash
./uos-dovecot-exporter 
```

### YAML 配置文件

也可以使用 YAML 配置文件：

```yaml

address: "0.0.0.0"
port: 9107
metricsPath: "/metrics"
log:
  level: "debug"
log_path: "/var/log/uos-exporter/dovecot-exporter.log"

```

## 监控指标

### 服务器连接指标

- Dovecot 服务是否在运行
- 当前活跃连接数
- 已连接会话数

### 系统资源信息

- CPU 使用率
- 虚拟内存使用情况
- 物理内存使用情况

## Prometheus 配置

在您的 `prometheus.yaml` 中添加以下配置：

```yaml
scrape_configs:
  - job_name: "uos-dovecot-exporter"
    static_configs:
      - targets: ["localhost:8090"]
```

## Dovecot 配置

为了允许导出器查询 Docecot 指标，请在您的 docecot.conf 中添加：

```
# 允许来自本地主机的缓存管理器访问
acl prometheus src 127.0.0.1
http_access allow manager prometheus
```

## 参与贡献

1. Fork 本仓库
2. 新建 feature 分支 (`git checkout -b feature/AmazingFeature`)
3. 提交修改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 提交 Pull Request

## 许可证

// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT
