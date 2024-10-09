package utils

import (
	"crypto/rand"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"shodan-proxy/internal/config"
)

// 删除这行，因为我们现在使用 config 包中的 ConfigMutex
// var configMutex sync.RWMutex

func GenerateSecureToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}

func IsValidSession(sessionID string) bool {
	// 这里应该实现会话验证逻辑
	// 为了简单起见，我们总是返回 true
	return true
}

func GetClientIP(r *http.Request) string {
	config.ConfigMutex.RLock()
	defer config.ConfigMutex.RUnlock()

	// 首先获取直接连接的 IP
	remoteIP, _, _ := net.SplitHostPort(r.RemoteAddr)

	// 检查 X-Forwarded-For 头
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		// 从右到左遍历 IP 列表
		for i := len(ips) - 1; i >= 0; i-- {
			ip := strings.TrimSpace(ips[i])
			// 如果这个 IP 不是受信任的代理，就返回它
			if !isIPTrusted(ip) {
				return ip
			}
		}
	}

	// 检查 X-Real-IP 头
	if xrip := r.Header.Get("X-Real-IP"); xrip != "" && isIPTrusted(remoteIP) {
		return xrip
	}

	// 如果没有找到可信的客户端 IP，返回直接连接的 IP
	return remoteIP
}

func isIPTrusted(ip string) bool {
	for _, trustedProxy := range config.GlobalConfig.TrustedProxies {
		_, ipNet, err := net.ParseCIDR(trustedProxy)
		if err == nil && ipNet.Contains(net.ParseIP(ip)) {
			return true
		}
	}
	return false
}

func IsPathBlocked(path string) bool {
	config.ConfigMutex.RLock()
	defer config.ConfigMutex.RUnlock()

	for _, blockedPath := range config.GlobalConfig.BlockedPaths {
		if strings.HasPrefix(path, blockedPath) {
			return true
		}
	}
	return false
}

func GetNextKey() string {
	config.ShodanKeysMutex.RLock()
	defer config.ShodanKeysMutex.RUnlock()

	if len(config.ShodanKeys) == 0 {
		log.Println("Warning: No Shodan API keys available. Please add keys in the admin panel.")
		return ""
	}

	key := config.ShodanKeys[config.CurrentKey]
	config.CurrentKey = (config.CurrentKey + 1) % len(config.ShodanKeys)
	return key
}

func IsIPTrusted(ip string) bool {
	config.ConfigMutex.RLock()
	defer config.ConfigMutex.RUnlock()

	return isIPTrusted(ip)
}

func IsIPAllowed(ip string) bool {
	config.ConfigMutex.RLock()
	defer config.ConfigMutex.RUnlock()

	// 如果白名单为空，允许所有 IP
	if len(config.GlobalConfig.AllowedIPs) == 0 {
		return true
	}

	for _, allowedIP := range config.GlobalConfig.AllowedIPs {
		if allowedIP == ip {
			return true
		}
		
		// 检查 CIDR
		_, ipNet, err := net.ParseCIDR(allowedIP)
		if err == nil {
			if ipNet.Contains(net.ParseIP(ip)) {
				return true
			}
		}
	}

	return false
}
