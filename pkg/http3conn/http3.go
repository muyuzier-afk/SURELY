package http3conn

import (
	"crypto/tls"
	"net/http"

	"github.com/quic-go/quic-go/http3"
)

// Client 定义 HTTP/3 客户端
type Client struct {
	client *http3.Client
}

// Server 定义 HTTP/3 服务器
type Server struct {
	server *http3.Server
}

// NewClient 创建一个新的 HTTP/3 客户端
func NewClient(tlsConfig *tls.Config) *Client {
	client := &http3.Client{
		TLSClientConfig: tlsConfig,
	}

	return &Client{
		client: client,
	}
}

// NewServer 创建一个新的 HTTP/3 服务器
func NewServer(addr string, tlsConfig *tls.Config, handler http.Handler) *Server {
	server := &http3.Server{
		Addr:      addr,
		TLSConfig: tlsConfig,
		Handler:   handler,
	}

	return &Server{
		server: server,
	}
}

// Get 发送 GET 请求
func (c *Client) Get(url string) (*http.Response, error) {
	return c.client.Get(url)
}

// ListenAndServe 启动 HTTP/3 服务器
func (s *Server) ListenAndServe() error {
	return s.server.ListenAndServe()
}
