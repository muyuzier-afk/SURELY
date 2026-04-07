package packet

import (
	"math/rand"
	"time"
)

// Packet 定义数据包结构
type Packet struct {
	Data []byte
}

// PacketLengthDist 定义包长分布
type PacketLengthDist struct {
	Lengths []int
	Weights []float64
}

// NewPacketLengthDist 创建一个新的包长分布
func NewPacketLengthDist() *PacketLengthDist {
	// 这里使用模拟的包长分布，实际应该从真实 pcap 中提取
	lengths := []int{64, 128, 256, 512, 1024, 1500}
	weights := []float64{0.1, 0.2, 0.3, 0.2, 0.1, 0.1}

	return &PacketLengthDist{
		Lengths: lengths,
		Weights: weights,
	}
}

// GetRandomLength 根据分布获取随机包长
func (p *PacketLengthDist) GetRandomLength() int {
	r := rand.Float64()
	cumulative := 0.0

	for i, weight := range p.Weights {
		cumulative += weight
		if r <= cumulative {
			return p.Lengths[i]
		}
	}

	return p.Lengths[0]
}

// PadPacket 对数据包进行填充
func PadPacket(packet *Packet, targetLength int) []byte {
	currentLength := len(packet.Data)
	if currentLength >= targetLength {
		return packet.Data
	}

	padding := make([]byte, targetLength-currentLength)
	// 使用低熵数据填充
	for i := range padding {
		padding[i] = byte(i % 256)
	}

	return append(packet.Data, padding...)
}

// AddNoise 添加低熵背景噪声
func AddNoise() []byte {
	// 生成低熵噪声数据，模拟 Keep-Alive 或视频黑屏流
	noise := make([]byte, 64)
	// 使用固定模式的低熵数据
	for i := range noise {
		noise[i] = byte(i % 16)
	}

	return noise
}

// GetJitter 获取随机抖动时间
func GetJitter() time.Duration {
	// 生成 5-50ms 的随机抖动
	jitter := rand.Intn(46) + 5
	return time.Duration(jitter) * time.Millisecond
}
