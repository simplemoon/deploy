package conf

type IntMap = map[string]int            // 名字对应的端口配置文件
type MultiStringMap = map[string]IntMap // 项目的文件

// 基础的配置
type CfgPortBase struct {
	PortStart int `json:"PortStart"` // 开始的端口
	PortStep  int `json:"PortStep"`  // 端口的跨度
}

// port template
type CfgPort struct {
	Base    CfgPortBase    `json:"base"`    // 基础的配置
	Replace MultiStringMap `json:"replace"` // 替代的开始端口
	Data    MultiStringMap `json:"data"`    // 具体的信息
}

// 获取端口信息
func (cp *CfgPort) GetPortInfo(name string) IntMap {
	v, ok := cp.Data[name]
	if !ok {
		return nil
	}
	return v
}

// 获取基础端口号
func (cp *CfgPort) GetBasePort(name, key string, idx int) (int, bool) {
	// 基础的端口信息
	p := cp.Base.PortStart + idx*cp.Base.PortStep

	d, ok := cp.Replace[name]
	if !ok {
		return p, false
	}
	r, ok := d[key]
	if !ok {
		return p, false
	}
	// 找到最大的数量
	cnt, ok := d[InfoKeyMaxCnt]
	if !ok {
		return p, false
	}
	return r + idx*cnt, true
}

// 需要发送到后台的信息
type ReportItem struct {
	Key       string `json:"key"`
	ReportKey string `json:"reportKey"`
	Count     string `json:"countKey"`
}

// 是否为null
func (ri *ReportItem) IsNil() bool {
	return ri.Key == "" || ri.ReportKey == ""
}
