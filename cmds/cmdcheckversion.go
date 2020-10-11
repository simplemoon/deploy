package cmds

import (
	"fmt"
	"github.com/simplemoon/deploy/conf"
	"github.com/simplemoon/deploy/report"
	"github.com/simplemoon/deploy/utils"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type CmdCheckVersion struct {
	*CommandBase
}

func NewCmdCheckVersion() *CmdCheckVersion {
	return &CmdCheckVersion{
		CommandBase: NewCommandBase(CmdNameCheckVersion),
	}
}

func (cmd *CmdCheckVersion) Copy() ICommand {
	return NewCmdCheckVersion()
}

func (cmd *CmdCheckVersion) Run(ds *conf.DataSet) (ret []report.RowInfo, err error) {
	ret = cmd.result
	// 检查对应的版本信息
	err = cmd.checkVersion(ds)
	return
}

func (cmd *CmdCheckVersion) checkVersion(ds *conf.DataSet) error {
	ver := ds.Get(conf.InfoKeyVersion)
	if ver == nil {
		return fmt.Errorf("don't found key of %s", conf.InfoKeyVersion)
	}
	v, ok := ver.(string)
	if !ok {
		return fmt.Errorf("%v convert to string failed", v)
	}
	// 读取文件
	idx := ds.GetIndex()
	if idx <= 0 {
		return fmt.Errorf("don't found server index")
	}
	// 读取文件信息
	verPath := filepath.Join(utils.GetServerDir(idx), conf.FileVersion)
	data, err := ioutil.ReadFile(verPath)
	if err != nil {
		return err
	}
	cv := strings.TrimSpace(string(data))
	if cv != v {
		return fmt.Errorf("now version is %s, verify version is %s", cv, v)
	}
	cmd.AddResult(&report.RowInfo{
		ServerIdx:  idx,
		Action:     cmd.GetName(),
		ActionType: "check version",
		State:      report.Success,
		Msg:        "version is same with oss shows",
	})
	return nil
}

func init() {
	Register(NewCmdCheckVersion())
}
