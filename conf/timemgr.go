package conf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/ini.v1"
	"net/http"
	"strconv"
	"time"
)

type TimeMgr struct {
	path     string // 文件的路径
	serverId int    // 服务器ID
	gameId   string // 游戏ID
	url      string // 基础的信息
}

func NewTimeMgr(p, gameId string, serverId int, url string) *TimeMgr {
	return &TimeMgr{
		path:     p,
		gameId:   gameId,
		serverId: serverId,
		url:      url,
	}
}

// 获取时间配置
func (tm *TimeMgr) GetTime() (int, error) {
	cfg, err := ini.Load(tm.path)
	if err != nil {
		return 0, err
	}
	s, err := cfg.GetSection(IniKeySectionMain)
	if err != nil {
		return 0, err
	}
	if !s.HasKey(IniKeyTimeAhead) {
		return 0, nil
	}
	key, err := s.GetKey(IniKeyTimeAhead)
	if err != nil {
		return 0, err
	}
	v, err := key.Int()
	if err != nil {
		return 0, err
	}
	return v, nil
}

// 保存对应的时间配置
func (tm *TimeMgr) SaveTime(seconds int) error {
	cfg, err := ini.Load(tm.path)
	if err != nil {
		return err
	}
	s := cfg.Section(IniKeySectionMain)
	s.Key(IniKeyTimeAhead).SetValue(strconv.Itoa(seconds))
	// 保存到文件之中
	if err = cfg.SaveToIndent(tm.path, "\r\n"); err != nil {
		return err
	}
	return tm.notify(seconds)
}

func (tm *TimeMgr) notify(t int) error {
	// 上传信息到web上面
	content := struct {
		Offset   int    `json:"offset"`
		ServerAt string `json:"serverAt"`
		ServerId int    `json:"serverId"`
	}{
		Offset:   t,
		ServerAt: time.Now().Format("2006-01-02 15-04-05"),
		ServerId: tm.serverId,
	}
	data, err := json.Marshal(content)
	if err != nil {
		return err
	}
	if tm.url == "" {
		return fmt.Errorf("send ")
	}
	fullPath := tm.url + fmt.Sprintf(UrlPathTime, tm.serverId)
	req, err := http.NewRequest(http.MethodPost, fullPath, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("gameid", tm.gameId)
	// 请求对应的数据
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("send to web code %d", resp.StatusCode)
	}
	return nil
}
