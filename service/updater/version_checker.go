package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/kkkkninezzz/maa-resource-updater/internal/http"
)

/*

{
    "activity": {
        "name": "泰拉饭",
        "time": 1725249600
    },
    "gacha": {
        "pool": "泰拉饭，呜呼，泰拉饭",
        "time": 1725249600
    },
    "last_updated": "2024-09-02 15:41:33.483"
}

*/

// 定义时间格式字符串
var timeFormat = "2006-01-02 15:04:05.000"

type ResourceVersionFileContent struct {
	Activity struct {
		Name string `json:"name"`
		Time int64  `json:"time"`
	} `json:"activity"`
	LastUpdated string `json:"last_updated"`
}

func ProxyTest() {
	client := http.Client()

	resp, err := client.Get("https://www.google.com/")
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Body)
}

// 检查资源是否有更新
func CheckResourceVersion(remoteVersionFileUrl string, localVersionFilePath string) bool {
	remoteLastUpdatedTime := getRemoteVerionLastUpdated(remoteVersionFileUrl)
	localLastUpdatedTime := getLocalVerionLastUpdated(localVersionFilePath)

	// 如果远端的最后更新时间晚于本地的更新时间，那么视为有更新
	return remoteLastUpdatedTime.After(localLastUpdatedTime)
}

func getRemoteVerionLastUpdated(versionFileUrl string) time.Time {
	client := http.Client()

	resp, err := client.Get(versionFileUrl)
	if err != nil {
		log.Panicf("Error download remote version file: %v", err)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Panicf("Error read remote version file: %v", err)
	}

	return getVersionLastUpdated(data)
}

func getVersionLastUpdated(data []byte) time.Time {
	var resourceVersion ResourceVersionFileContent
	err := json.Unmarshal(data, &resourceVersion)
	if err != nil {
		log.Panicf("Error parse version file: %v", err)
	}

	lastUpdatedTime, err := time.Parse(timeFormat, resourceVersion.LastUpdated)
	if err != nil {
		log.Panicf("Error parse version lastUpdatedTime: %v", err)
	}

	return lastUpdatedTime
}

func getLocalVerionLastUpdated(localVersionFilePath string) time.Time {
	data, err := os.ReadFile(localVersionFilePath)
	if err != nil {
		log.Panicf("Error read local version file: %v", err)
	}

	return getVersionLastUpdated(data)
}
