package controller

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// PingHost ping指定主机并返回是否可达以及往返时间
func PingHost(host string) (bool, time.Duration, error) {
	// 首先尝试使用net.DialTimeout快速检查主机是否可达
	startTime := time.Now()
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:80", host), 3*time.Second)
	if err == nil {
		// 连接成功，主机可达
		conn.Close()
		return true, time.Since(startTime), nil
	}

	// 如果TCP连接失败，使用系统ping命令
	// 根据不同操作系统使用不同的ping命令参数
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("ping", "-n", "1", "-w", "3000", host)
	default: // Linux, Darwin等
		cmd = exec.Command("ping", "-c", "1", "-W", "3", host)
	}

	startTime = time.Now()
	output, err := cmd.CombinedOutput()
	pingDuration := time.Since(startTime)

	if err != nil {
		// 检查输出以确定是错误还是主机不可达
		outputStr := string(output)
		if strings.Contains(outputStr, "timed out") || 
		   strings.Contains(outputStr, "100% packet loss") ||
		   strings.Contains(outputStr, "Destination Host Unreachable") {
			return false, 0, nil // 主机不可达，但不是错误
		}
		return false, 0, fmt.Errorf("ping error: %w, output: %s", err, outputStr)
	}

	// 简单地返回执行命令的时间作为RTT
	// 可以进一步解析输出以获取更准确的RTT
	return true, pingDuration, nil
}

// IsPortOpen 检查指定主机的指定端口是否开放
func IsPortOpen(host string, port int) (bool, error) {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, 3*time.Second)
	
	if err != nil {
		// 如果是超时或连接被拒绝，端口可能是关闭的
		if strings.Contains(err.Error(), "timeout") || 
		   strings.Contains(err.Error(), "refused") {
			return false, nil
		}
		// 其他错误
		return false, err
	}
	
	// 连接成功，关闭连接并返回
	conn.Close()
	return true, nil
}
