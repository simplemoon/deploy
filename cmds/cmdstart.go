package cmds

import (
	"context"
	"fmt"
	"github.com/simplemoon/deploy/report"
	"sync"
	"time"

	"github.com/simplemoon/deploy/conf"
	"github.com/simplemoon/deploy/utils"
)

// 启服
type CmdStart struct {
	*CommandBase     // 基础的功能
	idx          int // 服务器的序号
}

// 创建命令
func NewCmdStart() *CmdStart {
	return &CmdStart{
		CommandBase: NewCommandBase(CmdNameStart),
	}
}

func (cmd *CmdStart) Copy() ICommand {
	return NewCmdStart()
}

// 检查
func (cmd *CmdStart) Check(ds *conf.DataSet) error {
	// 基本的检查
	cmd.SetStep("check")
	serverId := ds.GetServerId()
	idx, err := conf.GetIndexByServerId(serverId)
	if err != nil {
		return err
	}
	// 设置以下idx
	cmd.idx = idx
	// 检查基础的是否存在
	if err := cmd.CheckByIndex(serverId, idx); err != nil {
		return err
	}
	// 检查服务器是否已经开启了
	s, err := GetServiceCmdResult(conf.ProcessKeyGame, idx, conf.EchoCmdStatus)
	if err != nil {
		return err
	}
	if s == utils.ConnectStateOpened {
		return fmt.Errorf("server %d is started, no to start", serverId)
	}
	return nil
}

// 运行对应的命令
func (cmd *CmdStart) Run(ds *conf.DataSet) ([]report.RowInfo, error) {
	// 读取配置文件的IP和端口,然后查询状态
	cmd.SetStep("running")
	var wg sync.WaitGroup

	// 执行对应的结果
	err := cmd.ExecServiceCmd(ds, &wg, ServiceCmdNameStart, conf.ProcTypeContainsAll)
	if err != nil {
		return nil, err
	}

	// 检测是否需要开启
	wg.Add(1)
	go cmd.waitOpened(&wg)
	// 等待所有的执行完毕
	wg.Wait()
	cmd.SetStep("end")
	return cmd.result, nil
}

// 等待着完全开启
func (cmd *CmdStart) waitOpened(wg *sync.WaitGroup) {
	defer wg.Done()
	// 发送 status 检查状态
	ctx, _ := context.WithTimeout(context.Background(), time.Minute*3)
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ctx.Done():
			cmd.GetLogger().Printf("wait %s finish timeout\n", cmd.GetName())
			ticker.Stop()
			return
		case <-ticker.C:
			s, err := GetServiceCmdResult(conf.ProcessKeyGame, cmd.idx, conf.EchoCmdStatus)
			if err != nil {
				cmd.GetLogger().Printf("query %s %s -> %v\n", conf.ProcessKeyGame, conf.EchoCmdStatus, err)
				continue
			}
			if s == utils.ConnectStateOpened {
				cmd.GetLogger().Printf("%s opened", conf.ProcessKeyGame)
				return
			}
			cmd.GetLogger().Printf("query %s %s -> %v\n", conf.ProcessKeyGame, conf.EchoCmdStatus, err)
		}
	}
}

// 初始化
func init() {
	Register(NewCmdStart())
}
