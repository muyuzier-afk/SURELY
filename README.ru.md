# SURELY

Прокси-протокол Surely для противодействия модели цензуры GFW.

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

## Другие языки / Other Languages / 其他语言

- [English](../README.md)
- [Русский](README.ru.md)
- [中文](../README.zh.md)
