package conf

import (
	"encoding/json"
	"fmt"
	"github.com/simplemoon/deploy/utils"
	"io/ioutil"
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
	filePath := utils.GetRuntimePath(FileNameServers)
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
