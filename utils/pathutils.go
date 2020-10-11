package utils

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"
)

var (
	rootDir = ""
	execDir = ""
)

const (
	DirRunCfg    = "cfg" // 配置路径
	DirRunSubIni = "ini" // ini文件配置路径
)

const (
	// 目录的名称
	DirNameRuntime = "soda_runtime" // 运行时目录
	DirNameLog     = "log"          // 日志目录
	DirNameRes     = "res"          // 资源目录
	DirNameLock    = "lock"         // 锁文件目录
	DirNameBackUp  = "backup"       // 备份目录
	DirNameServer  = "servers"      // 服务器所在的父目录
	DirNameScripts = "scripts"      // 脚本目录
	DirNameConfigs = "configs"      // config 的配置目录

	DirNameSvrRunAt = "svr"        // 服务器所在的目录
	DirNameBin64    = "bin64"      // 服务器运行的目录
	DirNameServices = "services"   // 服务程序所在的目录
	DirNameTools    = "tools"      // 工具目录
	DirNameSql      = "sql"        // 数据库目录
	DirNameExtra    = "extra"      // 可以改变的配置目录
	DirNameDebugger = "SdDebugger" // 蚂蚁工具所在的目录
)

const (
	FileNameLog  = "log"     // 文件的名称
	FileNameBase = "servers" // 服务器ID对应的文件名

	ServiceSuffixExe = "_service.exe" // service 执行的后缀
	ServiceSuffixXml = "_service.xml" // service 配置文件的后缀

	ServerSuffixExe = ".exe"
)

// 设置根目录
func SetRootPath(root string) {
	rootDir = root
}

// 文件是否存在
func IsPathExist(name string) bool {
	_, err := os.Stat(name)
	if err == nil {
		return true
	}
	return false
}

// 创建文件夹
func CreateDir(name string) error {
	if IsPathExist(name) {
		return nil
	}
	return os.Mkdir(name, 0666)
}

// 获取执行文件的路径
func GetExePath() (string, error) {
	if execDir != "" {
		return execDir, nil
	}
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}
	execDir = dir
	return execDir, nil
}

// 获取执行的路径
func GetCfgPath(name string) (string, error) {
	dir, err := GetExePath()
	if err == nil {
		return filepath.Join(dir, DirRunCfg, name), nil
	}
	return "", err
}

// 获取配置的ini路径
func GetCfgIniDir() (string, error) {
	ed, err := GetExePath()
	if err != nil {
		return "", err
	}
	return filepath.Join(ed, DirRunCfg, DirRunSubIni), nil
}

// 获取日志文件路径
func GetLogPath(serverId int) string {
	return path.Join(rootDir, DirNameRuntime, fmt.Sprintf("%s_%d.log", FileNameLog, serverId))
}

// 获取运行时目录的路径
func GetRuntimePath(fileName string) string {
	return path.Join(rootDir, DirNameRuntime, fileName)
}

// 获取运行的资源目录
func GetRuntimeResDir() string {
	return path.Join(rootDir, DirNameRuntime, DirNameRes)
}

// 获取运行的锁文件目录
func GetRuntimeLockDir() string {
	return path.Join(rootDir, DirNameRuntime, DirNameLock)
}

// 备份文件的目录
func GetRuntimeConfigBackupPath() string {
	timeStamp := time.Now().Format("20060102150405")
	name := fmt.Sprintf("%s_%s.json", FileNameBase, timeStamp)
	return path.Join(rootDir, DirNameRuntime, DirNameBackUp, DirNameConfigs, name)
}

// 备份文件的目录
func GetRuntimeServerBackupDir(sid int) string {
	timeStamp := time.Now().Format("20060102150405")
	name := fmt.Sprintf("%s%d_%s.zip", DirNameSvrRunAt, sid, timeStamp)
	return path.Join(rootDir, DirNameRuntime, DirNameBackUp, DirNameServer, name)
}

// 获取服务器的目录
func GetServerDir(idx int) string {
	return path.Join(rootDir, DirNameServer, fmt.Sprintf("%s%d", DirNameSvrRunAt, idx))
}

// 获取服务器的bin目录
func GetServerBinDir(idx int) string {
	return path.Join(GetServerDir(idx), DirNameBin64)
}

// 获取service的目录
func GetServiceDir(idx int) string {
	return path.Join(GetServerDir(idx), DirNameBin64, DirNameServices)
}

// 获取service文件的名称
func GetServiceFileName(name string) string {
	return name + ServiceSuffixExe
}

// 检查对应的目录是否存在啊
func CheckServerDirExist(idx int) bool {
	return IsPathExist(GetServerDir(idx))
}

// 获取ini配置路径
func GetServerIniConfigPath(name string, idx int) string {
	serverDir := GetServerDir(idx)
	return path.Join(serverDir, name)
}

// 创建需要的文件夹
func PrepareDirs() error {
	// 运行时目录
	runtimeDir := path.Join(rootDir, DirNameRuntime)
	// 服务器目录
	serverDir := path.Join(rootDir, DirNameServer)
	err := os.MkdirAll(serverDir, os.ModePerm)
	if err != nil {
		return err
	}
	// 日志目录
	logDir := path.Join(runtimeDir, DirNameLog)
	// 创建所有目录
	err = os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		return err
	}
	// 资源目录
	resDir := path.Join(runtimeDir, DirNameRes)
	err = os.MkdirAll(resDir, os.ModePerm)
	if err != nil {
		return err
	}
	// 锁文件目录
	lockDir := path.Join(runtimeDir, DirNameLock)
	err = os.MkdirAll(lockDir, os.ModePerm)
	if err != nil {
		return err
	}
	// 备份文件目录
	backupDir := path.Join(runtimeDir, DirNameBackUp)
	err = os.MkdirAll(backupDir, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
