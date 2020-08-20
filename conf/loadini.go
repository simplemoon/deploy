package conf

import (
	"gopkg.in/ini.v1"
	"sync"
)

var (
	telnetCfgCache = newCfgCache()
)

// 缓存的数据
type CfgCache struct {
	sync.Mutex
	data map[string]*TelnetCfg // 对应的记录
}

// 创建一个缓存
func newCfgCache() *CfgCache {
	return &CfgCache{
		data: make(map[string]*TelnetCfg),
	}
}

// 获取
func (cc *CfgCache) get(name string) *TelnetCfg {
	cfg, ok := cc.data[name]
	if !ok {
		return nil
	}
	return cfg
}

// 存储
func (cc *CfgCache) put(name string, cfg *TelnetCfg) {
	cc.data[name] = cfg
}

// 需要的数据
type TelnetCfg struct {
	EchoAddr string // 连接的地址
	EchoPort int    // 连接的端口号
}

// 加载对应的配置文件路径
func LoadIni(iniPath string) (*TelnetCfg, error) {
	data := telnetCfgCache.get(iniPath)
	if data != nil {
		return data, nil
	}
	// 继续获取
	telnetCfgCache.Lock()
	defer telnetCfgCache.Unlock()
	// 在获取一下
	data = telnetCfgCache.get(iniPath)
	if data != nil {
		return data, nil
	}
	// 获取对应的数据
	cfg, err := ini.Load(iniPath)
	if err != nil {
		return nil, err
	}
	// 获取对应的IP和地址
	section := cfg.Section(IniKeySectionEcho)
	addr := section.Key(IniKeyAddr).String()
	port, err := section.Key(IniKeyPort).Int()
	if err != nil {
		return nil, err
	}
	// 存入对应的数据
	data = &TelnetCfg{EchoAddr: addr, EchoPort: port}
	telnetCfgCache.put(iniPath, data)
	return data, nil
}
