package cmds

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/simplemoon/deploy/report"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/simplemoon/deploy/conf"
	"github.com/simplemoon/deploy/utils"
)

// extra 的文件结构
type keyMap = map[string]string

// 更新服务器进程配置
type CmdInit struct {
	*CommandBase
	ports  conf.CfgPort      // 端口的配置
	upload map[string]string // 上传的端口和参数信息
}

// 创建一个init
func NewCmdInit() *CmdInit {
	return &CmdInit{
		CommandBase: NewCommandBase(CmdNameInit),
		upload:      make(map[string]string),
	}
}

// 获取一个新的类型
func (cmd *CmdInit) Copy() ICommand {
	return NewCmdInit()
}

// 运行
func (cmd *CmdInit) Run(ds *conf.DataSet) (ret []report.RowInfo, err error) {
	ret = cmd.result
	// 获取 idx
	serverId := ds.GetServerId()
	// 配置 extra 之中的日至配置的信息
	idx, err := conf.GetIndexByServerId(serverId)
	if err != nil {
		return
	}
	// 1. load configs
	if err = cmd.loadPortCfg(); err != nil {
		return
	}
	// 2. 配置 extra 目录的信息
	if err = cmd.configExtra(ds, idx); err != nil {
		return
	}
	// 获取进程信息
	pl := conf.GetProcInfo(ds, conf.ProcTypeContainsAll)
	if len(pl) == 0 {
		return ret, ErrNotHaveProcess
	}
	// 3. 生成配置的端口信息
	if err = cmd.configUsePort(pl, ds, idx); err != nil {
		return
	}
	// 4. 报告端口信息
	if err = cmd.reportInfo(ds); err != nil {
		return
	}
	// 5. 发送端口信息
	if err = cmd.sendReport(ds.GetBaseUrl()); err != nil {
		return
	}
	// 6. 生成 service 的文件信息
	if err = cmd.configService(pl, idx); err != nil {
		return
	}
	// 7. 安装所有的服务
	if err = cmd.installService(ds); err != nil {
		return
	}

	return
}

// 加载端口号配置
func (cmd *CmdInit) loadPortCfg() error {
	cmd.SetStep("load port")
	// 获取文件名称
	fp, err := utils.GetCfgPath(conf.FileNamePort)
	if err != nil {
		return err
	}
	content, err := ioutil.ReadFile(fp)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(content, &cmd.ports); err != nil {
		return err
	}
	return nil
}

// 配置 extra 内容
func (cmd *CmdInit) configExtra(ds *conf.DataSet, idx int) error {
	cmd.SetStep("config extra")
	// 服务器的日志配置文件目录
	tmp, err := utils.GetCfgPath(conf.FileNameExtra)
	if err != nil {
		return err
	}
	// 读取文件内容
	content, err := ioutil.ReadFile(tmp)
	if err != nil {
		return err
	}
	// 解析的数据
	var data map[string]keyMap
	if err = json.Unmarshal(content, &data); err != nil {
		return err
	}
	// 执行目录
	closes := make([]*os.File, 0)
	defer func() {
		for _, f := range closes {
			err := f.Close()
			if err != nil {
				cmd.GetLogger().Printf(ErrTextCloseFailed, f.Name(), err)
			} else {
				cmd.GetLogger().Printf(TextNameCloseSuccess, f.Name())
			}
		}
	}()

	serverAt := utils.GetServerDir(idx)
	for f, d := range data {
		// 打开数据，并且替换
		p := path.Join(serverAt, f)
		file, err := os.OpenFile(p, os.O_RDWR, 0666)
		if err == nil {
			cmd.GetLogger().Printf(ErrTextOpenFailed, p, err)
			return err
		}
		// 加入关闭列表
		closes = append(closes, file)
		// 读取内容，并且替换
		c, err := ioutil.ReadAll(file)
		if err != nil {
			cmd.GetLogger().Printf(ErrTextReadFailed, p, err)
			return err
		}
		cs := string(c)
		// 替换内容
		for k, v := range d {
			rv, err := ds.GetValue(v)
			if err != nil {
				return err
			}
			strings.Replace(cs, k, rv, -1)
		}
		_, err = file.Seek(0, 0)
		if err != nil {
			cmd.GetLogger().Printf(ErrTextSeekFailed, p, err)
			return err
		}
		_, err = file.WriteString(cs)
		if err != nil {
			cmd.GetLogger().Printf(ErrTextWriteFailed, p, err)
			return err
		}
	}
	return nil
}

