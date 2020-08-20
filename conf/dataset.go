package conf

import (
	"encoding/json"
	"fmt"
	"github.com/simplemoon/deploy/utils"
	"io/ioutil"
	"regexp"
	"strconv"
)

const (
	KeyServerId = "serverId"
)

var (
	baseData = DataMap(make(map[string]interface{}))
)

// 数据的 MAP
type DataMap map[string]interface{}

// 所有需要的数据集合
type DataSet struct {
	Base    DataMap // 基础的数据，配置的数据
	Special DataMap // 实际的数据
}

// 获取对应的数据
func (ds *DataSet) Get(name string) interface{} {
	if ret, ok := ds.Special[name]; ok {
		return ret
	}
	if ret2, ok2 := ds.Base[name]; ok2 {
		return ret2
	}
	return nil
}

// 获取序号
func (ds *DataSet) GetIndex() int {
	// TODO: GET INDEX FROM DATASET
	return 0
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
