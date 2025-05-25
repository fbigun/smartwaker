package controller

import (
	"encoding/hex"
	"fmt"
	"net"
	"strings"
)

// WakeOnLAN 发送网络唤醒 Magic Packet
func WakeOnLAN(macAddr, ipAddr string, port int) error {
	// 如果没有指定端口，使用默认端口9
	if port <= 0 {
		port = 9
	}

	// 解析MAC地址
	mac, err := parseMACAddress(macAddr)
	if err != nil {
		return fmt.Errorf("invalid MAC address: %w", err)
	}

	// 创建Magic Packet
	mp, err := createMagicPacket(mac)
	if err != nil {
		return fmt.Errorf("failed to create magic packet: %w", err)
	}

	// 确定目标地址
	// 如果提供了IP地址，使用IP:端口作为目标
	// 否则使用广播地址255.255.255.255:端口
	var target string
	if ipAddr != "" {
		target = fmt.Sprintf("%s:%d", ipAddr, port)
	} else {
		target = fmt.Sprintf("255.255.255.255:%d", port)
	}

	// 解析目标地址
	addr, err := net.ResolveUDPAddr("udp", target)
	if err != nil {
		return fmt.Errorf("failed to resolve target address: %w", err)
	}

	// 创建UDP连接
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return fmt.Errorf("failed to create UDP connection: %w", err)
	}
	defer conn.Close()

	// 发送Magic Packet
	_, err = conn.Write(mp)
	if err != nil {
		return fmt.Errorf("failed to send magic packet: %w", err)
	}

	return nil
}

// parseMACAddress 解析MAC地址字符串为字节数组
func parseMACAddress(macAddr string) ([]byte, error) {
	// 统一MAC地址格式，移除可能的分隔符
	macAddr = strings.ReplaceAll(macAddr, ":", "")
	macAddr = strings.ReplaceAll(macAddr, "-", "")
	macAddr = strings.ReplaceAll(macAddr, ".", "")

	// 检查MAC地址长度
	if len(macAddr) != 12 {
		return nil, fmt.Errorf("invalid MAC address length: %s", macAddr)
	}

	// 将MAC地址从十六进制字符串转换为字节数组
	mac, err := hex.DecodeString(macAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode MAC address: %w", err)
	}

	return mac, nil
}

// createMagicPacket 创建网络唤醒的Magic Packet
func createMagicPacket(mac []byte) ([]byte, error) {
	if len(mac) != 6 {
		return nil, fmt.Errorf("invalid MAC address length: %d bytes", len(mac))
	}

	// Magic Packet格式：6字节的0xFF，然后是目标MAC地址重复16次
	packet := make([]byte, 102)

	// 前6字节设置为0xFF
	for i := 0; i < 6; i++ {
		packet[i] = 0xFF
	}

	// 接下来重复MAC地址16次
	for i := 0; i < 16; i++ {
		copy(packet[6+(i*6):], mac)
	}

	return packet, nil
}
