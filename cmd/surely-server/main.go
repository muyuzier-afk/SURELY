package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"flag"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"surely/internal/config"
	"github.com/quic-go/quic-go/http3"
)

func main() {
	configPath := flag.String("config", "config.toml", "Path to configuration file")
	createSSL := flag.Bool("creat-ssl", false, "Create TLS certificate and key files")
	flag.Parse()

	// 处理 --creat-ssl 选项
	if *createSSL {
		log.Println("Creating TLS certificate and key files...")
		if err := createTLSCertificates(); err != nil {
			log.Fatalf("Failed to create TLS certificates: %v", err)
		}
		log.Println("TLS certificate and key files created successfully!")
		log.Println("Files created:")
		log.Println("- server.key (private key)")
		log.Println("- server.crt (certificate)")
		return
	}

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
	log.Printf("Loaded TLS certificate from %s and %s", cfg.Server.TLSCertPath, cfg.Server.TLSKeyPath)

	// 配置 TLS
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

	// 创建 HTTP/3 服务器
	server := &http3.Server{
		Addr:      cfg.Server.ListenAddr,
		TLSConfig: tlsConfig,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Received request: %s %s", r.Method, r.URL.Path)
			w.Write([]byte("Hello from Surely Server v1.0!"))
		}),
	}

	// 启动服务器
	go func() {
		log.Printf("Starting HTTP/3 server on %s...", cfg.Server.ListenAddr)
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutting down...")
}

// createTLSCertificates 创建 TLS 证书和密钥文件
func createTLSCertificates() error {
	// 生成私钥
	privKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}

	// 创建证书模板
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Surely Proxy"},
			CommonName:   "localhost",
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
		DNSNames:     []string{"localhost"},
	}

	// 生成证书
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privKey.PublicKey, privKey)
	if err != nil {
		return err
	}

	// 保存证书
	if err := os.WriteFile("server.crt", certDER, 0644); err != nil {
		return err
	}

	// 保存私钥
	privKeyDER := x509.MarshalPKCS1PrivateKey(privKey)
	if err := os.WriteFile("server.key", privKeyDER, 0600); err != nil {
		return err
	}

	return nil
}

