package cmds

import (
	"fmt"
	"github.com/simplemoon/deploy/conf"
	"github.com/simplemoon/deploy/report"
	"github.com/simplemoon/deploy/utils"
)

// 启服
type CmdStart struct {
	CommandBase                  // 基础的功能
	result      []report.RowInfo // 结果
}

// 创建命令
func NewCmdStart() CmdStart {
	return CmdStart{
		CommandBase: NewCommandBase(CmdNameStart),
		result:      make([]report.RowInfo, 0),
	}
}

// 检查
func (cmd CmdStart) Check(ds *conf.DataSet) error {
	// 基本的检查
	serverId := ds.GetServerId()
	idx, err := conf.GetIndexByServerId(serverId)
	if err != nil {
		return err
	}
	// 检查基础的是否存在
	if err := cmd.CheckByIndex(serverId, idx); err != nil {
		return err
	}
	// 检查服务器是否已经开启了
	name := conf.GetIniFileName(conf.ProcessKeyGame)
	iniPath := utils.GetServerIniConfigPath(name, idx)
	if utils.IsPathExist(iniPath) {
		return fmt.Errorf("%s not exist", iniPath)
	}
	// 加载数据
	gameEchoData, err := conf.LoadIni(iniPath)
	if err != nil {
		return err
	}
	// 获取telnet
	telnetHelper := utils.NewTelnet(gameEchoData.EchoAddr, gameEchoData.EchoPort)
	if telnetHelper == nil {
		return fmt.Errorf("create telnet client failed")
	}
	s := telnetHelper.Send(conf.EchoCmdStatus)
	if s == utils.ConnectStateOpened {
		return fmt.Errorf("server %d is started, no to start", serverId)
	}
	return nil
}

// 运行对应的命令
func (cmd CmdStart) Run(ds *conf.DataSet) ([]report.RowInfo, error) {
	return nil, nil
}

// 初始化
func init() {
	Register(NewCmdStart())
}
