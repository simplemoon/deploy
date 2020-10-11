package utils

import (
	"bytes"
	"os/exec"
	"strconv"
	"strings"
)

// 执行命令
func StartProc(name, dir string, cmd ...string) (string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	// 设置一些参数
	c := exec.Command(name, cmd...)
	c.Stdout = &stdout
	c.Stderr = &stderr
	c.Dir = dir
	// 执行
	err := c.Run()
	if err != nil {
		return ServiceResultFailed, err
	}

	return stdout.String(), nil
}

// 获取对应的进程号
func GetProcId(dir string, mm *map[string]int) error {
	buf := bytes.Buffer{}
	// 获取所有的进程信息
	c := exec.Command("wmic", "process", "get", "processid,name,executablepath")
	c.Stdout = &buf

	// 执行命令
	err := c.Run()
	if err != nil {
		return err
	}
	// 获取对应的信息
	rows := strings.Fields(buf.String())
	for _, row := range rows {
		r := strings.Split(row, ",")
		if len(r) < 3 {
			continue
		}
		// 检查对应的路径是否是对应的目录
		if !strings.Contains(r[2], dir) {
			continue
		}
		name := strings.TrimSuffix(r[1], ".exe")
		if _, ok := (*mm)[name]; !ok {
			continue
		}
		id, err := strconv.Atoi(r[0])
		if err != nil {
			continue
		}
		(*mm)[name] = id
	}
	return nil
}

// 获取对应的进程号
func KillProcess(pidInfos string) error {
	buf := bytes.Buffer{}
	// 获取所有的进程信息
	c := exec.Command("taskkill", "/F", pidInfos)
	c.Stdout = &buf

	// 执行命令
	return c.Run()
}
