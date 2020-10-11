package conf

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"gopkg.in/ini.v1"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/simplemoon/deploy/utils"
)

var (
	errZipExist = fmt.Errorf("zip file already exist")
)

type KVMap = map[string]string
type SectionMap = map[string]KVMap
type CfgSvr = map[string]SectionMap

// 是否需要拷贝
func NeedCopy(name string, m map[string]bool) bool {
	s := strings.Split(name, FilterString)
	if len(s) < 3 {
		return true
	}
	if s[0] != utils.DirNameBin64 || s[len(s)-2] != utils.DirNameDebugger {
		return true
	}
	p := filepath.Base(s[len(s)-1])
	for k, v := range m {
		if strings.Contains(p, k) && v {
			return true
		}
	}
	return false
}

// 拷贝文件
func CopyFile(src, dst string) error {
	if src == "" || dst == "" {
		return fmt.Errorf("copy file source or target is nil")
	}
	if !utils.IsPathExist(src) {
		return fmt.Errorf("copy file %s not exist", src)
	}
	// 打开文件
	rd, err := os.Open(src)
	if err != nil {
		return err
	}
	defer rd.Close()

	// 创建需要写入的文件
	wd, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 644)
	if err != nil {
		return err
	}
	_, err = io.Copy(wd, rd)
	return err
}

// 替换 ini 内容
func ReplaceContent(src, dst string, ds *DataSet, d SectionMap) error {
	if src == "" || dst == "" {
		return fmt.Errorf("copy file source or target is nil")
	}
	if !utils.IsPathExist(src) {
		return fmt.Errorf("copy file %s not exist", src)
	}
	// 读取文件
	content, err := ini.Load(src)
	if err != nil {
		return err
	}
	// 写入文件
	for name, data := range d {
		sec, err := content.GetSection(name)
		if err != nil {
			return fmt.Errorf("%s get section %s, error: %v", src, name, err)
		}
		// 写入对应的数据
		for key, value := range data {
			v, err := ds.GetValue(value)
			if err != nil {
				return fmt.Errorf("%s get section: %s, key: %s, value: %s, error: %v", src, name, key, value, err)
			}
			tk, err := sec.GetKey(key)
			if err != nil {
				return fmt.Errorf("%s section: %s, key: %s not found, error: %v", src, name, key, err)
			}
			tk.SetValue(v)
		}
	}
	err = content.SaveToIndent(dst, IndentString)
	return err
}

func IsExistErr(err error) bool {
	return err == errZipExist
}

// 下载对应的文件
func DownloadZips(url string, verify interface{}) error {
	name := path.Base(url)
	// 获取res的目录
	resDir := utils.GetRuntimeResDir()
	// 验证一下文件内容
	resPath := path.Join(resDir, name)
	if checkFileFit(resPath, verify) {
		return errZipExist
	}
	// 锁文件的路径
	lockDir := utils.GetRuntimeLockDir()
	// lock 文件名称
	lockName := strings.Replace(name, FileSuffixZip, FileSuffixLock, 1)
	lockPath := path.Join(lockDir, lockName)
	// 创建一个锁
	fl := utils.NewFileLock(lockPath)
	err := fl.LockWithTime(time.Second * 30)
	if err != nil {
		return err
	}
	defer fl.UnLock()
	// 下载对应的文件
	rsp, err := http.Get(url)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}
	// 打开文件
	err = ioutil.WriteFile(resPath, data, 0666)
	if err != nil {
		return err
	}
	return nil
}

// 检查一下文件是否满足要求
func checkFileFit(name string, md5str interface{}) bool {
	if utils.IsPathExist(name) {
		return false
	}
	// 如果有md5就检查
	if md5str == nil {
		return true
	}
	val, ok := md5str.(string)
	if !ok {
		return false
	}
	// 检查文件的 md5
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return false
	}
	entity := md5.New()
	m := hex.EncodeToString(entity.Sum(data))
	if m == val {
		return true
	}
	return false
}

func GetIniFileName(name string) string {
	return fmt.Sprintf("%s%s.ini", PreIniFileName, name)
}

func GetNameWithoutSuffix(name string) string {
	return fmt.Sprintf("%s%s", PreIniFileName, name)
}
