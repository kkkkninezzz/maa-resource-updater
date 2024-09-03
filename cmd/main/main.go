package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/kkkkninezzz/maa-resource-updater/internal/config"
	"github.com/kkkkninezzz/maa-resource-updater/service/updater"
)

func main() {

	startUpdate()

	// 提示用户按任意键退出
	fmt.Println("Press 'Enter' to exit...")
	var input string
	fmt.Scanln(&input)
}

func startUpdate() {
	defer func() {
		if r := recover(); r != nil {

		}
	}()

	// 获取当前可执行文件的路径
	execPath, err := os.Executable()
	if err != nil {
		log.Panicf("Error getting executable path: %v", err)

	}

	// 获取可执行文件所在的目录
	execDir := filepath.Dir(execPath)

	hasUpdateFlag := updater.CheckResourceVersion(config.MaaResourceVersionFileUrl, filepath.Join(execDir, "resource", "version.json"))

	if !hasUpdateFlag {
		log.Println("It is currently the latest version.")
		return
	}

	log.Println("Discover the new version and start downloading the pack.")

	updater.UpdateResource(config.MaaResourceFileUrl, execDir)
	log.Println("Maa resource update success!")
}
