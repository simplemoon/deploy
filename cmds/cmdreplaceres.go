package cmds

import (
	"encoding/json"
	"fmt"
	"github.com/simplemoon/deploy/report"
	"io/ioutil"
	"path"
	"strings"

	"github.com/simplemoon/deploy/conf"
	"github.com/simplemoon/deploy/utils"
)

// 需要替换的配置
type ReplaceFile struct {
	ExcludePath []string `json:"excludePath"` // 需要排除的路径
	Includes    []string `json:"includes"`    // 包含的目录
}

// 是否需要替换
func (r *ReplaceFile) CanCopy(p string) bool {
	for _, s := range r.Includes {
		if strings.Contains(p, s) {
			return true
		}
	}
	for _, s := range r.ExcludePath {
		if strings.Contains(p, s) {
			return true
		}
	}
	return false
}

// 替换资源
type CmdReplaceRes struct {
	*CommandBase
	resName string // 资源的名称
}

func NewCmdReplaceRes() *CmdReplaceRes {
	return &CmdReplaceRes{
		CommandBase: NewCommandBase(CmdNameReplaceRes),
	}
}

func (cmd *CmdReplaceRes) Copy() ICommand {
	return NewCmdReplaceRes()
}

func (cmd *CmdReplaceRes) Run(ds *conf.DataSet) (ret []report.RowInfo, err error) {
	cmd.SetStep("running")
	ret = cmd.result
	// 1. 下载对应的资源包
	if err = cmd.downloadZip(ds); err != nil {
		return
	}
	// 2. 解压并且替换对应的资源包
	if err = cmd.replace(ds); err != nil {
		return
	}
	return
}

// 解压对应的资源包
func (cmd *CmdReplaceRes) downloadZip(ds *conf.DataSet) error {
	cmd.SetStep("download res")
	// 获取对应的url
	idx := ds.GetIndex()
	if idx <= 0 {
		return fmt.Errorf("server not exist")
	}
	value := ds.Get(conf.InfoKeyResUrl)
	url, ok := value.(string)
	if !ok {
		return fmt.Errorf("dataset key: %s data type must be string", conf.InfoKeyResUrl)
	}
	name := path.Base(url)
	cmd.resName = name
	// 下载对应的文件
	err := conf.DownloadZips(url, ds.Get(conf.InfoKeyMd5))
	if conf.IsExistErr(err) {
		cmd.AddResult(&report.RowInfo{
			ServerIdx:  idx,
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
		ServerIdx:  idx,
		Action:     cmd.GetName(),
		ActionType: "download res",
		State:      report.Success,
		Msg:        fmt.Sprintf("%d download %s succeed", ds.GetServerId(), name),
	})
	return nil
}

// 替换资源
func (cmd *CmdReplaceRes) replace(ds *conf.DataSet) error {
	cmd.SetStep("replace res")
	// 获取需要拷贝的文件列表
	rp, err := utils.GetCfgPath(conf.FileNameReplace)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadFile(rp)
	if err != nil {
		return err
	}
	var cfg ReplaceFile
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	// 解压对应的文件
	idx := ds.GetIndex()
	svrDir := utils.GetServerDir(idx)
	// 获取res的目录
	resDir := utils.GetRuntimeResDir()
	resPath := path.Join(resDir, cmd.resName)
	// 解压到对应的目录
	err = utils.DeCompress(resPath, svrDir, cfg.CanCopy)
	return err
}

func init() {
	Register(NewCmdReplaceRes())
}
