package main

import (
	"flag"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"surely/internal/config"
	"surely/pkg/surely-sdk"
)

func main() {
	configPath := flag.String("config", "client-config.toml", "Path to client configuration file")
	target := flag.String("target", "", "Target path to request")
	flag.Parse()

	log.Println("Surely Client v1.0 starting...")

	// 加载配置文件
	cfg, err := config.LoadClientConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Printf("Loaded configuration from %s", *configPath)

	// 创建 SDK 客户端
	serverAddr := "https://" + cfg.Client.ServerAddr
	client, err := surely_sdk.NewClient(serverAddr, cfg.Client.InsecureSkipVerify)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	log.Printf("Connected to server: %s", serverAddr)

	// 如果指定了目标，发送请求
	if *target != "" {
		log.Printf("Sending request to: %s", *target)
		resp, err := client.Get(*target)
		if err != nil {
			log.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		log.Printf("Response status: %s", resp.Status)

		// 读取响应
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Failed to read response: %v", err)
		}

		log.Printf("Response body: %s", body)
	} else {
		log.Println("No target specified, waiting...")
	}

	// 等待信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutting down...")
}

