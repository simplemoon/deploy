package cmds

import (
	"fmt"
	"github.com/simplemoon/deploy/conf"
	"github.com/simplemoon/deploy/log"
	"github.com/simplemoon/deploy/report"
)

// 命令链条
type CommandList struct {
	root ICommand // 最开始的命令
}

// 创建 CommandList
func NewCommandList(names []string) (*CommandList, error) {
	cl := &CommandList{}
	var last ICommand
	for i := 0; i < len(names); i++ {
		nc := GetCommand(names[i])
		if nc == nil {
			return nil, fmt.Errorf("command %s not found", names[i])
		}
		if cl.root == nil {
			cl.root = nc
			last = cl.root
		} else if last != nil {
			last.SetNext(nc)
			last = nc
		} else {
			return nil, fmt.Errorf("last point is nil please check it")
		}
	}
	return cl, nil
}

// 执行所有的命令
func (cl *CommandList) Exec(ds *conf.DataSet) {
	if cl.root == nil {
		log.FormatErr("there is no command need to exec")
		return
	}
	// 获取日志
	logger := log.CreateLogger(ds.GetServerId(), conf.IsDebugModel())
	// 创建一个保存结果的类型
	result := report.NewResult()
	// 顺序执行所有的命令
	for c := cl.root; c != nil; c = c.GetNext() {
		// 设置日志实例
		c.SetLogger(logger)
		// 准备工作
		err := c.Prepare(ds)
		if err != nil {
			result.AddFailed(ds.GetIndex(), c.GetName(), c.GetStep(), err.Error())
			break
		}
		// 检查能否执行
		err = c.Check(ds)
		if err != nil {
			result.AddFailed(ds.GetIndex(), c.GetName(), c.GetStep(), err.Error())
			break
		}
		// 执行对应的命令
		rows, err := c.Run(ds)
		if rows != nil {
			result.AddRows(rows)
		}
		if err != nil {
			result.AddFailed(ds.GetIndex(), c.GetName(), c.GetStep(), fmt.Sprintf("%v", err))
			break
		}
	}
	// 报告结果
	result.Report()
}

// 注册函数
func Register(cmd ICommand) {
	allCommand = append(allCommand, cmd)
}

// 获取函数
func GetCommand(name string) ICommand {
	for i := 0; i < len(allCommand); i++ {
		if name == allCommand[i].GetName() {
			return allCommand[i]
		}
	}
	return nil
}
