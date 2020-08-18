package cmds

import (
	"github.com/simplemoon/deploy/conf"
	"github.com/simplemoon/deploy/report"
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
	if err := cmd.CommandBase.Check(ds); err != nil {
		return err
	}
	// 检查服务器是否已经开启了

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
