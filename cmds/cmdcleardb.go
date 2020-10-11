package cmds

import (
	"fmt"
	"github.com/simplemoon/deploy/conf"
	"github.com/simplemoon/deploy/report"
	"github.com/simplemoon/deploy/utils"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	DbResultSuccess = "0"
)

// 数据库的信息
type DB struct {
	Addr string // 地址
	User string // 用户名
	Pwd  string // 密码
	Port string // 端口号
}

type CmdClearDB struct {
	*CommandBase
}

func NewCmdClearDB() *CmdClearDB {
	return &CmdClearDB{
		CommandBase: NewCommandBase(CmdNameClearDB),
	}
}

func (cmd *CmdClearDB) Copy() ICommand {
	return NewCmdClearDB()
}

func (cmd *CmdClearDB) Run(ds *conf.DataSet) (ret []report.RowInfo, err error) {
	ret = cmd.result
	if err = cmd.clearDataBase(ds); err != nil {
		return
	}
	if err = cmd.resetTime(ds); err != nil {
		return
	}
	return
}

// 清理数据库
func (cmd *CmdClearDB) clearDataBase(ds *conf.DataSet) error {
	db, err := cmd.getDb(ds)
	if err != nil {
		return err
	}
	// 执行对应的命令
	strSid := strconv.Itoa(ds.GetServerId())
	idx := ds.GetIndex()
	if idx <= 0 {
		return fmt.Errorf("don't find server index of %v", ds.GetServerId())
	}
	// 执行文件的路径
	workDir, err := utils.GetExePath()
	if err != nil {
		return err
	}
	svrDir := utils.GetServerDir(idx)
	// 原来的脚本id
	sqlDir := filepath.Join(svrDir, utils.DirNameTools, utils.DirNameSql)
	if !utils.IsPathExist(sqlDir) {
		sqlDir = filepath.Join(workDir, utils.DirNameScripts, utils.DirNameSql, conf.ProjectName)
	}
	// 主要的执行的文件
	scriptFile := filepath.Join(workDir, utils.DirNameScripts, conf.FileNameResetDb)
	// 执行对应的命令
	ret, err := utils.StartProc(scriptFile, sqlDir, db.User, db.Pwd, db.Addr, db.Port, strSid, strSid)
	if err != nil {
		return err
	}
	if ret != DbResultSuccess {
		return fmt.Errorf("exec clear_db failed code %s", ret)
	}
	return nil
}

// 重置时间配置
func (cmd *CmdClearDB) resetTime(ds *conf.DataSet) error {
	// 获取服务器的目录
	idx := ds.GetIndex()
	svrDir := utils.GetServerDir(idx)
	file := filepath.Join(svrDir, utils.DirNameExtra, conf.FileExtraSystemTime)
	if utils.IsPathExist(file) {
		return fmt.Errorf("can't found %s", file)
	}
	gameId, err := ds.GetString(conf.InfoKeyGameId)
	if err != nil {
		return err
	}
	// 重置时间配置
	tm := conf.NewTimeMgr(file, gameId, ds.GetServerId(), ds.GetBaseUrl())
	return tm.SaveTime(0)
}

// 获取数据库信息
func (cmd *CmdClearDB) getDb(ds *conf.DataSet) (*DB, error) {
	db := new(DB)
	var err error
	db.User, err = ds.GetString(conf.InfoKeyDBUser)
	if err != nil {
		return nil, err
	}
	db.User, err = ds.GetString(conf.InfoKeyDBPwd)
	if err != nil {
		return nil, err
	}
	addr, err := ds.GetString(conf.InfoKeyDBAddr)
	if err != nil {
		return nil, err
	}
	r := strings.Split(addr, ":")
	rn := len(r)
	if rn < 1 {
		return nil, fmt.Errorf("get database address error")
	}
	db.Addr = r[0]
	if rn >= 2 {
		db.Port = r[1]
	} else {
		db.Port = conf.DefaultDataBasePort
	}
	return db, nil
}

func init() {
	Register(NewCmdClearDB())
}
