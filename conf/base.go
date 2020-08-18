package conf

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/simplemoon/deploy/log"
	"github.com/simplemoon/deploy/utils"
)

var (
	// 命令行参数
	cmdArgs = CommandParams{}
	// 待执行的服务器的配置信息
	serverCfg = make([]map[string]interface{}, 0)
)

// 命令行参数
type CommandParams struct {
	action      string // 保存对应的命令列表
	jsonContent string // (-j) json 内容
	urlKey      string // (-k) 请求参数的加密内容
	filePath    string // (-f) json文件的路径

	actionList   []string // (-a) 需要执行的操作
	rootPath     string   // (-r) 目录的盘符
	jarPath      string   // (-jar) 合服的jar包路径
	timeStamp    string   // (-t) 期望设置的服务器时间
	telnetTarget string   // (-s) 期望发送到的服务器进程名称
	telnetParams string   // (-p) 执行命令需要的参数
	onceCount    int      // (-count) 每次执行的数量
	isDebug      bool     // (-d) 是否是需要打印debug日志
	isCreatCfg   bool     // (-c) 是否需要生成配置文件
	isRunOneKey  bool     // (-b) 是否是一键执行
}

// 初始化，添加参数处理逻辑
func init() {
	flag.StringVar(&cmdArgs.action, "a", "", "action list to exec")
	flag.StringVar(&cmdArgs.rootPath, "r", "", "dir path you want game server install")
	flag.StringVar(&cmdArgs.jsonContent, "j", "", "json content base64 encode")
	flag.StringVar(&cmdArgs.filePath, "f", "", "the file path can replace json content")
	flag.StringVar(&cmdArgs.urlKey, "k", "", "the url path encode by base64")
	flag.StringVar(&cmdArgs.jarPath, "jar", "", "merge jar path")
	flag.StringVar(&cmdArgs.timeStamp, "t", "", "the time you want to change to")
	flag.StringVar(&cmdArgs.telnetParams, "p", "", "the params of telnet")
	flag.IntVar(&cmdArgs.onceCount, "count", 3, "every time you want to run the server instance")
	flag.BoolVar(&cmdArgs.isDebug, "d", true, "if you need print logs, set it")
	flag.BoolVar(&cmdArgs.isCreatCfg, "c", false, "to create server config")
	flag.BoolVar(&cmdArgs.isRunOneKey, "b", false, "is it report result by https")
	// 解析参数
	flag.Parse()
}

// 参数太长的请求处理
func requestIfTooLong() bool {
	if cmdArgs.urlKey == "" {
		return true
	}
	// base64 解密数据
	data, err := base64.StdEncoding.DecodeString(cmdArgs.urlKey)
	if err != nil {
		log.FormatErr(err)
		return false
	}
	// 获取对应的内容
	resp, err := http.Get(string(data))
	if err != nil {
		log.FormatErr(err)
		return false
	}
	defer resp.Body.Close()
	// 获取得到的内容
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.FormatErr(err)
		return false
	}
	// 解析对应的参数
	result := strings.Fields(string(content))
	err = flag.CommandLine.Parse(result[2:])
	if err != nil {
		log.FormatErr(err)
		return false
	}
	return true
}

// 获取content
func getContent() ([]byte, error) {
	// 检查文件路径
	if cmdArgs.filePath != "" && utils.IsPathExist(cmdArgs.filePath) {
		data, err := ioutil.ReadFile(cmdArgs.filePath)
		if err == nil {
			return data, nil
		}
	}
	// 获取 content 内容
	return base64.StdEncoding.DecodeString(cmdArgs.jsonContent)
}

// 准备参数
func PrePare() bool {
	if !requestIfTooLong() {
		return false
	}
	// 设置根路径
	utils.SetRootPath(GetRootDir())
	// 命令列表
	cmdArgs.actionList = strings.Fields(cmdArgs.action)
	// 配置信息解析
	content, err := getContent()
	if err != nil {
		log.FormatErr(err)
		return false
	}
	// 获取 content 的内容
	switch content[0] {
	case '[':
		// 说明是一个数组list，单个
		err := json.Unmarshal(content, &serverCfg)
		if err != nil {
			log.FormatErr(err)
			return false
		}
	case '{':
		// 说明是个字典
		result := make(map[string]interface{})
		err := json.Unmarshal(content, &result)
		if err != nil {
			log.FormatErr(err)
			return false
		}
		serverCfg = append(serverCfg, result)
	default:
		// 说明解析错误了
		log.FormatErr("json content error")
		return false
	}
	return true
}

// 获取项目信息
func GetProject() int {
	if len(serverCfg) == 0 {
		return 0
	}
	val := serverCfg[0]["gameId"]
	if val == nil {
		return 0
	}
	if intVal, ok := val.(int); ok {
		return intVal
	}
	return 0
}

// 获取参数
func GetServerUnitCount() int {
	return len(serverCfg)
}

// 获取 action list
func GetActions() []string {
	return cmdArgs.actionList
}

// 获取 rootPath
func GetRootDir() string {
	return cmdArgs.rootPath
}

// 获取 jar 的路径
func GetJarPath() string {
	return cmdArgs.jarPath
}

// 获取服务器的期望时间
func GetTimeStampAt() string {
	return cmdArgs.timeStamp
}

// 获取目标服务器的进程名称
func GetTargetName() string {
	return cmdArgs.telnetTarget
}

// 获取目标的执行参数
func GetTelnetArgs() string {
	return cmdArgs.telnetParams
}

// 协程的个数
func GetRoutineCount() int {
	return cmdArgs.onceCount
}

// 是否是调试模式
func IsDebugModel() bool {
	return cmdArgs.isDebug
}

// 是否需要生成
func IsGenerate() bool {
	return cmdArgs.isCreatCfg
}

// 是否远程模式
func IsRemoteModel() bool {
	return cmdArgs.isRunOneKey
}

// 加载所有数据
func LoadAll() error {
	//err := loadBase()
	//if err != nil {
	//	return err
	//}
	return nil
}
