package conf

import "fmt"

// 配置文件名称
const (
	FileNameBase    = "base.json"
	FileNameServers = "servers.json"
)

// ini section 名称
const (
	IniKeySectionEcho = "Echo"
	IniKeyAddr        = "Addr"
	IniKeyPort        = "Port"
)

// ini文件的名称
const (
	PreIniFileName = "sd" // ini文件名称的前缀

	ProcessKeyGame = "game" // world配置的名称
)

// 命令的定义
const (
	EchoCmdStatus = "status"
)

func GetIniFileName(name string) string {
	return fmt.Sprintf("%s_%s.ini", PreIniFileName, name)
}
