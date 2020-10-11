package cmds

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/simplemoon/deploy/report"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/simplemoon/deploy/conf"
	"github.com/simplemoon/deploy/utils"
)

// 启服
type CmdStop struct {
	*CommandBase        // 基础的功能
	Contains     uint32 // 需要包含的类型
	idx          int    // 服务器的序号
}

// 创建命令
func NewCmdStop() *CmdStop {
	return &CmdStop{
		CommandBase: NewCommandBase(CmdNameStop),
		Contains:    conf.ProcTypeContainsAll,
	}
}

// 拷贝对应的命令
func (cmd *CmdStop) Copy() ICommand {
	return NewCmdStop()
}

// 检查
func (cmd *CmdStop) Check(ds *conf.DataSet) error {
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
	var name string
	if ds.HasServer() {
		name = conf.ProcessKeyGame
	} else if ds.HasPlatform() {
		name = conf.ProcessKeyPlatform
	} else {
		return fmt.Errorf("not found key process")
	}

	s, err := GetServiceCmdResult(name, idx, conf.EchoCmdStatus)
	if err != nil {
		return err
	}
	// 没有连接上
	if s == utils.ConnectStateFailed {
		return fmt.Errorf("server %d is closed, do not need to close", serverId)
	}
	return nil
}

// 运行对应的命令
func (cmd *CmdStop) Run(ds *conf.DataSet) ([]report.RowInfo, error) {
	// 读取配置文件的IP和端口,然后查询状态
	cmd.SetStep("running")

	// 1. 发送之前的数据
	if err := cmd.beforeRun(ds); err != nil {
		return nil, err
	}
	time.Sleep(time.Second * 5)

	// 2. 关闭关键的数据
	if err := cmd.killKeyProc(ds); err != nil {
		return nil, err
	}

	// 3. 关闭所有的进程
	var wg sync.WaitGroup
	// 执行对应的结果
	err := cmd.ExecServiceCmd(ds, &wg, ServiceCmdNameStop, cmd.Contains)
	if err != nil {
		return nil, err
	}
	wg.Done()

	// 4. 关闭所有属于当前服务器的进程
	err = cmd.CloseReserve(ds)
	if err != nil {
		return cmd.result, err
	}

	// 报告服务器的状态
	err = cmd.ReportStatus(ds)

	return cmd.result, err
}

// 关闭之前的处理
func (cmd *CmdStop) beforeRun(ds *conf.DataSet) error {
	// 发送telnet的消息，关闭服务器
	if ds.HasCharge() {
		_, err := GetServiceCmdResult(conf.ProcessKeyGm, cmd.idx, conf.EchoCmdBeforeQuit)
		if err != nil {
			return err
		}
	}
	if ds.HasServer() {
		_, err := GetServiceCmdResult(conf.ProcessKeyGame, cmd.idx, conf.EchoCmdDisconnectRoom)
		if err != nil {
			return err
		}
	}
	return nil
}

// 发送quit命令，等待进程关闭或者超时
func (cmd *CmdStop) killKeyProc(ds *conf.DataSet) error {
	var wg sync.WaitGroup

	if ds.HasServer() && cmd.Contains&conf.ProcTypeServer > 0 {
		wg.Add(1)
		go cmd.waitProcClose(&wg, conf.ProcessKeyGame)
	}

	if ds.HasPlatform() && cmd.Contains&conf.ProcTypePlatform > 0 {
		wg.Add(1)
		go cmd.waitProcClose(&wg, conf.ProcessKeyPlatform)
	}
	wg.Wait()

	return nil
}

// 关闭对应的进程
func (cmd *CmdStop) waitProcClose(wg *sync.WaitGroup, name string) {
	defer wg.Done()

	// 发送对应的命令
	_, err := GetServiceCmdResult(name, cmd.idx, conf.EchoCmdQuit)
	if err != nil {
		cmd.GetLogger().Infof("send to process %s quit command error %v", name, err)
		return
	}

	// 等待对应的进程关闭
	keyName := conf.GetNameWithoutSuffix(name)
	mm := map[string]int{keyName: 0}
	dirPath := utils.GetServerBinDir(cmd.idx)

	// 关闭的通知
	sig := make(chan struct{})

	go func() {
		defer close(sig)

		for i := 0; i < 90; i++ {
			err = utils.GetProcId(dirPath, &mm)
			if err != nil {
				cmd.GetLogger().Infof("get %s processId err %v", conf.ProcessKeyGame)
				return
			}
			v, ok := mm[keyName]
			if !ok {
				cmd.GetLogger().Infof("get %s processId not at map", keyName)
				return
			}
			if v <= 0 {
				cmd.GetLogger().Infof("get %s processId is zero", keyName)
				return
			}
			cmd.GetLogger().Infof("get %s processId is %d", keyName, v)

			// 等待一段时间
			time.Sleep(time.Second * 2)
		}
	}()

	// 等待关闭啊
	<-sig
}

// 关闭遗留的进程
func (cmd *CmdStop) CloseReserve(ds *conf.DataSet) error {
	// 等待对应的进程关闭
	pl := conf.GetProcInfo(ds, cmd.Contains)

	// 所有需要关闭的实例
	mm := make(map[string]int, len(pl))
	dirPath := utils.GetServerBinDir(cmd.idx)

	// 关闭的通知
	for _, p := range pl {
		mm[p.GetName()] = 0
	}
	// 获取所有的信息
	err := utils.GetProcId(dirPath, &mm)
	if err != nil {
		cmd.GetLogger().Infof("get %s processId err %v", conf.ProcessKeyGame)
		return err
	}

	// 关闭对应的信息
	buf := strings.Builder{}
	for k, v := range mm {
		if v <= 0 {
			continue
		}
		// 写入对应的数据
		c := fmt.Sprintf("/PID %d ", v)

		_, err := buf.WriteString(c)
		if err != nil {
			cmd.GetLogger().Infof("write to buf process %s, pid %d err %v", k, v, err)
			continue
		}
	}

	// 关闭pid
	return utils.KillProcess(buf.String())
}

// 报告服务器的状态
func (cmd *CmdStop) ReportStatus(ds *conf.DataSet) error {
	url := ds.GetBaseUrl()
	if url == "" {
		return fmt.Errorf("report status base url is null")
	}
	// 获取游戏id
	gameId, err := ds.GetString(conf.InfoKeyGameId)
	if err != nil {
		return err
	}
	gId, err := strconv.Atoi(gameId)
	if err != nil {
		return err
	}
	serverId := ds.GetServerId()
	// 创建对应的数据
	data := make(map[string]interface{})
	url += fmt.Sprintf(conf.UrlPathSvrState, serverId)
	// 具体的数据
	data["gameId"] = gId
	if cmd.Contains&conf.ProcTypeServer > 0 {
		data["state"] = "0"
	}
	if cmd.Contains&conf.ProcTypePlatform > 0 {
		data["platformState"] = "0"
	}
	// 具体的数据
	content, err := json.Marshal(data)
	if err != nil {
		return err
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(content))
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("send status to server code error %d", resp.StatusCode)
	}
	return nil
}

// 初始化
func init() {
	Register(NewCmdStop())
}
