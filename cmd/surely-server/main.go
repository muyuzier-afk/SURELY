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

	cert, err := tls.LoadX509KeyPair(cfg.Server.TLSCertPath, cfg.Server.TLSKeyPath)
	if err != nil {
		log.Fatalf("Failed to load TLS certificate: %v", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   cfg.TLS.NextProtos,
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
	}

	quicConfig := &quic.Config{
		MaxIdleTimeout: config.ParseDuration(cfg.QUIC.MaxIdleTimeout),
		Allow0RTT:      cfg.QUIC.Enable0RTT,
	}

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

func handleStream(stream quic.Stream) {
	defer stream.Close()

	buf := make([]byte, 1024)
	n, err := stream.Read(buf)
	if err != nil {
		log.Printf("Read error: %v", err)
		return
	}

	target := string(buf[:n])
	log.Printf("Connect request to: %s", target)

	conn, err := net.Dial("tcp", target)
	if err != nil {
		log.Printf("Failed to connect to %s: %v", target, err)
		return
	}
	defer conn.Close()

	_, err = stream.Write([]byte("OK\n"))
	if err != nil {
		log.Printf("Failed to send OK: %v", err)
		return
	}

	log.Printf("Connected to %s", target)

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
