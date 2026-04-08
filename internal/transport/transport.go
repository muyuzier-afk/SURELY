package transport

import (
	"context"
	"crypto/tls"
	"io"
	"log"
	"math/rand"
	"time"

	"github.com/quic-go/quic-go"
)

// Stream 简单的流接口
type Stream interface {
	io.Reader
	io.Writer
	io.Closer
}

// Config 传输层配置
type Config struct {
	TLSConfig  *tls.Config
	QUICConfig *quic.Config
}

// Server 传输层服务器
type Server struct {
	listener quic.Listener
	config   *Config
}

// Client 传输层客户端
type Client struct {
	conn quic.Connection
}

// NewServer 创建新的传输层服务器
func NewServer(addr string, config *Config) (*Server, error) {
	listener, err := quic.ListenAddr(addr, config.TLSConfig, config.QUICConfig)
	if err != nil {
		return nil, err
	}

	return &Server{
		listener: listener,
		config:   config,
	}, nil
}

// Accept 接受连接
func (s *Server) Accept(ctx context.Context) (*Client, error) {
	conn, err := s.listener.Accept(ctx)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
	}, nil
}

// Close 关闭服务器
func (s *Server) Close() error {
	return s.listener.Close()
}

// NewClient 创建新的传输层客户端
func NewClient(ctx context.Context, addr string, config *Config) (*Client, error) {
	conn, err := quic.DialAddr(ctx, addr, config.TLSConfig, config.QUICConfig)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
	}, nil
}

// OpenStream 打开新的流
func (c *Client) OpenStream() (Stream, error) {
	stream, err := c.conn.OpenStreamSync(context.Background())
	if err != nil {
		return nil, err
	}
	return &streamWrapper{stream: stream}, nil
}

// AcceptStream 接受新的流
func (c *Client) AcceptStream(ctx context.Context) (Stream, error) {
	stream, err := c.conn.AcceptStream(ctx)
	if err != nil {
		return nil, err
	}
	return &streamWrapper{stream: stream}, nil
}

// Close 关闭连接
func (c *Client) Close() error {
	return c.conn.CloseWithError(0, "")
}

// streamWrapper 包装 quic.Stream 以实现简单的接口
type streamWrapper struct {
	stream quic.Stream
}

func (w *streamWrapper) Read(p []byte) (n int, err error) {
	return w.stream.Read(p)
}

func (w *streamWrapper) Write(p []byte) (n int, err error) {
	return w.stream.Write(p)
}

func (w *streamWrapper) Close() error {
	return w.stream.Close()
}

// GeneratePadding 生成低熵填充数据
func GeneratePadding(targetLength int) []byte {
	padding := make([]byte, targetLength)
	for i := range padding {
		padding[i] = byte(i % 16)
	}
	return padding
}

// GetJitter 获取随机抖动时间
func GetJitter(min, max time.Duration) time.Duration {
	if min >= max {
		return min
	}
	jitter := rand.Int63n(int64(max-min)) + int64(min)
	return time.Duration(jitter)
}

// GenerateLowEntropyData 生成低熵数据
func GenerateLowEntropyData(length int) []byte {
	data := make([]byte, length)
	for i := range data {
		data[i] = byte(i % 256)
	}
	return data
}

// TLSConfigForServer 生成服务器 TLS 配置
func TLSConfigForServer(certPath, keyPath string, nextProtos []string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   nextProtos,
		MinVersion:   tls.VersionTLS13,
		MaxVersion:   tls.VersionTLS13,
		CipherSuites: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
		},
		CurvePreferences: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
			tls.CurveP384,
		},
	}, nil
}

// TLSConfigForClient 生成客户端 TLS 配置
func TLSConfigForClient(insecureSkipVerify bool, nextProtos []string) *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: insecureSkipVerify,
		NextProtos:         nextProtos,
		MinVersion:         tls.VersionTLS13,
		MaxVersion:         tls.VersionTLS13,
		CipherSuites: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
		},
		CurvePreferences: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
			tls.CurveP384,
		},
	}
}

// QUICConfig 生成 QUIC 配置
func QUICConfig(maxIdleTimeout time.Duration, enable0RTT bool) *quic.Config {
	return &quic.Config{
		MaxIdleTimeout: maxIdleTimeout,
		Allow0RTT:      enable0RTT,
	}
}

// Initialize 初始化传输层
func Initialize() {
	rand.Seed(time.Now().UnixNano())
	log.Println("Transport layer initialized")
}
