package utils

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

// GetHostname 获取主机名
func GetHostname() (string, error) {
	return os.Hostname()
}

// GetLocalIP 获取本地IP地址
func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	
	for _, addr := range addrs {
		// 检查是否是IP网络地址
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			// 只考虑IPv4地址
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	
	return "", fmt.Errorf("no IP address found")
}

// WaitForConnection 等待网络连接可用
func WaitForConnection(host string, port int, timeout time.Duration) bool {
	// 计算超时时间
	deadline := time.Now().Add(timeout)
	
	// 循环尝试连接，直到超时
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), time.Second)
		if err == nil {
			conn.Close()
			return true
		}
		
		// 等待一秒再试
		time.Sleep(time.Second)
	}
	
	return false
}

// FormatBytes 将字节数格式化为人类可读的形式
func FormatBytes(bytes uint64) string {
	const (
		_          = iota // 忽略0
		KB float64 = 1 << (10 * iota)
		MB
		GB
		TB
		PB
	)
	
	var size string
	var unit string
	
	switch {
	case bytes >= uint64(PB):
		size = fmt.Sprintf("%.2f", float64(bytes)/PB)
		unit = "PB"
	case bytes >= uint64(TB):
		size = fmt.Sprintf("%.2f", float64(bytes)/TB)
		unit = "TB"
	case bytes >= uint64(GB):
		size = fmt.Sprintf("%.2f", float64(bytes)/GB)
		unit = "GB"
	case bytes >= uint64(MB):
		size = fmt.Sprintf("%.2f", float64(bytes)/MB)
		unit = "MB"
	case bytes >= uint64(KB):
		size = fmt.Sprintf("%.2f", float64(bytes)/KB)
		unit = "KB"
	default:
		size = fmt.Sprintf("%d", bytes)
		unit = "bytes"
	}
	
	return strings.TrimSuffix(size, ".00") + " " + unit
}
