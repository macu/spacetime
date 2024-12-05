package net

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"spacetime/pkg/env"
)

// getUserIP returns the IP address of the remote connection.
func GetUserIP(r *http.Request) string {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

func IsAjax(r *http.Request) bool {
	return strings.ToLower(r.Header.Get("X-Requested-With")) == "xmlhttprequest"
}

// getLocalPort returns the server (local) port number for the given request.
func GetLocalPort(r *http.Request) (string, error) {
	a, ok := r.Context().Value(http.LocalAddrContextKey).(net.Addr)
	if !ok {
		return "", fmt.Errorf("getting local address from request context")
	}
	_, port, err := net.SplitHostPort(a.String())
	if err != nil {
		return "", fmt.Errorf("extracting port: %w", err)
	}
	return port, nil
}

// buildAbsoluteURL returns an absolute URL to the given path on the currently running server.
func BuildAbsoluteURL(r *http.Request, path string) (string, error) {
	path = strings.TrimPrefix(path, "/")

	if env.IsAppEngine() {
		return fmt.Sprintf("https://%s/%s", os.Getenv("DOMAIN"), path), nil
	}

	port, err := GetLocalPort(r)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("http://localhost:%s/%s", port, path), nil
}
