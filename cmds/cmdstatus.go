package cmds

import (
	"github.com/simplemoon/deploy/conf"
	"github.com/simplemoon/deploy/report"
)

// 状态查询
type CmdStatus struct {
	*CommandBase
}

func NewCmdStatus() *CmdStatus {
	return &CmdStatus{
		CommandBase: NewCommandBase(CmdNameStatus),
	}
}

func (cmd *CmdStatus) Copy() ICommand {
	return NewCmdStatus()
}

// 运行
func (cmd *CmdStatus) Run(ds *conf.DataSet) (ret []report.RowInfo, err error) {
	ret = cmd.result
	pl := conf.GetProcInfo(ds, conf.ProcTypeContainsAll)
	idx := ds.GetIndex()
	if idx <= 0 {
		return ret, ErrNotFoundServerIndex
	}
	for _, pi := range pl {
		// 检查服务器是否已经开启了
		s, err := GetServiceCmdResult(pi.GetName(), idx, conf.EchoCmdStatus)
		if err != nil {
			return ret, err
		}
		r := report.RowInfo{
			ServerIdx:  idx,
			Action:     cmd.GetName(),
			ActionType: pi.GetBaseName(),
			State:      s,
		}
		cmd.AddResult(&r)
	}
	return
}

func init() {
	Register(NewCmdStatus())
}
