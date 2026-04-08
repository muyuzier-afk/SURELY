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
	"syscall"
	"time"

	"surely/internal/config"
	"surely/pkg/socks5"

	"github.com/quic-go/quic-go"
)

func main() {
	configPath := flag.String("config", "client-config.toml", "Path to client configuration file")
	socksAddr := flag.String("socks", ":1080", "SOCKS5 proxy listen address")
	flag.Parse()

	log.Println("Surely Client v1.0.1 starting...")

	cfg, err := config.LoadClientConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Printf("Loaded configuration from %s", *configPath)

	// 创建 TLS 配置
	tlsConfig := &tls.Config{
		InsecureSkipVerify: cfg.Client.InsecureSkipVerify,
		NextProtos:         cfg.TLS.NextProtos,
		MinVersion:         tls.VersionTLS13,
		MaxVersion:         tls.VersionTLS13,
	}

	// 创建 QUIC 配置
	quicConfig := &quic.Config{
		MaxIdleTimeout: config.ParseDuration(cfg.QUIC.MaxIdleTimeout),
		Allow0RTT:      cfg.QUIC.Enable0RTT,
	}

	// 建立 QUIC 连接
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := quic.DialAddr(ctx, cfg.Client.ServerAddr, tlsConfig, quicConfig)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.CloseWithError(0, "")

	log.Printf("Connected to server: %s", cfg.Client.ServerAddr)

	socksServer := socks5.NewServer(*socksAddr, func(clientConn net.Conn, target string) error {
		log.Printf("Proxying connection to: %s", target)

		// 创建 QUIC 流
		stream, err := conn.OpenStreamSync(context.Background())
		if err != nil {
			log.Printf("Failed to open stream: %v", err)
			return err
		}
		defer stream.Close()

		// 发送 MASQUE CONNECT 请求
		_, err = stream.Write([]byte("CONNECT " + target + " HTTP/3\r\n\r\n"))
		if err != nil {
			log.Printf("Failed to send MASQUE request: %v", err)
			return err
		}

		// 读取 MASQUE 响应
		buf := make([]byte, 1024)
		n, err := stream.Read(buf)
		if err != nil {
			log.Printf("Failed to read MASQUE response: %v", err)
			return err
		}

		// 验证响应
		response := string(buf[:n])
		if len(response) < 13 || response[:12] != "HTTP/3 200 " {
			log.Printf("Unexpected MASQUE response: %s", response)
			return io.ErrUnexpectedEOF
		}

		log.Printf("Connected to %s via MASQUE proxy", target)

		// 双向数据转发
		go func() {
			_, err := io.Copy(stream, clientConn)
			if err != nil && err != io.EOF {
				log.Printf("Error copying from client to server: %v", err)
			}
		}()

		_, err = io.Copy(clientConn, stream)
		if err != nil && err != io.EOF {
			log.Printf("Error copying from server to client: %v", err)
		}

		return nil
	})

	go func() {
		if err := socksServer.ListenAndServe(); err != nil {
			log.Fatalf("SOCKS5 server error: %v", err)
		}
	}()

	log.Printf("SOCKS5 proxy listening on %s", *socksAddr)
	log.Println("Client is ready. Configure your browser/app to use SOCKS5 proxy at", *socksAddr)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutting down...")
}
