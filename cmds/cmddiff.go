package cmds

import (
	"encoding/json"
	"fmt"
	"github.com/simplemoon/deploy/report"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"

	"github.com/simplemoon/deploy/conf"
	"github.com/simplemoon/deploy/utils"
)

// 检查文件的配置
type CheckItem struct {
	conf.SectionMap        // 具体的数据
	ResetCnt        string `json:"ResetCnt"`
}

type CmdDiff struct {
	*CommandBase
	idx int // 服务器编号
}

func NewCmdDiff() *CmdDiff {
	return &CmdDiff{
		CommandBase: NewCommandBase(CmdNameDiff),
	}
}

func (cmd *CmdDiff) Copy() ICommand {
	return NewCmdDiff()
}

func (cmd *CmdDiff) Run(ds *conf.DataSet) (ret []report.RowInfo, err error) {
	ret = cmd.result
	// 对比游戏服务器的配置和后台的配置数值
	if err = cmd.check(ds); err != nil {
		return
	}
	return
}

func (cmd *CmdDiff) check(ds *conf.DataSet) error {
	idx := ds.GetIndex()
	if idx <= 0 {
		return fmt.Errorf("server not found")
	}
	cmd.idx = idx

	// 获取名称，并且检查
	cp, err := utils.GetCfgPath(conf.FileNameCheck)
	if err != nil {
		return err
	}
	content, err := ioutil.ReadFile(cp)
	if err != nil {
		return err
	}
	d := make(map[string]CheckItem)
	err = json.Unmarshal(content, &d)
	if err != nil {
		return err
	}
	// 读取配置并且解析
	var wg sync.WaitGroup
	rch := make(chan *report.RowInfo, 32)

	for name, item := range d {
		cnt, err := cmd.getCnt(ds, item.ResetCnt)
		if err != nil {
			return err
		}
		if cnt == 0 {
			wg.Add(1)
			p := conf.NewPI(name, 0)
			go cmd.diff(&wg, rch, ds, item.SectionMap, p)
		} else {
			for i := 0; i < cnt; i++ {
				wg.Add(1)
				p := conf.NewPI(name, i+1)
				go cmd.diff(&wg, rch, ds, item.SectionMap, p)
			}
		}
	}
	// 等待关闭的进程
	go func() {
		wg.Wait()
		close(rch)
	}()
	for v := range rch {
		cmd.result = append(cmd.result, *v)
	}
	return nil
}

func (cmd *CmdDiff) getCnt(ds *conf.DataSet, key string) (int, error) {
	if key == "" {
		return 0, nil
	}
	if key[0] == conf.FlagQuery {
		key = key[1:]
	}
	// 获取数据
	val := ds.Get(key)
	if val == nil {
		return 0, fmt.Errorf("%s not found at dataset", key)
	}
	switch v := val.(type) {
	case string:
		if r, err := strconv.Atoi(v); err == nil {
			return r, nil
		} else {
			return 0, err
		}
	case int:
		return v, nil
	}
	return 0, fmt.Errorf("%s type error. must be int or string to int", key)
}

func (cmd *CmdDiff) diff(wg *sync.WaitGroup, r chan<- *report.RowInfo,
	ds *conf.DataSet, m conf.SectionMap, p *conf.PI) {
	// 关闭
	defer wg.Done()
	fn := p.GetFileNameWithSuffix()
	fp := utils.GetServerIniConfigPath(fn, cmd.idx)
	// 结果信息
	ret := &report.RowInfo{
		ServerIdx:  cmd.idx,
		Action:     cmd.GetName(),
		ActionType: p.GetBaseName(),
		State:      report.Failed,
	}

	// 读取文件内容
	cfg, err := ini.Load(fp)
	if err != nil {
		ret.Msg = err.Error()
		r <- ret
		return
	}
	// 读取并且解析对于对应的文件
	for sn, data := range m {
		s, err := cfg.GetSection(sn)
		if err != nil {
			ret.Msg = err.Error()
			r <- ret
			return
		}
		// 遍历key 对应的 数值
		for key, val := range data {
			k, err := s.GetKey(key)
			if err != nil {
				ret.Msg = err.Error()
				r <- ret
				return
			}
			// 获取对应的key的数值
			v, err := cmd.getValue(ds, val, p.GetIndex())
			if err != nil {
				ret.Msg = fmt.Sprintf("get %s from dataset error", val)
				r <- ret
				return
			}

			if k.String() != v {
				ret.Msg = fmt.Sprintf("section: %s, key: %s, v1: %s, v2: %s don't same", sn, key, k.String(), v)
				r <- ret
				return
			}
		}
	}
	ret.State = report.Success
	ret.Msg = fmt.Sprintf("file %s verify ok", p.GetFileName())
	r <- ret
}

func (cmd *CmdDiff) getValue(ds *conf.DataSet, name string, idx int) (string, error) {
	if name == "" {
		return "", nil
	}
	switch name[0] {
	case conf.FlagQuery:
		// 获取数值
		return ds.GetString(name[1:])
	case conf.FlagSlice:
		// 获取数值
		v, err := ds.GetString(name[1:])
		if err != nil {
			return "", err
		}
		vl := strings.Split(v, ",")
		cnt := len(vl)
		if cnt <= idx-1 {
			return "", fmt.Errorf("key: %s has no index of %d", name, idx)
		}
		return vl[idx-1], err
	default:
		return ds.GetString(name)
	}
}

func init() {
	Register(NewCmdDiff())
}
