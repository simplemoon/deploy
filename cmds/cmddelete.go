package cmds

import (
	"encoding/json"
	"fmt"
	"github.com/simplemoon/deploy/report"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/simplemoon/deploy/conf"
	"github.com/simplemoon/deploy/utils"
)

type CmdDelete struct {
	*CommandBase
	idx int // 服务器的编号
}

// 删除服务器
func NewCmdDelete() *CmdDelete {
	return &CmdDelete{
		CommandBase: NewCommandBase(CmdNameDelete),
	}
}

// 检查服务器的状态
func (cmd *CmdDelete) Check(ds *conf.DataSet) error {
	cmd.SetStep("check")
	// 获取idx
	serverId := ds.GetServerId()
	idx, err := conf.GetIndexByServerId(serverId)
	if err != nil {
		return err
	}
	// 检查基础的是否存在
	if err := cmd.CheckByIndex(serverId, idx); err != nil {
		return err
	}
	cmd.idx = idx
	// 检查服务器是否已经开启了
	var keyName string
	if ds.HasServer() {
		keyName = conf.ProcessKeyGame
	} else {
		keyName = conf.ProcessKeyPlatform
	}
	s, err := GetServiceCmdResult(keyName, idx, conf.EchoCmdStatus)
	if err != nil {
		return err
	}
	if s == utils.ConnectStateOpened {
		return fmt.Errorf("server %d is started, can not delete", serverId)
	}
	return nil
}

// 拷贝命令
func (cmd *CmdDelete) Copy() ICommand {
	return NewCmdDelete()
}

func (cmd *CmdDelete) Run(ds *conf.DataSet) (ret []report.RowInfo, err error) {
	ret = cmd.result
	// 1. 删除对应的服务
	if err = cmd.deleteService(ds); err != nil {
		return
	}
	// 2. 备份对应的文件夹
	if err = cmd.backDir(ds); err != nil {
		return
	}
	// 3. 删除对应的 json 文件
	if err = cmd.removeJson(ds); err != nil {
		return
	}
	return
}

// 删除对应的服务
func (cmd *CmdDelete) deleteService(ds *conf.DataSet) error {
	cmd.SetStep("delete service")
	var wg sync.WaitGroup
	// 执行对应的结果
	err := cmd.ExecServiceCmd(ds, &wg, ServiceCmdNameUnInstall, conf.ProcTypeContainsAll)
	if err != nil {
		return err
	}
	wg.Wait()
	return nil
}

// 备份server 文件夹
func (cmd *CmdDelete) backDir(ds *conf.DataSet) error {
	cmd.SetStep("back up")
	// 备份对应的目录
	serverDir := utils.GetServerDir(cmd.idx)
	if !utils.IsPathExist(serverDir) {
		return nil
	}
	// 备份文件的目录
	sid := ds.GetServerId()
	backUpDir := utils.GetRuntimeServerBackupDir(sid)
	// 压缩到备份文件
	err := utils.Compress(serverDir, backUpDir)
	if err != nil {
		return err
	}
	return os.RemoveAll(serverDir)
}

// 删除 json 文件
func (cmd *CmdDelete) removeJson(ds *conf.DataSet) error {
	cmd.SetStep("remove json")
	// 获取 server.json
	filePath := utils.GetRuntimePath(conf.FileRuntimeServers)
	if !utils.IsPathExist(filePath) {
		return nil
	}
	// 删除 json 文件
	fl := utils.NewFileLock(utils.GetRuntimePath(conf.FileRuntimeServerLock))
	err := fl.LockWithTime(time.Second * 30)
	if err != nil {
		return err
	}
	defer fl.UnLock()
	// 获取拷贝后的文件名称
	destPath := utils.GetRuntimeConfigBackupPath()
	err = conf.CopyFile(filePath, destPath)
	if err != nil {
		return err
	}
	// 拷贝文件
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	// 获取对应的信息
	servers := make([]conf.ServerInfo, 0)
	err = json.Unmarshal(data, &servers)
	if err != nil {
		return err
	}
	sid := ds.GetServerId()
	for n, s := range servers {
		if s.ServerId == sid {
			servers = append(servers[:n], servers[n+1:]...)
			break
		}
	}
	// 写入到文件之中
	content, err := json.Marshal(servers)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filePath, content, 0666)
}

func init() {
	Register(NewCmdDelete())
}
