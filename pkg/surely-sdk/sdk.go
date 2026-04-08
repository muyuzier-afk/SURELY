package surely_sdk

import (
	"crypto/tls"
	"io"
	"net/http"
	"net/url"
)

// Client Surely 客户端 SDK
type Client struct {
	httpClient   *http.Client
	serverAddr   string
	insecureSkip bool
}

// NewClient 创建新的 Surely 客户端
func NewClient(serverAddr string, insecureSkipVerify bool) (*Client, error) {
	// 验证服务器地址
	u, err := url.Parse(serverAddr)
	if err != nil {
		return nil, err
	}

	// 创建 HTTP 客户端（暂时使用标准 HTTP 客户端，简化实现）
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: insecureSkipVerify,
				NextProtos:         []string{"h3", "http/1.1"},
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
			},
		},
	}

	return &Client{
		httpClient:   httpClient,
		serverAddr:   u.String(),
		insecureSkip: insecureSkipVerify,
	}, nil
}

// Get 发送 GET 请求
func (c *Client) Get(path string) (*http.Response, error) {
	fullURL := c.serverAddr + path
	return c.httpClient.Get(fullURL)
}

// Post 发送 POST 请求
func (c *Client) Post(path, contentType string, body io.Reader) (*http.Response, error) {
	fullURL := c.serverAddr + path
	return c.httpClient.Post(fullURL, contentType, body)
}

// Do 发送自定义请求
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.httpClient.Do(req)
}

// Close 关闭客户端
func (c *Client) Close() error {
	if c.httpClient != nil && c.httpClient.Transport != nil {
		if transport, ok := c.httpClient.Transport.(*http.Transport); ok {
			transport.CloseIdleConnections()
		}
	}
	return nil
}

// Version 获取 SDK 版本
func Version() string {
	return "1.0.0"
}

