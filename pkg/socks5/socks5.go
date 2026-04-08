package socks5

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

// ConnectionHandler 定义连接处理函数类型
type ConnectionHandler func(conn net.Conn, target string) error

// Server SOCKS5 服务器
type Server struct {
	addr    string
	handler ConnectionHandler
}

// NewServer 创建一个新的 SOCKS5 服务器
func NewServer(addr string, handler ConnectionHandler) *Server {
	return &Server{
		addr:    addr,
		handler: handler,
	}
}

// ListenAndServe 启动 SOCKS5 服务器
func (s *Server) ListenAndServe() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept error: %v", err)
			continue
		}

		go s.handleConnection(conn)
	}
}

// handleConnection 处理 SOCKS5 连接
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	// 读取客户端问候
	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	if err != nil {
		log.Printf("Read greeting error: %v", err)
		return
	}

	// 验证 SOCKS5 版本
	if buf[0] != 0x05 {
		log.Printf("Invalid SOCKS version: %d", buf[0])
		return
	}

	// 支持的认证方法数量
	nmethods := buf[1]
	methods := buf[2 : 2+nmethods]

	// 选择无认证方法
	hasNoAuth := false
	for _, method := range methods {
		if method == 0x00 {
			hasNoAuth = true
			break
		}
	}

	if !hasNoAuth {
		// 不支持的认证方法
		conn.Write([]byte{0x05, 0xFF})
		return
	}

	// 发送认证确认
	conn.Write([]byte{0x05, 0x00})

	// 读取请求
	n, err = conn.Read(buf)
	if err != nil {
		log.Printf("Read request error: %v", err)
		return
	}

	// 验证请求版本
	if buf[0] != 0x05 {
		log.Printf("Invalid SOCKS version: %d", buf[0])
		return
	}

	// 只支持 CONNECT 命令
	if buf[1] != 0x01 {
		log.Printf("Unsupported command: %d", buf[1])
		conn.Write([]byte{0x05, 0x07, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		return
	}

	// 解析地址
	target, err := parseAddress(buf[3:], n-3)
	if err != nil {
		log.Printf("Parse address error: %v", err)
		conn.Write([]byte{0x05, 0x08, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		return
	}

	// 发送成功响应
	conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	// 调用处理函数
	err = s.handler(conn, target)
	if err != nil {
		log.Printf("Handler error: %v", err)
	}
}

// parseAddress 解析 SOCKS5 地址
func parseAddress(data []byte, length int) (string, error) {
	if len(data) < 1 {
		return "", fmt.Errorf("invalid address data")
	}

	addrType := data[0]
	switch addrType {
	case 0x01: // IPv4
		if len(data) < 7 {
			return "", fmt.Errorf("invalid IPv4 address")
		}
		ip := net.IP(data[1:5])
		port := binary.BigEndian.Uint16(data[5:7])
		return fmt.Sprintf("%s:%d", ip, port), nil

	case 0x03: // 域名
		if len(data) < 2 {
			return "", fmt.Errorf("invalid domain address")
		}
		domainLen := data[1]
		if len(data) < int(domainLen)+3 {
			return "", fmt.Errorf("invalid domain address length")
		}
		domain := string(data[2 : 2+domainLen])
		port := binary.BigEndian.Uint16(data[2+domainLen:])
		return fmt.Sprintf("%s:%d", domain, port), nil

	case 0x04: // IPv6
		if len(data) < 19 {
			return "", fmt.Errorf("invalid IPv6 address")
		}
		ip := net.IP(data[1:17])
		port := binary.BigEndian.Uint16(data[17:19])
		return fmt.Sprintf("%s:%d", ip, port), nil

	default:
		return "", fmt.Errorf("unsupported address type: %d", addrType)
	}
}