package cmds

import (
	"fmt"
	"github.com/simplemoon/deploy/report"

	"github.com/simplemoon/deploy/conf"
	"github.com/simplemoon/deploy/log"
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
		conf.FormatErr("there is no command need to exec")
		return
	}
	// 获取日志
	logger := log.CreateLogger(ds.GetServerId(), conf.IsDebugModel())
	// 创建一个保存结果的类型
	result := report.NewResult(conf.GetProject(), conf.GetVersion())
	// url 的路径
	url := ds.GetBaseUrl()

	// 顺序执行所有的命令
	for c := cl.root; c != nil; c = c.GetNext() {
		// 设置日志实例
		c.SetLogger(logger)
		// 根据名称决定是否需要报告
		name := c.GetName()
		// 准备工作
		err := c.Prepare(ds)
		if err != nil {
			result.AddFailed(ds.GetIndex(), name, c.GetStep(), err.Error())
			break
		}
		// 检查能否执行
		err = c.Check(ds)
		if err != nil {
			result.AddFailed(ds.GetIndex(), name, c.GetStep(), err.Error())
			break
		}
		// 执行对应的命令
		rows, err := c.Run(ds)
		if rows != nil {
			result.AddRows(rows)
		}
		if err != nil {
			result.AddFailed(ds.GetIndex(), name, c.GetStep(), fmt.Sprintf("%v", err))
			break
		}
		// 发送钉钉消息
		if name == CmdNameStart || name == CmdNameStop || name == CmdNameUpdate {
			if url == "" {
				logger.Error("not found base url")
				continue
			}
			gameId, err := ds.GetString(conf.InfoKeyGameId)
			if err != nil {
				logger.Error("not found gameId")
				return
			}
			dingUrl := url + conf.UrlPathDing
			err = result.DingNotify(dingUrl, gameId, name, ds.GetServerId())
			if err != nil {
				logger.Info(err)
			}
		}
	}
	// 报告结果
	if url == "" {
		logger.Errorf("not found base url")
		return
	}
	resultUrl := url + conf.UrlPathResult
	// 报告对应的结果信息
	err := result.Report(resultUrl, conf.IsRemoteModel())
	if err != nil {
		logger.Errorf("report data err %v", err)
	}
}

// 注册函数
func Register(cmd ICommand) {
	allCommand = append(allCommand, cmd)
}

// 获取函数
func GetCommand(name string) ICommand {
	for i := 0; i < len(allCommand); i++ {
		if name == allCommand[i].GetName() {
			return allCommand[i].Copy()
		}
	}
	return nil
}
