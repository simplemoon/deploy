package conf

// 配置文件名称
const (
	FileNameBase       = "base.json"     // 基础的配置信息
	FileNameProcess    = "process.json"  // 需要启动的进程信息
	FileNameExtra      = "extra.json"    // 需要替换的文件
	FileNamePort       = "ports.json"    // 端口配置
	FileNameProcIni    = "proc_ini.json" // 服务器的ini配置信息
	FileNameServiceXml = "service.xml"   // 服务的进程配置文件
	FileNameServiceExe = "service.exe"   // 服务的执行文件
	FileNameReport     = "report.json"   // 需要上传的信息
	FileNameResetDb    = "reset_db.bat"  // 重置脚本的脚本
	FileNameReplace    = "replace.json"  // 替换的文件
	FileNameCheck      = "check.json"    // 需要检查的文件
)

// runtime 配置文件的名称
const (
	FileRuntimeServers    = "servers.json" // 服务器ID对应的idx的信息配置。
	FileRuntimeServerLock = "servers.lock" // 服务器ID对应的锁文件
)

// 游戏服务器下目录下面的文件
const (
	FileBin64LaunchExe  = "launch.exe"     // 执行的exe文件
	FileExtraSystemTime = "systemtime.ini" // 系统配置文件
	FileVersion         = "version.ini"    // 版本配置文件
)

// ini section 名称
const (
	IniKeySectionEcho = "Echo"
	IniKeySectionMain = "Main"

	IniKeyAddr      = "Addr"
	IniKeyPort      = "Port"
	IniKeyTimeAhead = "TimeAheadSeconds"
)

// ini文件的名称
const (
	ProjectName = "dq"

	PreIniFileName = "sd_" // ini文件名称的前缀

	ProcessKeyGame     = "game"    // world配置的名称
	ProcessKeyPlatform = "roommgr" // 平台服的配置
	ProcessKeyGm       = "gm"      // gm进程信息

	ProcessNameWorld    = "world"    // world服
	ProcessNamePlatform = "platform" // 平台服

	ProcessKeyPortFmt = "_%s_%s%dPort" // 进程端口号key的格式

	DefaultDataBasePort = "3306"
)

// 标识符号
const (
	FlagQuery    = '$'
	FlagSlice    = '*'
	FilterString = "/"
	IndentString = "\r\n"
)

// 命令的定义
const (
	EchoCmdStatus         = "status"                              // 状态查询
	EchoCmdQuit           = "quit"                                // 服务器退出
	EchoCmdBeforeQuit     = "bquit"                               // 通知玩家服务器要关闭了
	EchoCmdDisconnectRoom = "setp SoloWorld UseRoomService false" // 断开和平台服的连接

)

// 字段的关键字
const (
	InfoKeyGameType    = "type"       // 服务器类型
	InfoKeyMemberCount = "member"     // 场景服个数
	InfoKeyIsPlatform  = "isPlatform" // 是否是平台
	InfoKeyCrossType   = "crossType"  // 跨服类型
	InfoKeyIsCharge    = "isCharge"   // 是否开启计费
	InfoKeyMaxCnt      = "MaxCnt"     // 最大数量
	InfoKeyResUrl      = "resUrl"     // 资源的路径
	InfoKeyMd5         = "resMd5"     // md5验证码
	InfoKeyDBUser      = "dbUser"     // 数据库用户名
	InfoKeyDBPwd       = "dbPwd"      // 数据库密码
	InfoKeyDBAddr      = "dbAddr"     // 数据库的链接地址
	InfoKeyVersion     = "version"    // 版本信息
	InfoKeyToolArgs    = "toolArgs"   // 修改时间或者修改状态的参数
	InfoKeyGameId      = "gameId"     // 游戏编号
	InfoKeyDeployUrl   = "devopsApi"  // 运维的url
	InfoKeyChargeUrl   = "chargeApi"  // 计费的url
)

// 文件名称后缀
const (
	FileSuffixZip  = ".zip"
	FileSuffixLock = ".lock"
	FileSuffixIni  = ".ini"
)

// 服务器的类型定义
const (
	GSTypeGame  = "G" // 服务器类型
	GSTypeCross = "C" // 跨服类型
	GSTypeMix   = "M" // 混服类型

	CSTypePlatform = "P" // 平台服
	CSTypeCross    = "C" // 混服
	CSTypeRoom     = "R" // 房间服
)

// 进程的类型
const (
	ProcTypeBase     uint32 = 1 << iota // 所有的基础类型
	ProcTypeServer                      // 服务器的进程
	ProcTypeCharge                      // 计费的进程
	ProcTypePlatform                    // 平台的进程
	ProcTypeRoom                        // ROOM的进程

	ProcTypeContainsSvr      = ProcTypeServer | ProcTypeCharge
	ProcTypeContainsPlatform = ProcTypePlatform | ProcTypeRoom
	ProcTypeContainsAll      = ProcTypeBase | ProcTypeContainsSvr | ProcTypeContainsPlatform
)

// url路径定义
const (
	UrlPathDing     = "/admin/v1/xyz/notify"              // 钉钉通知
	UrlPathResult   = "/admin/v1/xyz/deployDetail"        // 报告的结果
	UrlPathSvrState = "/gs/v1/server/%d/syncState"        // 服务器状态
	UrlPathUpdate   = "/admin/v1/xyz/updateServer"        // 更新服务器的信息
	UrlPathTime     = "/gs/v1/server/%d/sync/change_time" // 时间的通知信息
)
