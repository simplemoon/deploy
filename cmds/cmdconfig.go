package cmds

import (
	"github.com/simplemoon/deploy/conf"
	"github.com/simplemoon/deploy/report"
	"github.com/simplemoon/deploy/utils"
	"time"
)

type CmdConfig struct {
	*CmdDelete
}

func NewCmdConfig() *CmdConfig {
	cmd := &CmdConfig{
		CmdDelete: NewCmdDelete(),
	}
	cmd.SetName(CmdNameConfig)
	return cmd
}

func (cmd *CmdConfig) Copy() ICommand {
	return NewCmdConfig()
}

func (cmd *CmdConfig) Run(ds *conf.DataSet) (ret []report.RowInfo, err error) {
	ret = cmd.result
	// 1. 删除服务
	if err = cmd.deleteService(ds); err != nil {
		return
	}
	// 2. 备份目录
	if err = cmd.backDir(ds); err != nil {
		return
	}
	// 3. 重置一下json文件
	if err = cmd.resetJson(ds); err != nil {
		return
	}
	return
}

func (cmd *CmdConfig) resetJson(ds *conf.DataSet) error {
	cmd.SetStep("reset json")
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
	// 拷贝文件
	destPath := utils.GetRuntimeConfigBackupPath()
	err = conf.CopyFile(filePath, destPath)
	if err != nil {
		return err
	}
	// 拷贝文件
	_, err = conf.GetIndexByServerId(ds.GetServerId())
	return err
}

func init() {
	Register(NewCmdConfig())
}
