package conf

import (
	"encoding/json"
	"fmt"
	"github.com/simplemoon/deploy/utils"
	"io/ioutil"
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
	Base    *DataMap // 基础的数据，配置的数据
	Special DataMap  // 实际的数据
}

// 获取对应的数据
func (ds *DataSet) Get(name string) interface{} {
	if ret, ok := ds.Special[name]; ok {
		return ret
	}
	if ret2, ok2 := (*ds.Base)[name]; ok2 {
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

func loadBase() error {
	// TODO: load config
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
	return &DataSet{
		Base:    &baseData,
		Special: serverCfg[idx],
	}
}
