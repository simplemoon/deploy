package cmds

import (
	"fmt"
	"github.com/simplemoon/deploy/report"
	"strings"

	"github.com/simplemoon/deploy/conf"
)

// 状态查询
type CmdCState struct {
	*CommandBase
}

func NewCmdCState() *CmdCState {
	return &CmdCState{
		CommandBase: NewCommandBase(CmdNameCState),
	}
}

func (cmd *CmdCState) Copy() ICommand {
	return NewCmdCState()
}

// 运行
func (cmd *CmdCState) Run(ds *conf.DataSet) (ret []report.RowInfo, err error) {
	ret = cmd.result
	// 时间
	idx := ds.GetIndex()
	if idx <= 0 {
		err = fmt.Errorf("don't found server config")
		return
	}

	var args string
	tn := conf.GetTargetName()
	if tn == "" {
		args, err = ds.GetString(conf.InfoKeyToolArgs)
		if err != nil {
			return
		}
		if args == "" {
			err = fmt.Errorf("telnet param not found")
			return
		}
		pos := strings.Index(args, " ")
		if pos <= 0 {
			err = fmt.Errorf("telnet param little than 2")
			return
		}
		tn = args[:pos]
		args = args[pos+1:]
	} else {
		args = conf.GetTelnetArgs()
	}
	if tn != conf.ProcessNameWorld && tn != conf.ProcessNamePlatform || args == "" {
		err = fmt.Errorf("process or param is nil")
		return
	}
	// 进程名称
	var procName string
	if tn == conf.ProcessNameWorld {
		procName = conf.ProcessKeyGame
	} else {
		procName = conf.ProcessKeyPlatform
	}
	// 参数
	args = fmt.Sprintf("%s %s", tn, args)
	// 具体的命令
	s, err := GetServiceCmdResult(procName, idx, args)
	if err != nil {
		cmd.GetLogger().Printf("query %s %s -> %v\n", procName, args, err)
		return
	}
	cmd.result = append(cmd.result, report.RowInfo{
		ServerIdx:  idx,
		Action:     cmd.Name,
		ActionType: "param",
		State:      report.Success,
		Msg:        args,
	}, report.RowInfo{
		ServerIdx:  idx,
		Action:     cmd.Name,
		ActionType: "result",
		State:      report.Success,
		Msg:        s,
	})
	return
}

func init() {
	Register(NewCmdCState())
}
