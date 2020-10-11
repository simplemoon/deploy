package cmds

import (
	"fmt"
	"github.com/simplemoon/deploy/report"
	"os"
	"path"
	"time"

	"github.com/simplemoon/deploy/conf"
	"github.com/simplemoon/deploy/utils"
)

type CmdUpdate struct {
	*CmdInit
	idx int // 服务器的序号
}

func NewCmdUpdate() *CmdUpdate {
	c := &CmdUpdate{
		CmdInit: NewCmdInit(),
	}
	c.SetName(CmdNameUpdate)
	return c
}

func (cmd *CmdUpdate) Copy() ICommand {
	return NewCmdUpdate()
}

// 检查，直接通过，或者检查一下版本信息
func (cmd *CmdUpdate) Check(ds *conf.DataSet) error {
	return nil
}

// 更新
func (cmd *CmdUpdate) Run(ds *conf.DataSet) (ret []report.RowInfo, err error) {
	ret = cmd.result
	// 1. 生成配置文件
	if err = cmd.tryCreateIndex(ds); err != nil {
		return
	}
	// 2. 下载对应的资源到对应的目录之中
	if err = cmd.downloadRes(ds); err != nil {
		return
	}
	// 3. 解压到对应的目录
	if err = cmd.unzipFile(ds); err != nil {
		return
	}
	// 4. 执行 init 的操作
	if err = cmd.CmdInit.Check(ds); err != nil {
		return
	}
	return
}

// 创建一个index,如果没有的话
func (cmd *CmdUpdate) tryCreateIndex(ds *conf.DataSet) error {
	// 获取文件锁
	cmd.SetStep("create index")
	fl := utils.NewFileLock(utils.GetRuntimePath(conf.FileRuntimeServerLock))
	err := fl.LockWithTime(time.Second * 30)
	if err != nil {
		return err
	}
	defer fl.UnLock()

	// 创建服务器的序号
	idx, err := conf.CreateServerIdx(ds.GetServerId(), false)
	if err != nil {
		return err
	}
	cmd.idx = idx
	cmd.AddResult(&report.RowInfo{
		ServerIdx:  idx,
		Action:     cmd.GetName(),
		ActionType: "create index",
		State:      report.Success,
		Msg:        fmt.Sprintf("%v create index %d", ds.GetServerId(), idx),
	})
	return nil
}

// 下载资源
func (cmd *CmdUpdate) downloadRes(ds *conf.DataSet) error {
	cmd.SetStep("download res")
	// 获取对应的url
	value := ds.Get(conf.InfoKeyResUrl)
	url, ok := value.(string)
	if !ok {
		return fmt.Errorf("dataset key: %s data type must be string", conf.InfoKeyResUrl)
	}
	name := path.Base(url)
	// 下载对应的文件
	err := conf.DownloadZips(url, ds.Get(conf.InfoKeyMd5))
	if conf.IsExistErr(err) {
		cmd.AddResult(&report.RowInfo{
			ServerIdx:  cmd.idx,
			Action:     cmd.GetName(),
			ActionType: "download res",
			State:      report.Success,
			Msg:        fmt.Sprintf("%s alread exist", name),
		})
		return nil
	}
	// 正常的逻辑
	if err != nil {
		return err
	}
	cmd.AddResult(&report.RowInfo{
		ServerIdx:  cmd.idx,
		Action:     cmd.GetName(),
		ActionType: "download res",
		State:      report.Success,
		Msg:        fmt.Sprintf("%d download %s succeed", ds.GetServerId(), name),
	})
	return nil
}

// 解压对应的文件
func (cmd *CmdUpdate) unzipFile(ds *conf.DataSet) error {
	cmd.SetStep("unzip res")
	// 获取对应的url
	value := ds.Get(conf.InfoKeyResUrl)
	url, ok := value.(string)
	if !ok {
		return fmt.Errorf("dataset key: %s data type must be string", conf.InfoKeyResUrl)
	}
	name := path.Base(url)
	// 获取res的目录
	resDir := utils.GetRuntimeResDir()
	resPath := path.Join(resDir, name)
	// 获取服务器的目录
	svrDir := utils.GetServiceDir(cmd.idx)
	// 获取备份的目录
	backDir := utils.GetRuntimeServerBackupDir(cmd.idx)

	if utils.IsPathExist(svrDir) {
		// 备份对应的文件
		err := utils.Compress(svrDir, backDir)
		if err != nil {
			return err
		}
		// 删除目录
		err = os.RemoveAll(svrDir)
		if err != nil {
			return err
		}
	}
	// 解压文件到相应的目录
	if err := utils.DeCompress(resPath, svrDir, nil); err != nil {
		return err
	}
	cmd.AddResult(&report.RowInfo{
		ServerIdx:  cmd.idx,
		Action:     cmd.GetName(),
		ActionType: "unzip res",
		State:      report.Success,
		Msg:        fmt.Sprintf("%d unzip %s succeed", ds.GetServerId(), name),
	})
	return nil
}

// 注册
func init() {
	Register(NewCmdUpdate())
}
