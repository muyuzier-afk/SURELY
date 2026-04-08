package masque

import (
	"context"

	"github.com/quic-go/quic-go"
)

// Client 定义 MASQUE 客户端
type Client struct {
	conn quic.Connection
}

// Server 定义 MASQUE 服务器
type Server struct {
	conn quic.Connection
}

// NewClient 创建一个新的 MASQUE 客户端
func NewClient(conn quic.Connection) *Client {
	return &Client{
		conn: conn,
	}
}

// NewServer 创建一个新的 MASQUE 服务器
func NewServer(conn quic.Connection) *Server {
	return &Server{
		conn: conn,
	}
}

// Connect 建立 CONNECT 连接
func (c *Client) Connect(target string) (quic.Stream, error) {
	// 创建双向流
	stream, err := c.conn.OpenStreamSync(context.Background())
	if err != nil {
		return nil, err
	}

	// 发送 CONNECT 请求
	_, err = stream.Write([]byte("CONNECT " + target + " HTTP/3\r\n\r\n"))
	if err != nil {
		stream.Close()
		return nil, err
	}

	return stream, nil
}

// AcceptStream 接受流
func (s *Server) AcceptStream(ctx context.Context) (quic.Stream, error) {
	return s.conn.AcceptStream(ctx)
}
