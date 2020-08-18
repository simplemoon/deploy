package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"
)

// 链接的状态
const (
	ConnectStateFailed   = "NotConnect"
	ConnectStateWriteErr = "WriteFailed"
	ConnectStateReadErr  = "ReadFailed"
	ConnectStateParseErr = "ParsedFailed"
	ConnectStateKeyErr   = "FoundKeyFailed"
)

var (
	helpers = sync.Map{}
)

// 连接的参数
type TelnetKey struct {
	addr string // 地址
	port int    // 端口号
}

func (tk *TelnetKey) GetConnectString() string {
	return fmt.Sprintf("%s:%d", tk.addr, tk.port)
}

// telnet 连接地址
type TelnetHelper struct {
	key    *TelnetKey    // 键值
	conn   net.Conn      // 链接
	reader *bufio.Reader // 读取
}

// 创建一个helper
func NewTelnet(addr string, port int) *TelnetHelper {
	key := TelnetKey{addr: addr, port: port}
	val, ok := helpers.Load(key)
	if ok {
		return val.(*TelnetHelper)
	}
	helper := &TelnetHelper{
		key:    &key,
		conn:   nil,
		reader: nil,
	}
	helpers.Store(key, helper)
	return helper
}

// 发送命令
func (t *TelnetHelper) Send(cmd string) string {
	// 连接一下
	err := t.connect()
	if err != nil {
		return ConnectStateFailed
	}
	// 写入对应的数据
	cmd = fmt.Sprintf("%s\r\n", cmd)
	_, err = t.conn.Write([]byte(cmd))
	if err != nil {
		return ConnectStateWriteErr
	}
	// 等待返回啊
	result, err := t.reader.ReadString('}')
	if err != nil {
		return ConnectStateReadErr
	}
	// 加载对应的数据
	infos := make(map[string]string, 0)
	err = json.Unmarshal([]byte(result), &infos)
	if err != nil {
		return ConnectStateParseErr
	}
	s, ok := infos["status"]
	if !ok {
		return ConnectStateKeyErr
	}
	return s
}

// 连接
func (t *TelnetHelper) connect() error {
	if t.conn != nil {
		conn, err := net.DialTimeout("tcp", t.key.GetConnectString(), time.Second*3)
		if err != nil {
			return err
		}
		t.conn = conn
		t.reader = bufio.NewReader(t.conn)
	}
	return nil
}
