package report

import (
	"github.com/simplemoon/deploy/conf"
)

const (
	ResponseCodeSuccess = "10000" // 执行成功
	ResponseCodeFailed  = "10001" // 执行失败

	Success = "Succeed" // 成功的状态
	Failed  = "Failed"  // 失败状态
)

// 详细信息
type RowInfo struct {
	ServerIdx  int    `json:"idx"`     // 目录的序号
	Action     string `json:"act"`     // 命令名称
	ActionType string `json:"type"`    // 命令类型
	State      string `json:"state"`   // 状态
	Msg        string `json:"message"` // 输出信息
}

// 返回的结果
type ResponseResult struct {
	Code       string    `json:"code"`    // 返回码
	Details    []RowInfo `json:"obj"`     // 详细的信息
	Production int       `json:"proj"`    // 项目的序号
	Version    string    `json:"version"` // 版本号
}

// 创建一个结果数据
func NewResult() *ResponseResult {
	return &ResponseResult{
		Code:       ResponseCodeSuccess,
		Details:    make([]RowInfo, 0),
		Production: conf.GetProject(),
		Version:    conf.GetVersion(),
	}
}

// 报告错误的实例
func NewErrReport(msg string) *ResponseResult {
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
		Production: conf.GetProject(),
		Version:    conf.GetVersion(),
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

// 报告结果
func (r *ResponseResult) Report() {
	// TODO: TO REPORT RESULT TO STDERR OR WEB OSS
}
