# SURELY

A proxy protocol for Surely to counter the GFW censor model.

## Language / Язык / 语言

- [English](#english)
- [Русский](#русский)
- [中文](#中文)

---

## English

### Overview

SURELY is a secure proxy protocol designed to provide reliable and censorship-resistant network communication. It implements advanced encryption and traffic obfuscation techniques to evade detection and blocking mechanisms.

### Architecture

SURELY consists of two primary components:

- **Surely Server**: The server-side implementation that accepts and processes client connections
- **Surely Client SDK**: The client-side software development kit that provides programmatic access to the protocol

### Technical Features

#### Transport Layer Security

- TLS 1.3 protocol support
- Modern cipher suite configuration
- Strict version enforcement
- X25519, CurveP256, CurveP384 elliptic curve support

#### Network Protocol

- HTTP/3 (QUIC) transport
- Connection migration support
- 0-RTT handshake capability
- Configurable idle timeout

#### Traffic Obfuscation

- Packet length distribution control
- Timing jitter management
- Background noise generation
- Padding optimization

### Configuration

#### Server Configuration

The server configuration is specified in a TOML file. The default configuration file is `config.toml`.

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

#### Client Configuration

The client configuration is specified in a TOML file. The default configuration file is `client-config.toml`.

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

### Installation

#### From Source

Requires Go 1.25 or later.

```bash
git clone https://github.com/muyuzier-afk/SURELY.git
cd SURELY
go build -o surely-server ./cmd/surely-server
go build -o surely-client ./cmd/surely-client
```

#### From Releases

Pre-built binaries are available on the GitHub Releases page.

### Usage

#### Server

```bash
./surely-server -config config.toml
```

#### Client

```bash
./surely-client -config client-config.toml -target /
```

### Client SDK

The Surely Client SDK provides a programmatic interface for integrating SURELY into applications.

#### Example Usage

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

### Security Considerations

This software is provided for educational and research purposes only. Users are responsible for ensuring compliance with all applicable laws and regulations in their jurisdiction.

### Development

#### Project Structure

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

#### Building

```bash
go build -o surely-server ./cmd/surely-server
go build -o surely-client ./cmd/surely-client
```

### License

This project is licensed under the MIT License.

For more information, see the [LICENSE](LICENSE) file.

---

## Русский

### Обзор

SURELY — это безопасный прокси-протокол, предназначенный для обеспечения надежной и устойчивой к цензуре сетевой коммуникации. Он реализует передовые методы шифрования и маскировки трафика для обхода механизмов обнаружения и блокировки.

### Архитектура

SURELY состоит из двух основных компонентов:

- **Surely Server**: Серверная реализация, которая принимает и обрабатывает подключения клиентов
- **Surely Client SDK**: Клиентский программный набор для интеграции протокола в приложения

### Технические особенности

#### Безопасность транспортного уровня

- Поддержка протокола TLS 1.3
- Современная конфигурация наборов шифрования
- Строгое соблюдение версии протокола
- Поддержка эллиптических кривых X25519, CurveP256, CurveP384

#### Сетевой протокол

- Транспорт HTTP/3 (QUIC)
- Поддержка миграции соединений
- Возможность рукопожатия 0-RTT
- Настраиваемый таймаут простоя

#### Маскировка трафика

- Управление распределением длины пакетов
- Управление задержками
- Генерация фонового шума
- Оптимизация дополнения пакетов

### Конфигурация

#### Конфигурация сервера

Конфигурация сервера задается в файле TOML. Стандартный файл конфигурации — `config.toml`.

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

#### Конфигурация клиента

Конфигурация клиента задается в файле TOML. Стандартный файл конфигурации — `client-config.toml`.

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

### Установка

#### Из исходного кода

Требуется Go 1.25 или новее.

```bash
git clone https://github.com/muyuzier-afk/SURELY.git
cd SURELY
go build -o surely-server ./cmd/surely-server
go build -o surely-client ./cmd/surely-client
```

#### Из релизов

Предварительно собранные двоичные файлы доступны на странице GitHub Releases.

### Использование

#### Сервер

```bash
./surely-server -config config.toml
```

#### Клиент

```bash
./surely-client -config client-config.toml -target /
```

### Клиентский SDK

Surely Client SDK предоставляет программный интерфейс для интеграции SURELY в приложения.

#### Пример использования

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

### Безопасные соображения

Данное программное обеспечение предоставляется исключительно в образовательных и исследовательских целях. Пользователи несут ответственность за соблюдение всех применимых законов и нормативных актов в своей юрисдикции.

### Разработка

#### Структура проекта

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

#### Сборка

```bash
go build -o surely-server ./cmd/surely-server
go build -o surely-client ./cmd/surely-client
```

### Лицензия

Этот проект лицензирован под MIT License.

Для получения более подробной информации см. файл [LICENSE](LICENSE).

---

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
