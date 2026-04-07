package main

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/quic-go/quic-go/http3"
)

func main() {
	// 加载 TLS 证书
	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Fatalf("Failed to load TLS certificate: %v", err)
	}

	// 配置 TLS
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"h3", "http/1.1"},
	}

	// 创建 HTTP 处理函数
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		w.Write([]byte("Hello from server!"))
	})

	// 创建 HTTP/3 服务器
	http3Server := &http3.Server{
		Addr:      ":8443",
		TLSConfig: tlsConfig,
		Handler:   handler,
	}

	// 创建 HTTP/1.1 服务器
	http1Server := &http.Server{
		Addr:      ":8081",
		Handler:   handler,
	}

	// 启动 HTTP/3 服务器
	go func() {
		log.Println("Starting HTTP/3 server on :8443...")
		if err := http3Server.ListenAndServe(); err != nil {
			log.Fatalf("Failed to start HTTP/3 server: %v", err)
		}
	}()

	// 启动 HTTP/1.1 服务器
	go func() {
		log.Println("Starting HTTP/1.1 server on :8081...")
		if err := http1Server.ListenAndServe(); err != nil {
			log.Fatalf("Failed to start HTTP/1.1 server: %v", err)
		}
	}()

	// 等待
	select {}
}
