package http

import (
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/sys/windows/registry"
)

func Client() *http.Client {
	if !proxyEnable() {
		return http.DefaultClient
	}

	proxyServer := getProxyServer()
	if proxyServer == "" {
		return http.DefaultClient
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				return url.Parse(proxyServer)
			},
		},
	}
	return client

}

func proxyEnable() bool {
	key, err := registry.OpenKey(
		registry.CURRENT_USER,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Internet Settings`,
		registry.QUERY_VALUE)
	if err != nil {
		return false
	}

	proxyEnable, _, err := key.GetIntegerValue("ProxyEnable")
	if err != nil {
		return false
	}

	return proxyEnable == 1
}

func getProxyServer() string {
	key, err := registry.OpenKey(
		registry.CURRENT_USER,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Internet Settings`,
		registry.QUERY_VALUE)

	if err != nil {
		return ""
	}

	proxyServer, _, err := key.GetStringValue("ProxyServer")
	if err != nil {
		return ""
	}

	// 说明是基于http的代理，但是windows在http代理中不会添加http协议头
	if !strings.HasPrefix(proxyServer, "http") && !strings.HasPrefix(proxyServer, "socks") {
		proxyServer = "http://" + proxyServer
	}
	return proxyServer
}
