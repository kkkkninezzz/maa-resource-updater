package config

type MaaResourceUpdaterConfig struct {
	ResourceVersionFileUrl string `json:"resourceVersionFileUrl"`
	ResourceFileUrl        string `json:"resourceFileUrl"`
}

var MaaResourceVersionFileUrl = "https://raw.githubusercontent.com/MaaAssistantArknights/MaaResource/main/resource/version.json"
var MaaResourceFileUrl = "https://github.com/MaaAssistantArknights/MaaResource/archive/refs/heads/main.zip"
