package cmds

import (
	"fmt"
	"github.com/simplemoon/deploy/report"
	"path/filepath"
	"time"

	"github.com/simplemoon/deploy/conf"
	"github.com/simplemoon/deploy/utils"
)

// 状态查询
type CmdCTime struct {
	*CommandBase
}

func NewCmdCTime() *CmdCTime {
	return &CmdCTime{
		CommandBase: NewCommandBase(CmdNameCTime),
	}
}

func (cmd *CmdCTime) Copy() ICommand {
	return NewCmdCTime()
}

// 运行
func (cmd *CmdCTime) Run(ds *conf.DataSet) (ret []report.RowInfo, err error) {
	ret = cmd.result
	// 时间
	ts := conf.GetTimeStampAt()
	if ts == "" {
		ts, err = ds.GetString(conf.InfoKeyToolArgs)
	}
	if err != nil {
		return
	}
	// 转换成时间格式
	t, err := time.Parse("2006-01-02 15:04:05", ts)
	if err != nil {
		return
	}
	wt := t.Unix()
	nt := time.Now().Unix()
	if wt < nt {
		err = fmt.Errorf("%s is before than now", ts)
		return
	}
	diff := int(wt - nt)
	// 获取服务器的目录
	idx := ds.GetIndex()
	svrDir := utils.GetServerDir(idx)
	file := filepath.Join(svrDir, utils.DirNameExtra, conf.FileExtraSystemTime)
	if utils.IsPathExist(file) {
		err = fmt.Errorf("can't found %s", file)
		return
	}
	gameId, err := ds.GetString(conf.InfoKeyGameId)
	if err != nil {
		return
	}
	// 重置时间配置
	tm := conf.NewTimeMgr(file, gameId, ds.GetServerId(), ds.GetBaseUrl())
	err = tm.SaveTime(diff)
	if err != nil {
		return
	}
	cmd.result = append(cmd.result, report.RowInfo{
		ServerIdx:  idx,
		Action:     cmd.Name,
		ActionType: "change diff time",
		State:      report.Success,
		Msg:        fmt.Sprintf("server ahead second %d, time %s", diff, ts),
	})
	return
}

func init() {
	Register(NewCmdCTime())
}
