package cmds

import "github.com/simplemoon/deploy/conf"

// 启服
type CmdStopPlatform struct {
	*CmdStop
}

// 创建命令
func NewCmdStopPlatform() *CmdStopServers {
	cmd := &CmdStop{
		CommandBase: NewCommandBase(CmdNameCloseGame),
		Contains:    conf.ProcTypeContainsPlatform,
	}

	return &CmdStopServers{
		CmdStop: cmd,
	}
}

// 拷贝对应的命令
func (cmd *CmdStopPlatform) Copy() ICommand {
	return NewCmdStopPlatform()
}

func init() {
	Register(NewCmdStopPlatform())
}
