package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

func main() {
	// 加载 TLS 证书
	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Fatalf("Failed to load TLS certificate: %v", err)
	}

	// 配置 TLS，模拟 Chrome 141 指纹
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"h3"}, // 仅兼容 HTTP/3
		// 模拟 Chrome 141 TLS 指纹
		CipherSuites: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		},
		CurvePreferences: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
			tls.CurveP384,
		},
		MinVersion: tls.VersionTLS13,
		MaxVersion: tls.VersionTLS13,
	}

	// 配置 QUIC，使用 BBR 算法（模拟 Chrome 浏览器）
	quicConfig := &quic.Config{
		// 其他配置
		MaxIdleTimeout: 30 * time.Second,
		Allow0RTT:      true,
	}

	// 创建 HTTP 处理函数
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		w.Write([]byte("Hello from server!"))
	})

	// 创建 HTTP/3 服务器
	http3Server := &http3.Server{
		Addr:        ":8443",
		TLSConfig:   tlsConfig,
		QUICConfig:  quicConfig,
		Handler:     handler,
	}

	// 启动 HTTP/3 服务器
	go func() {
		log.Println("Starting HTTP/3 server on :8443...")
		if err := http3Server.ListenAndServe(); err != nil {
			log.Fatalf("Failed to start HTTP/3 server: %v", err)
		}
	}()

	// 等待
	select {}
}
