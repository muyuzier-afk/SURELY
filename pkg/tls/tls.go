package tls

import (
	"crypto/tls"
)

// Config 定义 TLS 配置
type Config struct {
	TLSConfig *tls.Config
}

// NewTLSConfig 创建一个新的 TLS 配置
func NewTLSConfig() *Config {
	// 创建基础 TLS 配置
	tlsConfig := &tls.Config{
		// 只包含常见的密码套件
		CipherSuites: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		},
		NextProtos: []string{"h3"},
	}

	return &Config{
		TLSConfig: tlsConfig,
	}
}

// NewServerTLSConfig 创建一个新的服务器 TLS 配置
func NewServerTLSConfig(certFile, keyFile string) (*Config, error) {
	// 加载证书
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	// 创建基础 TLS 配置
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		// 只包含常见的密码套件
		CipherSuites: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		},
		NextProtos: []string{"h3"},
	}

	return &Config{
		TLSConfig: tlsConfig,
	}, nil
}
