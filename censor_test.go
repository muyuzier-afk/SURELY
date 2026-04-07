package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

// Censor 审查者结构体
type Censor struct {
	interfaceName string
	filter        string
	handle        *pcap.Handle
	// 审查规则
	rules []Rule
}

// Rule 审查规则接口
type Rule interface {
	Name() string
	Check(packet gopacket.Packet) (bool, string)
}

// QUICRule QUIC 流量审查规则
type QUICRule struct {}

func (r *QUICRule) Name() string {
	return "QUIC 流量审查"
}

func (r *QUICRule) Check(packet gopacket.Packet) (bool, string) {
	// 检查是否是 QUIC 数据包
	udpLayer := packet.Layer(layers.LayerTypeUDP)
	if udpLayer == nil {
		return false, ""
	}

	udp, _ := udpLayer.(*layers.UDP)
	// 检查 UDP 负载长度
	if len(udp.Payload) < 1 {
		return false, ""
	}

	// 检查是否是 QUIC 长包头（前 2 位为 1）
	firstByte := udp.Payload[0]
	if (firstByte & 0xC0) != 0xC0 {
		return false, ""
	}

	// 检查版本号是否为 0x00000001
	if len(udp.Payload) >= 5 {
		version := udp.Payload[1:5]
		if !bytes.Equal(version, []byte{0x00, 0x00, 0x00, 0x01}) {
			return true, "QUIC 版本号不符合要求"
		}
	}

	// 检查连接 ID 是否符合可预测模式
	// 这里简化实现，实际应该根据具体的可预测模式进行检查

	return false, ""
}

// TLSRule TLS 流量审查规则
type TLSRule struct {}

func (r *TLSRule) Name() string {
	return "TLS 流量审查"
}

func (r *TLSRule) Check(packet gopacket.Packet) (bool, string) {
	// 检查是否是 TLS 数据包
	udpLayer := packet.Layer(layers.LayerTypeUDP)
	if udpLayer == nil {
		return false, ""
	}

	udp, _ := udpLayer.(*layers.UDP)
	// 检查 UDP 负载长度
	if len(udp.Payload) < 1 {
		return false, ""
	}

	// 检查是否是 QUIC 长包头
	firstByte := udp.Payload[0]
	if (firstByte & 0xC0) != 0xC0 {
		return false, ""
	}

	// 提取 QUIC 有效负载（跳过 QUIC 头部）
	// 这里简化实现，实际应该根据 QUIC 头部长度进行提取
	if len(udp.Payload) < 10 {
		return false, ""
	}

	// 检查是否包含 TLS ClientHello
	// 这里简化实现，实际应该根据 TLS 协议规范进行检查

	return false, ""
}

// PacketLengthRule 数据包长度审查规则
type PacketLengthRule struct {
	expectedDistribution []int
}

func (r *PacketLengthRule) Name() string {
	return "数据包长度审查"
}

func (r *PacketLengthRule) Check(packet gopacket.Packet) (bool, string) {
	// 检查数据包长度是否符合预期分布
	udpLayer := packet.Layer(layers.LayerTypeUDP)
	if udpLayer == nil {
		return false, ""
	}

	udp, _ := udpLayer.(*layers.UDP)
	length := len(udp.Payload)

	// 检查长度是否在预期范围内
	expected := false
	for _, expectedLen := range r.expectedDistribution {
		if length == expectedLen {
			expected = true
			break
		}
	}

	if !expected {
		return true, fmt.Sprintf("数据包长度不符合预期分布: %d", length)
	}

	return false, ""
}

// NewCensor 创建一个新的审查者
func NewCensor(interfaceName string) (*Censor, error) {
	handle, err := pcap.OpenLive(interfaceName, 65536, true, pcap.BlockForever)
	if err != nil {
		return nil, err
	}

	// 设置过滤器，只捕获 UDP 流量
	filter := "udp"
	err = handle.SetBPFFilter(filter)
	if err != nil {
		return nil, err
	}

	// 初始化审查规则
	rules := []Rule{
		&QUICRule{},
		&TLSRule{},
		&PacketLengthRule{
			expectedDistribution: []int{64, 128, 256, 512, 1024, 1500},
		},
		&EntropyRule{
			maxEntropy: 4.0, // 香农熵阈值，过高表示随机性太强
		},
		&ByteDistributionRule{},
		&TimingRule{},
	}

	return &Censor{
		interfaceName: interfaceName,
		filter:        filter,
		handle:        handle,
		rules:         rules,
	}, nil
}

