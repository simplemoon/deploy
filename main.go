package deploy

import (
	"github.com/simplemoon/deploy/cmds"
	"github.com/simplemoon/deploy/conf"
	"github.com/simplemoon/deploy/log"
	"github.com/simplemoon/deploy/utils"
	"sync"
)

// 同步执行
func syncRun(cl *cmds.CommandList, count int) {
	var wg sync.WaitGroup
	token := make(chan struct{}, conf.GetRoutineCount())
	for i := 0; i < count; i++ {
		token <- struct{}{}
		wg.Add(1)
		// 执行对应的单元功能
		go func(idx int) {
			// 释放token
			defer func() {
				wg.Done()
				<-token
			}()
			// 执行对应的逻辑
			cl.Exec(conf.GetDataSet(idx))
		}(i)
	}
	wg.Wait()
}

// 获取
func main() {
	// 处理输入参数，预备
	if !conf.PrePare() {
		return
	}
	// 加载所有的数据
	err := conf.LoadAll()
	if err != nil {
		log.FormatErr(err)
		return
	}
	// 创建一个执行的队列
	cmdList, err := cmds.NewCommandList(conf.GetActions())
	if err != nil {
		log.FormatErr(err)
		return
	}
	// 创建一些必要的文件夹
	err = utils.CreateDirs(conf.GetRootDir())
	if err != nil {
		log.FormatErr(err)
		return
	}
	// 获取所有的配置
	count := conf.GetServerUnitCount()
	switch {
	// 只有一个执行的单元
	case count == 1:
		cmdList.Exec(conf.GetDataSet(0))
	// 多个执行的单元
	case count > 1:
		syncRun(cmdList, count)
	// 木有执行单元
	default:
		log.FormatErr("do not have exec server config")
		return
	}
}
