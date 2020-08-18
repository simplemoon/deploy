package log

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/simplemoon/deploy/report"
	"github.com/simplemoon/deploy/utils"
	"github.com/sirupsen/logrus"
)

var (
	// 缓存的日志实例
	cache = Caches{
		loggers: make(map[int]*logrus.Logger),
	}
	// 默认的格式
	defaultFormatter = &logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
		DisableQuote:  true,
	}
)

// 日志缓存
type Caches struct {
	mutex   sync.Mutex             // 锁文件
	loggers map[int]*logrus.Logger // 日志记录实例
}

// 获取日志实例
func (c *Caches) GetLogger(serverId int, isDebug bool) *logrus.Logger {
	// 获取 logger
	l := c.loggers[serverId]
	if l != nil {
		return l
	}
	// 加入锁文件
	c.mutex.Lock()
	defer c.mutex.Unlock()
	l = c.loggers[serverId]
	if l != nil {
		return l
	}
	// 创建一个实例
	filePath := utils.GetLogPath(serverId)
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil
	}
	logger := logrus.New()
	logger.SetOutput(f)
	logger.SetFormatter(defaultFormatter)
	if isDebug {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}
	// 记录一下
	c.loggers[serverId] = logger
	return logger
}

// 创建一个logger
func CreateLogger(serverId int, isDebug bool) *logrus.Logger {
	return cache.GetLogger(serverId, isDebug)
}

// 报告错误
func FormatErr(v interface{}) {
	msg := fmt.Sprintf("%v", v)
	msg = strings.Replace(msg, "'", "", -1)
	// 创建一个结果
	r := report.NewErrReport(msg)
	// 反序列化成字符串
	info, err := json.Marshal(r)
	if err != nil {
		fmt.Println("json marshal failed ", err)
		return
	}
	fmt.Fprintf(os.Stderr, "%s", info)
}
