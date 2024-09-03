package main

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

func main() {
	// updater.CheckResourceVersion(config.MaaResourceVersionFileUrl)
	//updater.ProxyTest()
	key, err := registry.OpenKey(
		registry.CURRENT_USER,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Internet Settings`,
		registry.QUERY_VALUE)
	if err != nil {
		panic(err)
	}

	proxyEnable, _, err := key.GetIntegerValue("ProxyEnable")
	if err != nil {
		panic(err)
	}
	fmt.Println(proxyEnable)

	proxyServer, _, err := key.GetStringValue("ProxyServer")
	if err != nil {
		panic(err)
	}
	fmt.Println(proxyServer)
}
