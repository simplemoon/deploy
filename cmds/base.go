package cmds

import (
	"fmt"
	"github.com/simplemoon/deploy/conf"
	"github.com/simplemoon/deploy/report"
	"github.com/simplemoon/deploy/utils"
	"github.com/sirupsen/logrus"
)

var (
	// 所有命令的集合类
	allCommand = make([]ICommand, 20)
)

// 步骤设置
type IStep interface {
	SetStep(name string)
	GetStep() string
}

// 日志接口
type ILogger interface {
	SetLogger(logger *logrus.Logger) // 设置日志实例
	GetLogger() *logrus.Logger       // 获取日志实例
}

// 链表接口
type ICmdList interface {
	GetNext() ICommand  // 获取下一个命令
	SetNext(c ICommand) // 设置下一个
}

// 命令接口
type ICommandBase interface {
	GetName() string                                // 名称
	Prepare(ds *conf.DataSet) error                 // 检查的一些工作
	Check(ds *conf.DataSet) error                   // 检查能否运行
	Run(ds *conf.DataSet) ([]report.RowInfo, error) // 执行命令
}

// 命令
type ICommand interface {
	IStep        // 步骤设置
	ILogger      // 日志接口
	ICmdList     // 设置链表接口
	ICommandBase // 命令接口
}

// 状态机
type StepMachine struct {
	step string // 状态的名称
}

func (s StepMachine) SetStep(name string) {
	s.step = name
}

func (s StepMachine) GetStep() string {
	return s.step
}

// 基础的命令
type CommandBase struct {
	StepMachine          // 步骤名称
	Name        string   // 命令的名称
	Next        ICommand // 下一个命令

	logger *logrus.Logger // 日志实例
}

func NewCommandBase(name string) CommandBase {
	if name == "" {
		name = CmdNameNone
	}
	return CommandBase{Name: name, Next: nil}
}

// 获取命令的名称
func (bc CommandBase) GetName() string {
	return CmdNameNone
}

// 准备
func (bc CommandBase) Prepare(ds *conf.DataSet) error {
	bc.SetStep("prepare")
	return nil
}

// 检查能否运行
func (bc CommandBase) Check(ds *conf.DataSet) error {
	bc.SetStep("check")
	// 检查配置文件是否存在
	serverId := ds.GetServerId()
	idx, err := conf.GetIndexByServerId(serverId)
	if err != nil {
		return err
	}
	// 检查文件夹是否存在
	if idx <= 0 {
		return fmt.Errorf("get %d index failed", serverId)
	}
	// 检查server文件夹是否存在
	if !utils.CheckServerDirExist(idx) {
		return fmt.Errorf("server%d path not found", serverId)
	}
	return nil
}

// 运行
func (bc CommandBase) Run(ds *conf.DataSet) ([]report.RowInfo, error) {
	return nil, nil
}

// 获取下一个
func (bc CommandBase) GetNext() ICommand {
	return bc.Next
}

// 设置下一个命令
func (bc CommandBase) SetNext(c ICommand) {
	bc.Next = c
}

// 设置实例
func (bc CommandBase) SetLogger(logger *logrus.Logger) {
	bc.logger = logger
}

// 获取实例
func (bc CommandBase) GetLogger() *logrus.Logger {
	return bc.logger
}
