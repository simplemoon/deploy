package utils

import (
	"strings"
)

// 执行结果
const (
	ServiceResultSuccess      = "Success"
	ServiceResultFailed       = "Failed"
	ServiceResultRunning      = "AlreadyRunning"
	ServiceResultNotFound     = "NoSuchService"
	ServiceResultAlreadyStop  = "AlreadyStop"
	ServiceResultAccessDenied = "AccessDenied"
	ServiceResultInstalled    = "AlreadyInstall"
)

// 执行状态
const (
	ServiceStatusAlreadyRunning = "ServiceAlreadyRunning"
	ServiceStatusNoFound        = "NoSuchService"
	ServiceStatusCantControl    = "ServiceCannotAcceptControl"
	ServiceStatusAccessDenied   = "AccessDenied"
	ServiceStatusAlreadyExists  = "already exists"
)

// 默认的对应的map
var (
	// 检查的类型
	globalKeyResult = map[string]string{
		ServiceStatusAlreadyRunning: ServiceResultRunning,
		ServiceStatusNoFound:        ServiceResultNotFound,
		ServiceStatusCantControl:    ServiceResultAlreadyStop,
		ServiceStatusAccessDenied:   ServiceResultAccessDenied,
		ServiceStatusAlreadyExists:  ServiceResultInstalled,
	}
)

func getResult(s string) string {
	for k, v := range globalKeyResult {
		if strings.Contains(s, k) {
			return v
		}
	}
	return ServiceResultSuccess
}

// 执行命令
func Exec(name, dir, cmd string) (string, error) {
	result, err := StartProc(name, dir, cmd)
	if err != nil {
		return ServiceResultFailed, err
	}

	return getResult(result), nil
}
