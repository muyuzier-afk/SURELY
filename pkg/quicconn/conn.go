package quicconn

import (
	"context"
	"crypto/tls"

	"github.com/quic-go/quic-go"
)

// Config 定义 QUIC 连接配置
type Config struct {
	TLSConfig *tls.Config
	QuicConfig *quic.Config
}

// NewClientConn 创建一个新的 QUIC 客户端连接
func NewClientConn(ctx context.Context, addr string, cfg *Config) (quic.Connection, error) {
	// 创建 QUIC 连接
	conn, err := quic.DialAddr(ctx, addr, cfg.TLSConfig, cfg.QuicConfig)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// NewServerConn 创建一个新的 QUIC 服务器连接
func NewServerConn(ctx context.Context, listener *quic.Listener) (quic.Connection, error) {
	// 接受连接
	conn, err := listener.Accept(ctx)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// NewListener 创建一个新的 QUIC 监听器
func NewListener(addr string, cfg *Config) (*quic.Listener, error) {
	// 创建监听器
	listener, err := quic.ListenAddr(addr, cfg.TLSConfig, cfg.QuicConfig)
	if err != nil {
		return nil, err
	}

	return listener, nil
}