// 配置端口信息
func (cmd *CmdInit) configUsePort(pl conf.PIList, ds *conf.DataSet, idx int) error {
	cmd.SetStep("config use port")
	// 生成端口号
	pm := make(conf.IntMap)
	for _, p := range pl {
		name := p.GetBaseName()
		data := cmd.ports.GetPortInfo(name)
		if data == nil {
			continue
		}
		// 获取基础的端口号
		i := p.GetIndex()
		cnt := len(data)
		for key, port := range data {
			bp, isReplace := cmd.ports.GetBasePort(name, key, idx)
			if isReplace {
				port = bp + cnt*i + port
			} else {
				port = bp + i
			}
			pk := fmt.Sprintf(conf.ProcessKeyPortFmt, name, key, i)
			pm[pk] = port
			cmd.GetLogger().Printf(LogNameSetPort, pk, port)
		}
	}
	// 设置port的配置信息
	ds.SetMap(pm)
	return nil
}

// 配置 ini 文件
func (cmd *CmdInit) configIniFile(pl conf.PIList, ds *conf.DataSet, idx int) error {
	cmd.SetStep("config ini files")
	// 获取所有需要生成的配置的名称
	procMap := make(map[string]bool, len(pl))
	for _, p := range pl {
		procMap[p.GetName()] = true
	}
	// 加载数据
	svr := make(conf.CfgSvr)
	name, err := utils.GetCfgPath(conf.FileNameProcIni)
	if err != nil {
		return err
	}
	content, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(content, &svr); err != nil {
		return err
	}
	// 对应的目录
	svrDir := utils.GetServerDir(idx)
	srcDir, err := utils.GetCfgIniDir()
	if err != nil {
		return err
	}
	// 遍历一下数据
	for p, d := range svr {
		// 拷贝对应的文件
		if !conf.NeedCopy(p, procMap) {
			cmd.GetLogger().Printf("copy file.ignore %s", p)
			continue
		}
		// 获取目录和文件名
		pos := strings.LastIndex(p, conf.FilterString)
		if pos == -1 {
			cmd.GetLogger().Printf("%s not found %s. ignore", p, conf.FilterString)
			continue
		}
		// 目录名和文件名
		fd := p[:pos]
		fn := p[pos+1:]
		// 目标的目录
		targetDir := path.Join(svrDir, fd)
		err := os.MkdirAll(targetDir, 0666)
		if err != nil {
			cmd.GetLogger().Printf("create %s failed", targetDir)
			return err
		}
		// 源文件
		sf := path.Join(srcDir, fn)
		// 目标文件
		tf := path.Join(svrDir, p)

		// 获取对应的目录
		if len(d) == 0 {
			// 拷贝文件就行了
			err := conf.CopyFile(sf, tf)
			if err != nil {
				return err
			}
			continue
		}
		// 替换对应的内容
		err = conf.ReplaceContent(sf, tf, ds, d)
		if err != nil {
			return err
		}
	}
	return nil
}

// 报告端口信息
func (cmd *CmdInit) reportInfo(ds *conf.DataSet) error {
	cmd.SetStep("report info")
	reportPath, err := utils.GetCfgPath(conf.FileNameReport)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadFile(reportPath)
	if err != nil {
		return err
	}
	items := make([]conf.ReportItem, 0)
	err = json.Unmarshal(data, &items)
	if err != nil {
		return err
	}
	for _, r := range items {
		if r.IsNil() {
			continue
		}
		// 如果需要报告多个
		if r.Count != "" {
			val := ds.Get(r.Count)
			cnt, ok := val.(int)
			if !ok {
				return fmt.Errorf("key: %s, value: %v convert to int failed", r.Count, val)
			}
			if cnt <= 0 {
				cmd.GetLogger().Printf("key: %s value is zero, don't report", r.Count)
				continue
			}
			for i := 0; i < cnt; i++ {
				k := fmt.Sprintf(r.Key, i+1)
				err := cmd.addReport(ds, k, r.ReportKey)
				if err != nil {
					return err
				}
			}
		} else {
			if err := cmd.addReport(ds, r.Key, r.ReportKey); err != nil {
				return err
			}
		}
	}
	return nil
}

