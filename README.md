# maa-resource-updater
maa的资源更新工具，通过比较[MaaResource](https://github.com/MaaAssistantArknights/MaaResource)仓库中的[resource/version.json](https://github.com/MaaAssistantArknights/MaaResource/blob/main/resource/version.json)中最新的更新时间，与本地maa的资源更新时间。如果存在更新，那么将下载MaaResource的main分支压缩包，并解压覆盖到maa本地目录中。

## 使用方式
1. `maa-resource-updater.exe`放置于maa的根目录
2. 双击启动`maa-resource-updater.exe`