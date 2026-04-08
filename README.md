# SURELY

A proxy protocol for Surely to counter the GFW censor model.

## Overview

SURELY is a secure proxy protocol designed to provide reliable and censorship-resistant network communication. It implements advanced encryption and traffic obfuscation techniques to evade detection and blocking mechanisms.

## Architecture

SURELY consists of two primary components:

- **Surely Server**: The server-side implementation that accepts and processes client connections
- **Surely Client SDK**: The client-side software development kit that provides programmatic access to the protocol

## Technical Features

### Transport Layer Security

- TLS 1.3 protocol support
- Modern cipher suite configuration
- Strict version enforcement
- X25519, CurveP256, CurveP384 elliptic curve support

### Network Protocol

- HTTP/3 (QUIC) transport
- Connection migration support
- 0-RTT handshake capability
- Configurable idle timeout

### Traffic Obfuscation

- Packet length distribution control
- Timing jitter management
- Background noise generation
- Padding optimization

## Configuration

### Server Configuration

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

### Client Configuration

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

## Installation

### From Source

Requires Go 1.25 or later.

```bash
git clone https://github.com/muyuzier-afk/SURELY.git
cd SURELY
go build -o surely-server ./cmd/surely-server
go build -o surely-client ./cmd/surely-client
```

### From Releases

Pre-built binaries are available on the GitHub Releases page.

## Usage

### Server

```bash
./surely-server -config config.toml
```

### Client

```bash
./surely-client -config client-config.toml -target /
```

## Client SDK

The Surely Client SDK provides a programmatic interface for integrating SURELY into applications.

### Example Usage

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

## Security Considerations

This software is provided for educational and research purposes only. Users are responsible for ensuring compliance with all applicable laws and regulations in their jurisdiction.

## Development

### Project Structure

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

### Building

```bash
go build -o surely-server ./cmd/surely-server
go build -o surely-client ./cmd/surely-client
```

## License

This project is licensed under the MIT License.

For more information, see the [LICENSE](LICENSE) file.