// Start 开始审查
func (c *Censor) Start() {
	log.Println("审查者开始工作...")

	packetSource := gopacket.NewPacketSource(c.handle, c.handle.LinkType())
	for packet := range packetSource.Packets() {
		c.checkPacket(packet)
	}
}

// checkPacket 检查数据包
func (c *Censor) checkPacket(packet gopacket.Packet) {
	for _, rule := range c.rules {
		isViolation, reason := rule.Check(packet)
		if isViolation {
			c.reportViolation(rule.Name(), reason, packet)
			break
		}
	}
}

// reportViolation 报告违规
func (c *Censor) reportViolation(ruleName, reason string, packet gopacket.Packet) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] 违规 detected by %s: %s\n", timestamp, ruleName, reason)

	// 打印数据包信息
	if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
		udp, _ := udpLayer.(*layers.UDP)
		fmt.Printf("  Source: %d -> Destination: %d\n",
			udp.SrcPort, udp.DstPort)
		fmt.Printf("  Payload length: %d\n", len(udp.Payload))
		fmt.Printf("  Payload: %s\n", hex.EncodeToString(udp.Payload[:min(32, len(udp.Payload))]))
	}

	// 回炉重造：这里可以添加处理逻辑，如断开连接、发送警告等
	fmt.Println("  处理：回炉重造")
	fmt.Println()
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// CalculateShannonEntropy 计算香农熵
func CalculateShannonEntropy(data []byte) float64 {
	if len(data) == 0 {
		return 0
	}

	// 计算字节频率
	freq := make(map[byte]float64)
	for _, b := range data {
		freq[b]++
	}

	// 计算香农熵
	entropy := 0.0
	for _, count := range freq {
		p := count / float64(len(data))
		entropy -= p * math.Log2(p)
	}

	return entropy
}

// EntropyRule 香农熵审查规则
type EntropyRule struct {
	maxEntropy float64
}

func (r *EntropyRule) Name() string {
	return "香农熵审查"
}

func (r *EntropyRule) Check(packet gopacket.Packet) (bool, string) {
	// 检查数据包的香农熵
	udpLayer := packet.Layer(layers.LayerTypeUDP)
	if udpLayer == nil {
		return false, ""
	}

	udp, _ := udpLayer.(*layers.UDP)
	entropy := CalculateShannonEntropy(udp.Payload)

	if entropy > r.maxEntropy {
		return true, fmt.Sprintf("香农熵过高: %.2f", entropy)
	}

	return false, ""
}

// ByteDistributionRule 字节分布审查规则
type ByteDistributionRule struct {}

func (r *ByteDistributionRule) Name() string {
	return "字节分布审查"
}

func (r *ByteDistributionRule) Check(packet gopacket.Packet) (bool, string) {
	// 检查字节分布是否符合预期
	udpLayer := packet.Layer(layers.LayerTypeUDP)
	if udpLayer == nil {
		return false, ""
	}

	udp, _ := udpLayer.(*layers.UDP)
	payload := udp.Payload

	// 检查是否包含过多的随机字节
	// 这里简化实现，实际应该根据具体的预期分布进行检查
	if len(payload) > 0 {
		// 检查前 8 字节是否为时间戳（可预测）
		if len(payload) >= 8 {
			timestampBytes := payload[:8]
			// 检查是否为有效的时间戳范围
			// 这里简化实现，实际应该检查时间戳是否在合理范围内
		}
	}

	return false, ""
}

// TimingRule 时序审查规则
type TimingRule struct {
	lastPacketTime time.Time
}

func (r *TimingRule) Name() string {
	return "时序审查"
}

func (r *TimingRule) Check(packet gopacket.Packet) (bool, string) {
	// 检查数据包的到达时间间隔
	currentTime := time.Now()

	if !r.lastPacketTime.IsZero() {
		interval := currentTime.Sub(r.lastPacketTime)
		// 检查是否存在完美的周期性
		if interval == 100*time.Millisecond {
			return true, "数据包到达间隔过于规律"
		}
	}

	r.lastPacketTime = currentTime
	return false, ""
}

// 主函数
func main() {
	// 获取网络接口
	interfaces, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatalf("Failed to find interfaces: %v", err)
	}

	// 选择第一个网络接口
	if len(interfaces) == 0 {
		log.Fatalf("No interfaces found")
	}

	interfaceName := interfaces[0].Name
	log.Printf("Using interface: %s", interfaceName)

	// 创建审查者
	censor, err := NewCensor(interfaceName)
	if err != nil {
		log.Fatalf("Failed to create censor: %v", err)
	}
	defer censor.handle.Close()

	// 启动审查
	censor.Start()
}
