package report

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const (
	ResponseCodeSuccess = "10000" // 执行成功
	ResponseCodeFailed  = "10001" // 执行失败

	Success = "Succeed" // 成功的状态
	Failed  = "Failed"  // 失败状态
)

// 详细信息
type RowInfo struct {
	ServerIdx  int         `json:"idx"`     // 目录的序号
	Action     string      `json:"act"`     // 命令名称
	ActionType string      `json:"type"`    // 命令类型
	State      string      `json:"state"`   // 状态
	Msg        interface{} `json:"message"` // 输出信息
}

// 返回的结果
type ResponseResult struct {
	Code       string    `json:"code"`    // 返回码
	Details    []RowInfo `json:"obj"`     // 详细的信息
	Production int       `json:"proj"`    // 项目的序号
	Version    string    `json:"version"` // 版本号
}

// 创建一个结果数据
func NewResult(project int, version string) *ResponseResult {
	return &ResponseResult{
		Code:       ResponseCodeSuccess,
		Details:    make([]RowInfo, 0),
		Production: project,
		Version:    version,
	}
}

// 报告错误的实例
func NewErrReport(msg string, project int, version string) *ResponseResult {
	r := &ResponseResult{
		Code: ResponseCodeFailed,
		Details: []RowInfo{
			{
				ServerIdx:  0,
				Action:     "main",
				ActionType: "begin",
				State:      Failed,
				Msg:        msg,
			},
		},
		Production: project,
		Version:    version,
	}
	return r
}

// 结果
func (r *ResponseResult) AddSuc(idx int, act, actType, msg string) {
	// 单行的信息
	row := RowInfo{
		ServerIdx:  idx,
		Action:     act,
		ActionType: actType,
		State:      Success,
		Msg:        msg,
	}
	// 具体的信息
	r.Details = append(r.Details, row)
}

// 失败的结果
func (r *ResponseResult) AddFailed(idx int, act, actType, msg string) {
	row := RowInfo{
		ServerIdx:  idx,
		Action:     act,
		ActionType: actType,
		State:      Failed,
		Msg:        msg,
	}
	// 设置状态为失败状态
	r.Code = ResponseCodeFailed
	r.Details = append(r.Details, row)
}

// 增加多行
func (r *ResponseResult) AddRows(rows []RowInfo) {
	for i := 0; i < len(rows); i++ {
		if rows[i].State != Success {
			r.Code = ResponseCodeFailed
		}
		r.Details = append(r.Details, rows[i])
	}
}

// 获取状态
func (r *ResponseResult) GetResult() string {
	if r.Code == ResponseCodeSuccess {
		return Success
	}
	return Failed
}

// 报告结果
func (r *ResponseResult) Report(url string, model bool) error {
	content, err := json.Marshal(r)
	if err != nil {
		return err
	}
	// 是否是远程模式
	if model {
		_, err = os.Stderr.Write(content)
	} else {
		ret, err := http.Post(url, "application/json", bytes.NewReader(content))
		if err != nil {
			return err
		}
		if ret.StatusCode != 200 {
			return fmt.Errorf("send %v to %s failed, code: %d", content, url, ret.StatusCode)
		}
	}
	return err
}

// 报告结果
func (r *ResponseResult) DingNotify(url, gameId, action string, serverId int) error {
	// 发送钉钉通知
	// 写入对应的数据
	ddInfo := map[string]interface{}{
		"type":     "Ops",
		"serverId": serverId,
		"message":  fmt.Sprintf("%s %s", r.Version, r.GetResult()),
		"action":   action,
		"gameId":   gameId,
	}
	// 获取对应的数据
	content, err := json.Marshal(ddInfo)
	if err != nil {
		return fmt.Errorf("marshal info %v err %v", ddInfo, err)
	}
	// 发送消息
	ret, err := http.Post(url, "application/json", bytes.NewReader(content))
	if err != nil {
		return err
	}
	if ret.StatusCode != 200 {
		return fmt.Errorf("send %v to %s failed, code: %d", content, url, ret.StatusCode)
	}
	return nil
}
