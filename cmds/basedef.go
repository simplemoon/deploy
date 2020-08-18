package cmds

const (
	CmdNameNone          = "base"           // 未知的名称
	CmdNameStart         = "start"          // 启动
	CmdNameCTime         = "change_time"    // 修改时间
	CmdNameCState        = "change_state"   // 修改状态
	CmdNameClearDB       = "clear_db"       // 清理数据库
	CmdNameClearRankList = "clear_ranklist" // 清理排行榜
	CmdNameCloseGame     = "close_member"   // 关闭游戏服
	CmdNameClosePlatform = "close_platform" // 关闭平台服
	CmdNameDiff          = "diff"           // 检查配置文件是否正确
	CmdNameInit          = "init"           // 重置数据库
	CmdNamePlayers       = "players"        // 查询玩家数量
	CmdNameReConfig      = "reconfig"       // 重新配置
	CmdNameDelete        = "remove_cfg"     // 删除配置
	CmdNameReplaceRes    = "reset_ini"      // 重置游戏服务器配置
	CmdNameStatus        = "status"         // 查询状态
	CmdNameStop          = "stop"           // 关闭服务器
	CmdNameUninstall     = "uninstall"      // 删除服务
	CmdNameUpdate        = "update"         // 更新服务器
	CmdNameCheckVersion  = "version_check"  // 检查版本信息
)