// 发送信息
func (cmd *CmdInit) sendReport(url string) error {
	if url == "" {
		return fmt.Errorf("report server port failed url is empty")
	}
	// 发送数据
	content, err := json.Marshal(cmd.upload)
	if err != nil {
		return err
	}
	// 发送对应的数据
	url += conf.UrlPathUpdate
	r, err := http.Post(url, "application/json", bytes.NewReader(content))
	if err != nil {
		return err
	}
	if r.StatusCode != 200 {
		return fmt.Errorf("send %v to %s failed, code: %d", content, url, r.StatusCode)
	}
	return nil
}

// 加入数值
func (cmd *CmdInit) addReport(ds *conf.DataSet, key, reportKey string) error {
	val, err := ds.GetValue(key)
	if err != nil {
		return err
	}
	content := cmd.upload[reportKey]
	if content == "" {
		cmd.upload[reportKey] = fmt.Sprintf("%v", val)
	} else {
		cmd.upload[reportKey] = fmt.Sprintf("%s,%v", content, val)
	}
	return nil
}

// 配置端口信息
func (cmd *CmdInit) configService(pl conf.PIList, idx int) error {
	cmd.SetStep("config service files")
	// 服务器的bin目录
	svrBinPath := utils.GetServerBinDir(idx)
	launchPath := path.Join(svrBinPath, conf.FileBin64LaunchExe)
	// service的目录
	servicePath := utils.GetServiceDir(idx)
	if err := utils.CreateDir(servicePath); err != nil {
		return err
	}
	// exe 文件的路径
	serviceExePath, err := utils.GetCfgPath(conf.FileNameServiceExe)
	if err != nil {
		return err
	}
	// 读取文件
	serviceCfgPath, err := utils.GetCfgPath(conf.FileNameServiceXml)
	if err != nil {
		return err
	}
	content, err := ioutil.ReadFile(servicePath)
	if err != nil {
		return err
	}

	// 配置对应的 service 信息
	for _, p := range pl {
		// 复制文件内容
		c := string(content)
		// 替换文件内容
		fileName := p.GetFileName()
		// 替换对应的内容
		serviceName := p.GetServiceName(idx)
		c = strings.ReplaceAll(c, "${service_id}", p.GetServiceId(idx))
		c = strings.ReplaceAll(c, "${service_name}", serviceName)
		c = strings.ReplaceAll(c, "${target}", fileName)
		c = strings.ReplaceAll(c, "${ini_target}", fileName)
		// 上述的都替换完成了，那么复制文件进去
		launchTarget := path.Join(svrBinPath, fmt.Sprintf("%s%s", fileName, utils.ServerSuffixExe))
		serviceXmlTarget := path.Join(servicePath, fmt.Sprintf("%s%s", fileName, utils.ServiceSuffixXml))
		serviceExeTarget := path.Join(servicePath, fmt.Sprintf("%s%s", fileName, utils.ServiceSuffixExe))
		// 拷贝文件
		err := ioutil.WriteFile(serviceXmlTarget, content, 0666)
		if err != nil {
			cmd.GetLogger().Printf("write to %s err %v", serviceXmlTarget, err)
			return err
		}
		// 拷贝文件
		if err = conf.CopyFile(launchPath, launchTarget); err != nil {
			cmd.GetLogger().Printf("copy %s to %s err %v", launchPath, launchTarget, err)
			return err
		}
		if err = conf.CopyFile(serviceExePath, serviceExeTarget); err != nil {
			cmd.GetLogger().Printf("copy %s to %s err %v", serviceCfgPath, serviceExeTarget, err)
			return err
		}
		cmd.GetLogger().Printf("create %s service config success")

		// 添加结果信息
		info := report.RowInfo{
			ServerIdx:  idx,
			Action:     cmd.GetName(),
			ActionType: cmd.GetName(),
			Msg:        fmt.Sprintf("create service %s success", serviceName),
		}
		cmd.AddResult(&info)
	}
	return nil
}

// 安装所有的服务
func (cmd *CmdInit) installService(ds *conf.DataSet) error {
	cmd.SetStep("install service")
	var wg sync.WaitGroup
	// 执行对应的结果
	err := cmd.ExecServiceCmd(ds, &wg, ServiceCmdNameInstall, conf.ProcTypeContainsAll)
	if err != nil {
		return err
	}
	wg.Wait()

	return nil
}

func init() {
	Register(NewCmdInit())
}
