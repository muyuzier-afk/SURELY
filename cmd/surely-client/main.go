package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"surely/internal/config"
	"surely/internal/transport"
	"surely/pkg/socks5"
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

	tlsConfig := transport.TLSConfigForClient(cfg.Client.InsecureSkipVerify, cfg.TLS.NextProtos)
	quicConfig := transport.QUICConfig(
		config.ParseDuration(cfg.QUIC.MaxIdleTimeout),
		cfg.QUIC.Enable0RTT,
	)

	transportConfig := &transport.Config{
		TLSConfig:  tlsConfig,
		QUICConfig: quicConfig,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	transClient, err := transport.NewClient(ctx, cfg.Client.ServerAddr, transportConfig)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer transClient.Close()

	log.Printf("Connected to server: %s", cfg.Client.ServerAddr)

	socksServer := socks5.NewServer(*socksAddr, func(conn net.Conn, target string) error {
		log.Printf("Proxying connection to: %s", target)

		stream, err := transClient.OpenStream()
		if err != nil {
			log.Printf("Failed to open stream: %v", err)
			return err
		}
		defer stream.Close()

		_, err = stream.Write([]byte("CONNECT " + target + "\n"))
		if err != nil {
			log.Printf("Failed to send connect: %v", err)
			return err
		}

		buf := make([]byte, 1024)
		n, err := stream.Read(buf)
		if err != nil {
			log.Printf("Failed to read response: %v", err)
			return err
		}

		if string(buf[:n]) != "OK\n" {
			log.Printf("Unexpected response: %s", string(buf[:n]))
			return io.ErrUnexpectedEOF
		}

		log.Printf("Connected to %s via proxy", target)

		go func() {
			_, err := io.Copy(stream, conn)
			if err != nil && err != io.EOF {
				log.Printf("Error copying from client to server: %v", err)
			}
		}()

		_, err = io.Copy(conn, stream)
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
