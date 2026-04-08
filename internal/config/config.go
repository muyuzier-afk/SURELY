package config

import (
	"time"

	"github.com/BurntSushi/toml"
)

// ServerConfig 服务器配置
type ServerConfig struct {
	Server    ServerSection    `toml:"server"`
	TLS       TLSSection       `toml:"tls"`
	QUIC      QUICSection      `toml:"quic"`
	Traffic   TrafficSection   `toml:"traffic"`
	MASQUE    MASQUESection    `toml:"masque"`
}

// ClientConfig 客户端配置
type ClientConfig struct {
	Client   ClientSection    `toml:"client"`
	TLS      TLSSection       `toml:"tls"`
	QUIC     QUICSection      `toml:"quic"`
	Traffic  TrafficSection   `toml:"traffic"`
}

// ServerSection 服务器配置段
type ServerSection struct {
	ListenAddr   string `toml:"listen_addr"`
	TLSCertPath  string `toml:"tls_cert_path"`
	TLSKeyPath   string `toml:"tls_key_path"`
	EnableHTTP3  bool   `toml:"enable_http3"`
}

// ClientSection 客户端配置段
type ClientSection struct {
	ServerAddr          string `toml:"server_addr"`
	InsecureSkipVerify  bool   `toml:"insecure_skip_verify"`
}

// TLSSection TLS 配置段
type TLSSection struct {
	MinVersion  string   `toml:"min_version"`
	MaxVersion  string   `toml:"max_version"`
	NextProtos  []string `toml:"next_protos"`
}

// QUICSection QUIC 配置段
type QUICSection struct {
	MaxIdleTimeout    string `toml:"max_idle_timeout"`
	Enable0RTT        bool   `toml:"enable_0rtt"`
	ConnectionMigration bool `toml:"connection_migration,omitempty"`
}

// TrafficSection 流量配置段
type TrafficSection struct {
	PacketLengthDistribution []int   `toml:"packet_length_distribution,omitempty"`
	JitterMin              string  `toml:"jitter_min"`
	JitterMax              string  `toml:"jitter_max"`
	BackgroundNoise        bool    `toml:"background_noise,omitempty"`
	NoiseInterval          string  `toml:"noise_interval,omitempty"`
	EnablePadding          bool    `toml:"enable_padding,omitempty"`
	EnableJitter           bool    `toml:"enable_jitter,omitempty"`
}

// MASQUESection MASQUE 配置段
type MASQUESection struct {
	EnableConnect  bool `toml:"enable_connect"`
	EnableDatagram bool `toml:"enable_datagram"`
}

// LoadServerConfig 加载服务器配置
func LoadServerConfig(path string) (*ServerConfig, error) {
	var config ServerConfig
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// LoadClientConfig 加载客户端配置
func LoadClientConfig(path string) (*ClientConfig, error) {
	var config ClientConfig
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// ParseDuration 解析持续时间字符串
func ParseDuration(s string) time.Duration {
	if s == "" {
		return 0
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0
	}
	return d
}

// DefaultServerConfig 获取默认服务器配置
func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Server: ServerSection{
			ListenAddr:   ":8443",
			TLSCertPath:  "server.crt",
			TLSKeyPath:   "server.key",
			EnableHTTP3:  true,
		},
		TLS: TLSSection{
			MinVersion: "TLSv1.3",
			MaxVersion: "TLSv1.3",
			NextProtos: []string{"h3"},
		},
		QUIC: QUICSection{
			MaxIdleTimeout:      "30s",
			Enable0RTT:          true,
			ConnectionMigration: true,
		},
		Traffic: TrafficSection{
			PacketLengthDistribution: []int{64, 128, 256, 512, 1024, 1500},
			JitterMin:              "5ms",
			JitterMax:              "50ms",
			BackgroundNoise:        true,
			NoiseInterval:          "30s",
		},
		MASQUE: MASQUESection{
			EnableConnect:  true,
			EnableDatagram: true,
		},
	}
}

// DefaultClientConfig 获取默认客户端配置
func DefaultClientConfig() *ClientConfig {
	return &ClientConfig{
		Client: ClientSection{
			ServerAddr:         "localhost:8443",
			InsecureSkipVerify: true,
		},
		TLS: TLSSection{
			MinVersion: "TLSv1.3",
			MaxVersion: "TLSv1.3",
			NextProtos: []string{"h3"},
		},
		QUIC: QUICSection{
			MaxIdleTimeout: "30s",
			Enable0RTT:     true,
		},
		Traffic: TrafficSection{
			EnablePadding: true,
			EnableJitter:  true,
			JitterMin:     "5ms",
			JitterMax:     "50ms",
		},
	}
}
