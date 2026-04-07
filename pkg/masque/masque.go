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
	listener *quic.Listener
}

// NewClient 创建一个新的 MASQUE 客户端
func NewClient(conn quic.Connection) *Client {
	return &Client{
		conn: conn,
	}
}

// NewServer 创建一个新的 MASQUE 服务器
func NewServer(listener *quic.Listener) *Server {
	return &Server{
		listener: listener,
	}
}

// Connect 建立 CONNECT 连接
func (c *Client) Connect(target string) (quic.Stream, error) {
	// 创建双向流
	stream, err := c.conn.OpenStream()
	if err != nil {
		return nil, err
	}

	// 发送 CONNECT 请求
	// 这里简化实现，实际应该遵循 MASQUE 规范
	_, err = stream.Write([]byte("CONNECT " + target + " HTTP/3\r\n\r\n"))
	if err != nil {
		return nil, err
	}

	return stream, nil
}

// Accept 接受 MASQUE 连接
func (s *Server) Accept(ctx context.Context) (*Client, error) {
	// 接受 QUIC 连接
	conn, err := s.listener.Accept(ctx)
	if err != nil {
		return nil, err
	}

	return NewClient(conn), nil
}

// HandleStream 处理流
func (s *Server) HandleStream(stream quic.Stream) error {
	// 读取 CONNECT 请求
	buf := make([]byte, 1024)
	_, err := stream.Read(buf)
	if err != nil {
		return err
	}

	// 解析请求
	// 这里简化实现，实际应该遵循 MASQUE 规范
	// 提取目标地址

	// 发送响应
	_, err = stream.Write([]byte("HTTP/3 200 OK\r\n\r\n"))
	if err != nil {
		return err
	}

	// 处理数据流
	// 这里简化实现，实际应该根据目标地址建立连接并转发数据
	go func() {
		// 读取数据并丢弃
		buf := make([]byte, 1024)
		for {
			_, err := stream.Read(buf)
			if err != nil {
				break
			}
		}
		stream.Close()
	}()

	return nil
}
