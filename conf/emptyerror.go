package conf

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/simplemoon/deploy/report"
)

// 报告错误
func FormatErr(v interface{}) {
	msg := fmt.Sprintf("%v", v)
	msg = strings.Replace(msg, "'", "", -1)
	// 创建一个结果
	r := report.NewErrReport(msg, GetProject(), GetVersion())
	// 反序列化成字符串
	info, err := json.Marshal(r)
	if err != nil {
		fmt.Println("json marshal failed ", err)
		return
	}
	os.Stderr.Write(info)
}
