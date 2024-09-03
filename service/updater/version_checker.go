package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
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

type remoteResourceVersion struct {
	Activity struct {
		Name string `json:"name"`
		Time int64  `json:"time"`
	} `json:"activity"`
	LastUpdated string `json:"last_updated"`
}

func ProxyTest() {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment, // 使用系统代理
		},
	}

	resp, err := client.Get("https://www.google.com/")
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Body)
}

// 检查资源是否有更新
func CheckResourceVersion(versionFileUrl string) bool {
	remoteLastUpdatedTime := getRemoteVerionLastUpdated(versionFileUrl)
	fmt.Println(remoteLastUpdatedTime)
	return true
}

func getRemoteVerionLastUpdated(versionFileUrl string) time.Time {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment, // 使用系统代理
		},
	}

	resp, err := client.Get(versionFileUrl)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var remoteResourceVersion remoteResourceVersion
	err = json.Unmarshal(data, &remoteResourceVersion)
	if err != nil {
		panic(err)
	}

	lastUpdatedTime, err := time.Parse(timeFormat, remoteResourceVersion.LastUpdated)
	if err != nil {
		panic(err)
	}

	return lastUpdatedTime
}
