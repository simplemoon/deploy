package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/simplemoon/deploy/utils"
)

type ServerInfo struct {
	Idx      int `json:"_idx"`
	ServerId int `json:"serverId"`
}

func Check(serverId int) error {
	_, err := GetIndexByServerId(serverId)
	if err != nil {
		return err
	}
	return nil
}

func GetIndexByServerId(serverId int) (int, error) {
	filePath := utils.GetRuntimePath(FileRuntimeServers)
	if !utils.IsPathExist(filePath) {
		return 0, fmt.Errorf("path %s not exist", filePath)
	}
	// 加载对应的程序
	servers := make([]ServerInfo, 0)
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return 0, err
	}
	err = json.Unmarshal(data, &servers)
	if err != nil {
		return 0, err
	}
	// 找一下对应的数据
	for i := 0; i < len(servers); i++ {
		if servers[i].ServerId == serverId {
			return servers[i].Idx, nil
		}
	}
	return 0, fmt.Errorf("do not find %d config", serverId)
}

// 创建对应的index
func CreateServerIdx(serverId int, isDelete bool) (int, error) {
	// 加载对应的程序
	servers := make([]ServerInfo, 0)
	filePath := utils.GetRuntimePath(FileRuntimeServers)
	if utils.IsPathExist(filePath) {
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			return 0, err
		}
		err = json.Unmarshal(data, &servers)
		if err != nil {
			return 0, err
		}
	}
	// 找一下对应的数据
	cnt := len(servers)
	ns := servers[:0]
	idxMap := make(map[int]bool)
	for i := 0; i < cnt; i++ {
		if servers[i].ServerId == serverId {
			if !isDelete {
				return servers[i].Idx, nil
			} else {
				continue
			}
		}
		ns = append(ns, servers[i])
		idxMap[servers[i].Idx] = true
	}
	// 创建一个index
	newIdx := cnt
	// 增加数据
	for i := 0; i < cnt; i++ {
		if v, ok := idxMap[i+1]; ok && v {
			continue
		}
		newIdx = i + 1
		break
	}
	ns = append(ns, ServerInfo{
		ServerId: serverId,
		Idx:      newIdx,
	})
	// 写入对应的数据
	data, err := json.Marshal(ns)
	if err != nil {
		return newIdx, err
	}
	err = ioutil.WriteFile(filePath, data, 0666)
	if err != nil {
		return newIdx, err
	}
	return newIdx, nil
}
