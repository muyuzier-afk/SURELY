# SURELY

Surely 的代理协议，用于对抗 GFW 审查模型。

## 中文

### 概述

SURELY 是一种安全的代理协议，旨在提供可靠且抗审查的网络通信。它实现了先进的加密和流量混淆技术，以规避检测和阻塞机制。

### 架构

SURELY 由两个主要组件组成：

- **Surely Server**：服务器端实现，接受并处理客户端连接
- **Surely Client SDK**：客户端软件开发工具包，提供协议的编程访问接口

### 技术特性

#### 传输层安全

- TLS 1.3 协议支持
- 现代密码套件配置
- 严格的版本执行
- X25519、CurveP256、CurveP384 椭圆曲线支持

#### 网络协议

- HTTP/3 (QUIC) 传输
- 连接迁移支持
- 0-RTT 握手能力
- 可配置的空闲超时

#### 流量混淆

- 数据包长度分布控制
- 时序抖动管理
- 背景噪声生成
- 填充优化

### 配置

#### 服务器配置

服务器配置在 TOML 文件中指定。默认配置文件为 `config.toml`。

```toml
[server]
listen_addr = ":8443"
tls_cert_path = "server.crt"
tls_key_path = "server.key"
enable_http3 = true

[tls]
min_version = "TLSv1.3"
max_version = "TLSv1.3"
next_protos = ["h3"]

[quic]
max_idle_timeout = "30s"
enable_0rtt = true
connection_migration = true

[traffic]
packet_length_distribution = [64, 128, 256, 512, 1024, 1500]
jitter_min = "5ms"
jitter_max = "50ms"
background_noise = true
noise_interval = "30s"

[masque]
enable_connect = true
enable_datagram = true
```

#### 客户端配置

客户端配置在 TOML 文件中指定。默认配置文件为 `client-config.toml`。

```toml
[client]
server_addr = "localhost:8443"
insecure_skip_verify = true

[tls]
min_version = "TLSv1.3"
max_version = "TLSv1.3"
next_protos = ["h3"]

[quic]
max_idle_timeout = "30s"
enable_0rtt = true

[traffic]
enable_padding = true
enable_jitter = true
jitter_min = "5ms"
jitter_max = "50ms"
```

### 安装

#### 从源码安装

需要 Go 1.25 或更高版本。

```bash
git clone https://github.com/muyuzier-afk/SURELY.git
cd SURELY
go build -o surely-server ./cmd/surely-server
go build -o surely-client ./cmd/surely-client
```

#### 从发布版安装

预编译的二进制文件可在 GitHub Releases 页面获取。

### 使用

#### 服务器

```bash
./surely-server -config config.toml
```

#### 客户端

```bash
./surely-client -config client-config.toml -target /
```

### 客户端 SDK

Surely Client SDK 提供了一个编程接口，用于将 SURELY 集成到应用程序中。

#### 使用示例

```go
package main

import (
	"log"
	"surely/pkg/surely-sdk"
)

func main() {
	client, err := surely_sdk.NewClient("https://localhost:8443", true)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	resp, err := client.Get("/")
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	log.Printf("Response status: %s", resp.Status)
}
```

### 安全考虑

本软件仅用于教育和研究目的。用户有责任确保遵守其管辖区内的所有适用法律和法规。

### 开发

#### 项目结构

```
SURELY/
├── cmd/
│   ├── surely-client/
│   │   └── main.go
│   └── surely-server/
│       └── main.go
├── internal/
│   ├── config/
│   ├── protocol/
│   └── transport/
├── pkg/
│   └── surely-sdk/
├── config.toml
├── client-config.toml
└── ver
```

#### 构建

```bash
go build -o surely-server ./cmd/surely-server
go build -o surely-client ./cmd/surely-client
```

### 许可证

本项目采用 MIT License 许可证。

更多信息，请参阅 [LICENSE](LICENSE) 文件。

---

## 其他语言 / Other Languages / Другие языки

- [English](../README.md)
- [Русский](../README.ru.md)
- [中文](README.zh.md)
