package main

import (
	"crypto/tls"
	"log"
	"net/http"
)

func main() {
	// 配置 HTTP 客户端，使用 HTTP/3
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				NextProtos:         []string{"h3"},
			},
		},
	}

	// 发送请求
	resp, err := client.Get("https://localhost:443")
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body := make([]byte, 1024)
	n, err := resp.Body.Read(body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	log.Printf("Response: %s", body[:n])
}
