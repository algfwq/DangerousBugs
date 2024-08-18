package main

import (
	"net"
	"net/http"
	"strings"
)

func getClientIP(r *http.Request) string {
	// 尝试从X-Forwarded-For头获取IP地址
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		ips := strings.Split(ip, ", ")
		if len(ips) > 0 {
			return ips[0]
		}
	}

	// 尝试从X-Real-IP头获取IP地址
	ip = r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	// 从RemoteAddr获取IP地址
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
