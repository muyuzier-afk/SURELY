package main

import (
	"context"
	"crypto/tls"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"surely/internal/config"

	"github.com/quic-go/quic-go"
)

func main() {
	configPath := flag.String("config", "config.toml", "Path to configuration file")
	flag.Parse()

	log.Println("Surely Server v1.0.1 starting...")

	cfg, err := config.LoadServerConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Printf("Loaded configuration from %s", *configPath)

	// 加载 TLS 证书
	cert, err := tls.LoadX509KeyPair(cfg.Server.TLSCertPath, cfg.Server.TLSKeyPath)
	if err != nil {
		log.Fatalf("Failed to load TLS certificate: %v", err)
	}

	// 创建 TLS 配置
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   cfg.TLS.NextProtos,
		MinVersion:   tls.VersionTLS13,
		MaxVersion:   tls.VersionTLS13,
	}

	// 创建 QUIC 配置
	quicConfig := &quic.Config{
		MaxIdleTimeout: config.ParseDuration(cfg.QUIC.MaxIdleTimeout),
		Allow0RTT:      cfg.QUIC.Enable0RTT,
	}

	// 启动 QUIC 监听器
	listener, err := quic.ListenAddr(cfg.Server.ListenAddr, tlsConfig, quicConfig)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer listener.Close()

	log.Printf("Listening on %s", cfg.Server.ListenAddr)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			conn, err := listener.Accept(ctx)
			if err != nil {
				log.Printf("Accept error: %v", err)
				continue
			}

			log.Printf("New connection accepted")

			go func() {
				for {
					stream, err := conn.AcceptStream(ctx)
					if err != nil {
						log.Printf("Accept stream error: %v", err)
						return
					}

					go handleStream(stream)
				}
			}()
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutting down...")
}

func handleStream(stream *quic.Stream) {
	defer stream.Close()

	// 设置读取超时
	stream.SetReadDeadline(time.Now().Add(10 * time.Second))

	buf := make([]byte, 1024)
	n, err := stream.Read(buf)
	if err != nil {
		log.Printf("Read error: %v", err)
		return
	}

	// 解析 MASQUE CONNECT 请求
	request := string(buf[:n])
	lines := strings.Split(request, "\r\n")
	if len(lines) == 0 {
		log.Printf("Invalid MASQUE request")
		stream.Write([]byte("HTTP/3 400 Bad Request\r\n\r\n"))
		return
	}

	// 提取目标地址
	parts := strings.Split(lines[0], " ")
	if len(parts) < 2 {
		log.Printf("Invalid MASQUE request format")
		stream.Write([]byte("HTTP/3 400 Bad Request\r\n\r\n"))
		return
	}

	target := parts[1]
	log.Printf("MASQUE connect request to: %s", target)

	// 建立到目标的连接（设置超时）
	conn, err := net.DialTimeout("tcp", target, 5*time.Second)
	if err != nil {
		log.Printf("Failed to connect to %s: %v", target, err)
		stream.Write([]byte("HTTP/3 502 Bad Gateway\r\n\r\n"))
		return
	}
	defer conn.Close()

	// 发送成功响应
	_, err = stream.Write([]byte("HTTP/3 200 OK\r\n\r\n"))
	if err != nil {
		log.Printf("Failed to send response: %v", err)
		return
	}

	log.Printf("Connected to %s", target)

	// 双向数据转发
	go func() {
		_, err := io.Copy(conn, stream)
		if err != nil && err != io.EOF {
			log.Printf("Error copying from client to target: %v", err)
		}
	}()

	_, err = io.Copy(stream, conn)
	if err != nil && err != io.EOF {
		log.Printf("Error copying from target to client: %v", err)
	}
}
