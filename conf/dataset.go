package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"

	"github.com/simplemoon/deploy/utils"
)

const (
	KeyServerId = "serverId"
)

var (
	baseData = make(DataMap)
)

// 数据的 MAP
type DataMap map[string]interface{}

// 所有需要的数据集合
type DataSet struct {
	Base    DataMap // 基础的数据，配置的数据
	Special DataMap // 实际的数据
	portMap IntMap  // 端口信息
}

// 获取对应的数据
func (ds *DataSet) Get(name string) interface{} {
	if ret, ok := ds.Special[name]; ok {
		return ret
	}
	if ret2, ok2 := ds.Base[name]; ok2 {
		return ret2
	}
	if ret3, ok3 := ds.portMap[name]; ok3 {
		return ret3
	}
	return nil
}

// 获取string类型的数据
func (ds *DataSet) GetString(name string) (string, error) {
	data := ds.Get(name)
	if data == nil {
		return "", nil
	}
	if v, ok := data.(string); ok {
		return v, nil
	}
	return "", fmt.Errorf("%v convert to string failed", data)
}

// 获取string类型的数据
func (ds *DataSet) GetInt(name string) (int, error) {
	data := ds.Get(name)
	if data == nil {
		return 0, nil
	}
	if v, ok := data.(int); ok {
		return v, nil
	}
	return 0, fmt.Errorf("%v convert to string failed", data)
}

// 获取内容
func (ds *DataSet) GetValue(key string) (string, error) {
	if key[0] != FlagQuery {
		return key, nil
	}
	// 去掉之前的字符
	key = key[1:]
	// 获取对应的数值
	v := ds.Get(key)
	if v == nil {
		return "", fmt.Errorf(`%s not exist in dataset`, key)
	}
	if s, ok := v.(string); ok {
		return s, nil
	}
	return "", fmt.Errorf(`%s must be string`, key)
}

// 获取序号
func (ds *DataSet) GetIndex() int {
	serverId := ds.GetServerId()
	idx, err := GetIndexByServerId(serverId)
	if err != nil {
		return 0
	}
	return idx
}

// 获取序列号
func (ds *DataSet) GetServerId() int {
	val := ds.Get(KeyServerId)
	switch t := val.(type) {
	case int:
		return t
	case string:
		ret, err := strconv.Atoi(t)
		if err != nil {
			return 0
		}
		return ret
	default:
		return 0
	}
}

// 设置map的数值
func (ds *DataSet) SetMap(pm IntMap) {
	if len(ds.portMap) == 0 {
		ds.portMap = make(map[string]int, len(pm))
	}
	for k, v := range pm {
		ds.portMap[k] = v
	}
}

// 是否有服务器
func (ds *DataSet) HasServer() bool {
	t := ds.Get(InfoKeyGameType)
	if t == nil {
		mc := ds.Get(InfoKeyMemberCount)
		if v, ok := mc.(int); ok && v > 0 {
			return true
		}
	} else {
		if v, ok := t.(string); ok && v != GSTypeCross {
			return true
		}
	}
	return false
}

// 是否有平台
func (ds *DataSet) HasPlatform() bool {
	t := ds.Get(InfoKeyGameType)
	if t == nil {
		isp := ds.Get(InfoKeyIsPlatform)
		if v, ok := isp.(bool); ok && v {
			return true
		}
	} else {
		if v, ok := t.(string); ok && v == GSTypeGame {
			return false
		}
		ct := ds.Get(InfoKeyCrossType)
		if ct == nil {
			return true
		}
		if v, ok := ct.(string); ok && v != CSTypeRoom {
			return true
		}
	}
	return false
}

// 是否有ROOM
func (ds *DataSet) HasRoom() bool {
	t := ds.Get(InfoKeyGameType)
	if t != nil {
		if v, ok := t.(string); ok && v == GSTypeGame {
			return false
		}
		ct := ds.Get(InfoKeyCrossType)
		if ct == nil {
			return true
		}
		if v, ok := ct.(string); ok && v != CSTypePlatform {
			return true
		}
	}
	return false
}

// 是否有计费
func (ds *DataSet) HasCharge() bool {
	if !ds.HasServer() {
		return false
	}
	ic := ds.Get(InfoKeyIsCharge)
	if ic == nil {
		// 默认开启
		return true
	}
	if v, ok := ic.(bool); ok && v {
		return true
	}
	return false
}

// 获取上传的url
func (ds *DataSet) GetBaseUrl() string {
	v, err := ds.GetString(InfoKeyDeployUrl)
	if err == nil && v != "" {
		return v
	}
	v, err = ds.GetString(InfoKeyChargeUrl)
	if err != nil {
		return ""
	}
	return v
}

// 加载基础的配置
func loadBase() error {
	baseJson, err := utils.GetCfgPath(FileNameBase)
	if err != nil {
		return err
	}
	if !utils.IsPathExist(baseJson) {
		return fmt.Errorf("%s not exist", baseJson)
	}
	// 加载数据
	data, err := ioutil.ReadFile(baseJson)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &baseData)
	if err != nil {
		return err
	}
	return nil
}

// 获取data set
func GetDataSet(idx int) *DataSet {
	// 获取 dataSet
	spec := serverCfg[idx]
	return &DataSet{
		Base:    merge(spec),
		Special: spec,
	}
}

// 合并对应的数据
func merge(dm DataMap) DataMap {
	reg := regexp.MustCompile(`\$\{([^\{\}\s]+)\}`)
	// 创建数据
	data := make(DataMap)
	for key, val := range baseData {
		switch val.(type) {
		case string:
			// 匹配字符串
			result := reg.ReplaceAllStringFunc(val.(string), func(src string) string {
				// 获取对应的关键字
				src = src[2 : len(src)-1]
				// 获取对应的数据
				vd, ok := dm[src]
				if !ok {
					return src
				}
				return fmt.Sprintf("%v", vd)
			})
			data[key] = result
		default:
			data[key] = val
		}
	}
	return data
}
