package protocol

import (
	"bufio"
	"context"
	"encoding/binary"
	"io"
	"log"
	"net"

	"github.com/quic-go/quic-go"
)

// Version 协议版本
const Version = "1.0"

// MessageType 消息类型
type MessageType uint8

const (
	MessageTypeConnect    MessageType = 0x01
	MessageTypeConnectAck MessageType = 0x02
	MessageTypeData       MessageType = 0x03
	MessageTypePing       MessageType = 0x04
	MessageTypePong       MessageType = 0x05
)

// Message 协议消息
type Message struct {
	Type    MessageType
	Payload []byte
}

// Transport 传输层接口
type Transport interface {
	OpenStream() (quic.Stream, error)
	AcceptStream(ctx context.Context) (quic.Stream, error)
	Close() error
}

// Server 协议服务器
type Server struct {
	transport Transport
}

// Client 协议客户端
type Client struct {
	transport Transport
}

// NewServer 创建新的协议服务器
func NewServer(t Transport) *Server {
	return &Server{
		transport: t,
	}
}

// NewClient 创建新的协议客户端
func NewClient(t Transport) *Client {
	return &Client{
		transport: t,
	}
}

// HandleConnection 处理连接
func (s *Server) HandleConnection(ctx context.Context) error {
	for {
		stream, err := s.transport.AcceptStream(ctx)
		if err != nil {
			return err
		}

		go s.HandleStream(stream)
	}
}

// HandleStream 处理流
func (s *Server) HandleStream(stream quic.Stream) {
	defer stream.Close()

	reader := bufio.NewReader(stream)
	writer := bufio.NewWriter(stream)

	for {
		msg, err := ReadMessage(reader)
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading message: %v", err)
			}
			return
		}

		switch msg.Type {
		case MessageTypeConnect:
			s.handleConnect(writer, stream, msg.Payload)
		case MessageTypeData:
			s.handleData(stream, msg.Payload)
		case MessageTypePing:
			s.handlePing(writer)
		default:
			log.Printf("Unknown message type: %d", msg.Type)
		}
	}
}

// handleConnect 处理连接请求
func (s *Server) handleConnect(writer *bufio.Writer, stream quic.Stream, payload []byte) {
	target := string(payload)
	log.Printf("Connect request to: %s", target)

	conn, err := net.Dial("tcp", target)
	if err != nil {
		log.Printf("Failed to connect to %s: %v", target, err)
		return
	}
	defer conn.Close()

	ackMsg := &Message{
		Type:    MessageTypeConnectAck,
		Payload: []byte("OK"),
	}

	if err := WriteMessage(writer, ackMsg); err != nil {
		log.Printf("Failed to send connect ack: %v", err)
		return
	}

	writer.Flush()

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				break
			}
			dataMsg := &Message{
				Type:    MessageTypeData,
				Payload: buf[:n],
			}
			if err := WriteMessage(writer, dataMsg); err != nil {
				break
			}
			writer.Flush()
		}
	}()

	buf := make([]byte, 4096)
	for {
		msg, err := ReadMessage(reader)
		if err != nil {
			break
		}
		if msg.Type == MessageTypeData {
			_, err = conn.Write(msg.Payload)
			if err != nil {
				break
			}
		}
	}
}

// handleData 处理数据
func (s *Server) handleData(stream quic.Stream, payload []byte) {
	if _, err := stream.Write(payload); err != nil {
		log.Printf("Error writing data: %v", err)
	}
}

// handlePing 处理 ping
func (s *Server) handlePing(writer *bufio.Writer) {
	pongMsg := &Message{
		Type:    MessageTypePong,
		Payload: nil,
	}

	if err := WriteMessage(writer, pongMsg); err != nil {
		log.Printf("Failed to send pong: %v", err)
	}

	writer.Flush()
}

// Connect 连接到服务器
func (c *Client) Connect(target string) error {
	stream, err := c.transport.OpenStream()
	if err != nil {
		return err
	}
	defer stream.Close()

	connectMsg := &Message{
		Type:    MessageTypeConnect,
		Payload: []byte(target),
	}

	writer := bufio.NewWriter(stream)
	if err := WriteMessage(writer, connectMsg); err != nil {
		return err
	}

	writer.Flush()

	reader := bufio.NewReader(stream)
	msg, err := ReadMessage(reader)
	if err != nil {
		return err
	}

	if msg.Type != MessageTypeConnectAck {
		return io.ErrUnexpectedEOF
	}

	return nil
}

// SendData 发送数据
func (c *Client) SendData(data []byte) error {
	stream, err := c.transport.OpenStream()
	if err != nil {
		return err
	}
	defer stream.Close()

	dataMsg := &Message{
		Type:    MessageTypeData,
		Payload: data,
	}

	writer := bufio.NewWriter(stream)
	if err := WriteMessage(writer, dataMsg); err != nil {
		return err
	}

	return writer.Flush()
}

// Ping 发送 ping
func (c *Client) Ping() error {
	stream, err := c.transport.OpenStream()
	if err != nil {
		return err
	}
	defer stream.Close()

	pingMsg := &Message{
		Type:    MessageTypePing,
		Payload: nil,
	}

	writer := bufio.NewWriter(stream)
	if err := WriteMessage(writer, pingMsg); err != nil {
		return err
	}

	writer.Flush()

	reader := bufio.NewReader(stream)
	msg, err := ReadMessage(reader)
	if err != nil {
		return err
	}

	if msg.Type != MessageTypePong {
		return io.ErrUnexpectedEOF
	}

	return nil
}

// Close 关闭连接
func (c *Client) Close() error {
	return c.transport.Close()
}

// ReadMessage 读取消息
func ReadMessage(reader *bufio.Reader) (*Message, error) {
	header := make([]byte, 5)
	if _, err := io.ReadFull(reader, header); err != nil {
		return nil, err
	}

	msgType := MessageType(header[0])
	payloadLen := binary.BigEndian.Uint32(header[1:5])

	payload := make([]byte, payloadLen)
	if _, err := io.ReadFull(reader, payload); err != nil {
		return nil, err
	}

	return &Message{
		Type:    msgType,
		Payload: payload,
	}, nil
}

// WriteMessage 写入消息
func WriteMessage(writer *bufio.Writer, msg *Message) error {
	header := make([]byte, 5)
	header[0] = byte(msg.Type)
	binary.BigEndian.PutUint32(header[1:5], uint32(len(msg.Payload)))

	if _, err := writer.Write(header); err != nil {
		return err
	}

	if _, err := writer.Write(msg.Payload); err != nil {
		return err
	}

	return nil
}
