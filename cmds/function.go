package cmds

import (
	"fmt"

	"github.com/simplemoon/deploy/conf"
	"github.com/simplemoon/deploy/utils"
)

func GetServiceCmdResult(procName string, idx int, cmd string) (string, error) {
	// 检查服务器是否已经开启了
	if cmd == "" {
		cmd = conf.EchoCmdStatus
	}
	// 获取名称
	name := conf.GetIniFileName(procName)
	iniPath := utils.GetServerIniConfigPath(name, idx)
	if utils.IsPathExist(iniPath) {
		return "", fmt.Errorf("%s not exist", iniPath)
	}
	// 加载数据
	gameEchoData, err := conf.LoadIni(iniPath)
	if err != nil {
		return "", err
	}
	// 获取telnet
	telnetHelper := utils.NewTelnet(gameEchoData.EchoAddr, gameEchoData.EchoPort)
	if telnetHelper == nil {
		return "", fmt.Errorf("create telnet client failed")
	}
	s := telnetHelper.SendQueryStatus(cmd)
	return s, nil
}
