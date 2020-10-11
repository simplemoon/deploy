package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/simplemoon/deploy/utils"
)

// 初始化的类型
var (
	// 配置的文件
	procCfg = new(ProcInfoCfg)
	// 进程信息缓存
	procCache = new(sync.Map)
)

// 进程信息
type PIList []*PI

// 进程信息
type PI struct {
	name  string // 名称
	index int    // 序号
	pt    uint32 // 对应的消息
}

func NewPI(name string, idx int) *PI {
	return &PI{
		name:  name,
		index: idx,
	}
}

// 获取基础的名称
func (pi *PI) GetBaseName() string {
	n := strings.TrimPrefix(pi.name, PreIniFileName)
	if pi.index == 0 {
		return n
	}
	return fmt.Sprintf("%s%d", n, pi.index)
}

// 名称
func (pi *PI) GetName() string {
	if pi.index == 0 {
		return pi.name
	} else {
		return fmt.Sprintf("%s%d", pi.name, pi.index)
	}
}

// 序号
func (pi *PI) GetIndex() int {
	return pi.index
}

// 获取服务器名称
func (pi *PI) GetServiceName(idx int) string {
	return fmt.Sprintf("%s-%d-%s", PreIniFileName, idx, pi.GetName())
}

// 获取id
func (pi *PI) GetServiceId(idx int) string {
	return pi.GetServiceName(idx)
}

// 获取文件名称
func (pi *PI) GetFileName() string {
	return fmt.Sprintf("%s%s", PreIniFileName, pi.GetName())
}

// 获取文件名称包含后缀
func (pi *PI) GetFileNameWithSuffix() string {
	return fmt.Sprintf("%s%s%s", PreIniFileName, pi.GetName(), FileSuffixIni)
}

// 服务器进程信息
type GameProcCfg struct {
	Name string `json:"name"`          // 进程的名称
	Key  string `json:"key"`           // 进程个数的关键字
	Max  int    `json:"max,omitempty"` // 最大的数量
}

// 服务器进程信息
type GameProcList []GameProcCfg

// 配置的信息
type ProcInfoCfg struct {
	Base     GameProcList `json:"base"`     // 基础信息，所有的进程都需要
	Server   GameProcList `json:"server"`   // 服务器的信息
	Platform GameProcList `json:"platform"` // 平台服需要的进程
	Charge   GameProcList `json:"charge"`   // 计费进程
	Room     GameProcList `json:"room"`     // 房间服进程
}

// 加载进程配置
func loadProcCfg() error {
	procJson, err := utils.GetCfgPath(FileNameProcess)
	if err != nil {
		return err
	}
	if !utils.IsPathExist(procJson) {
		return fmt.Errorf("%s not exist", procJson)
	}
	// 加载数据
	data, err := ioutil.ReadFile(procJson)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, procCfg)
	if err != nil {
		return err
	}
	return nil
}

// 获取所有的进程信息
func GetProcInfo(ds *DataSet, containsType uint32) PIList {
	serverId := ds.GetServerId()
	val, ok := procCache.Load(serverId)
	if ok {
		return val.(PIList)
	}
	// 获取对应的进程信息了
	r := make(PIList, 0)
	// 添加进程信息
	getPIList(&r, &procCfg.Base, ds, ProcTypeBase, containsType)
	// 检查是否需要server
	if ds.HasServer() {
		getPIList(&r, &procCfg.Server, ds, ProcTypeServer, containsType)
	}
	if ds.HasCharge() {
		getPIList(&r, &procCfg.Charge, ds, ProcTypeCharge, containsType)
	}
	if ds.HasPlatform() {
		getPIList(&r, &procCfg.Platform, ds, ProcTypePlatform, containsType)
	}
	if ds.HasRoom() {
		getPIList(&r, &procCfg.Room, ds, ProcTypeRoom, containsType)
	}
	return r
}

// 获取信息列表
func getPIList(r *PIList, gl *GameProcList, ds *DataSet, pt, contains uint32) {
	if contains&pt == 0 {
		return
	}
	for _, v := range *gl {
		if v.Key != "" && v.Key != "null" {
			count := 1
			val := ds.Get(v.Key)
			if c, ok := val.(int); ok {
				count = c
			} else {
				count = 1
			}
			// 限制最大数量，防止配置错误
			if v.Max != 0 && count > v.Max {
				count = v.Max
			}
			for j := 0; j < count; j++ {
				*r = append(*r, &PI{name: v.Name, index: j + 1, pt: pt})
			}
		} else {
			*r = append(*r, &PI{name: v.Name, index: 0, pt: pt})
		}
	}
}
