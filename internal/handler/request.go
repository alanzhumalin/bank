package handler

import (
	"net"
	"net/http"
	"strings"
)

type RequestMetadata struct {
	Ip    string
	Agent string
}

func GetRequestMetadata(r *http.Request) *RequestMetadata {
	userAgent := GetUserAgentFromRequest(r)
	ip := GetIpFromRequest(r)

	return &RequestMetadata{
		Ip:    ip,
		Agent: userAgent,
	}
}

func GetUserAgentFromRequest(r *http.Request) string {
	return r.Header.Get("User-Agent")
}

func GetIpFromRequest(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	ip = r.Header.Get("X-Forwarded-For") //array of ips
	if ip != "" {
		return strings.TrimSpace(strings.Split(ip, ",")[0])
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)

	if err != nil {
		return r.RemoteAddr
	}

	return host
}
