package app

import (
	"math/rand"
	"net/http"
	"time"
)

// UserAgents 定义常见的 User-Agent 列表
var UserAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.1 Safari/605.1.15",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:89.0) Gecko/20100101 Firefox/89.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36",
}

// FakeRequests 定义伪装请求列表
var FakeRequests = []string{
	"GET / HTTP/3",
	"GET /health HTTP/3",
	"GET /favicon.ico HTTP/3",
	"GET /robots.txt HTTP/3",
}

// ErrorCodes 定义标准 HTTP 错误码
var ErrorCodes = []int{
	404, // Not Found
	503, // Service Unavailable
	403, // Forbidden
	500, // Internal Server Error
}

// GetRandomUserAgent 获取随机 User-Agent
func GetRandomUserAgent() string {
	return UserAgents[rand.Intn(len(UserAgents))]
}

// GetRandomFakeRequest 获取随机伪装请求
func GetRandomFakeRequest() string {
	return FakeRequests[rand.Intn(len(FakeRequests))]
}

// ShouldReturnError 判断是否应该返回错误
func ShouldReturnError() bool {
	// 1% 的概率返回错误
	return rand.Float64() < 0.01
}

// GetRandomErrorCode 获取随机错误码
func GetRandomErrorCode() int {
	return ErrorCodes[rand.Intn(len(ErrorCodes))]
}

// FakeRequestMiddleware 模拟真实请求的中间件
func FakeRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 检查是否应该返回错误
		if ShouldReturnError() {
			errorCode := GetRandomErrorCode()
			http.Error(w, http.StatusText(errorCode), errorCode)
			return
		}

		// 继续处理请求
		next.ServeHTTP(w, r)
	})
}

// ScheduleFakeRequests 定期发送伪装请求
func ScheduleFakeRequests(client *http.Client) {
	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for range ticker.C {
			// 获取随机请求和 User-Agent
			GetRandomFakeRequest()
			GetRandomUserAgent()
			
			// 这里简化实现，实际应该发送 HTTP/3 请求
			// req, err := http.NewRequest("GET", "https://localhost:443", nil)
			// if err != nil {
			// 	continue
			// }
			// req.Header.Set("User-Agent", userAgent)
			// client.Do(req)
		}
	}()
}
