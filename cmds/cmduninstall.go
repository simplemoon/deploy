package cmds

import (
	"github.com/simplemoon/deploy/conf"
	"github.com/simplemoon/deploy/report"
	"sync"
)

// 卸载服务
type CmdUninstall struct {
	*CommandBase
}

func NewCmdUninstall() *CmdUninstall {
	return &CmdUninstall{
		CommandBase: NewCommandBase(CmdNameUninstall),
	}
}

func (cmd *CmdUninstall) Copy() ICommand {
	return NewCmdUninstall()
}

// 运行
func (cmd *CmdUninstall) Run(ds *conf.DataSet) (ret []report.RowInfo, err error) {
	ret = cmd.result
	cmd.SetStep("running")
	var wg sync.WaitGroup

	// 执行对应的结果
	err = cmd.ExecServiceCmd(ds, &wg, ServiceCmdNameUnInstall, conf.ProcTypeContainsAll)
	if err != nil {
		return
	}
	return
}

func init() {
	Register(NewCmdUninstall())
}
