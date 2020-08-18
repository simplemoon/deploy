package utils

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

var (
	rootDir = ""
	execDir = ""
)

const (
	// 目录的名称
	DirNameRuntime  = "soda_runtime" // 运行时目录
	DirNameLog      = "log"          // 日志目录
	DirNameRes      = "res"          // 资源目录
	DirNameLock     = "lock"         // 锁文件目录
	DirNameBackUp   = "backup"       // 备份目录
	DirNameServer   = "servers"      // 服务器所在的目录
	DirNameSvrRunAt = "svr"          // 服务器运行的目录

	// 文件的名称
	FileNameLog = "log"
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
		return filepath.Join(dir, name), nil
	}
	return "", err
}

// 获取日志文件路径
func GetLogPath(serverId int) string {
	return path.Join(rootDir, DirNameRuntime, fmt.Sprintf("%s_%d.log", FileNameLog, serverId))
}

// 获取运行时目录的路径
func GetRuntimePath(fileName string) string {
	return path.Join(rootDir, DirNameRuntime, fileName)
}

// 获取服务器的目录
func GetServerDir(idx int) string {
	return path.Join(rootDir, DirNameServer, fmt.Sprintf("%s%d", DirNameSvrRunAt, idx))
}

// 检查对应的目录是否存在啊
func CheckServerDirExist(idx int) bool {
	return IsPathExist(GetServerDir(idx))
}

// 创建需要的文件夹
func CreateDirs() error {
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
