package cmds

import "github.com/simplemoon/deploy/conf"

// 启服
type CmdStopServers struct {
	*CmdStop
}

// 创建命令
func NewCmdStopServers() *CmdStopServers {
	cmd := &CmdStop{
		CommandBase: NewCommandBase(CmdNameCloseGame),
		Contains:    conf.ProcTypeContainsSvr,
	}

	return &CmdStopServers{
		CmdStop: cmd,
	}
}

// 拷贝对应的命令
func (cmd *CmdStopServers) Copy() ICommand {
	return NewCmdStopServers()
}

func init() {
	Register(NewCmdStopServers())
}
