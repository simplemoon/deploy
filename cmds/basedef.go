package cmds

// 后台配置的命令
const (
	CmdNameNone          = "base"            // 未知的名称
	CmdNameUpdate        = "update"          // 更新服务器
	CmdNameStart         = "start"           // 启动
	CmdNameStatus        = "status"          // 查询状态
	CmdNameStop          = "stop"            // 关闭服务器
	CmdNameUninstall     = "uninstall"       // 删除服务
	CmdNameClearDB       = "clear_db"        // 清理数据库
	CmdNameInit          = "init"            // 重置服务器进程配置文件
	CmdNameCTime         = "change_time"     // 修改时间
	CmdNameCState        = "change_state"    // 修改状态
	CmdNameClearRankList = "clear_rank_list" // 清理排行榜
	CmdNameCloseGame     = "close_member"    // 关闭游戏服
	CmdNameClosePlatform = "close_platform"  // 关闭平台服
	CmdNameDiff          = "diff"            // 检查配置文件是否正确
	CmdNameConfig        = "config"          // 重新配置
	CmdNameDelete        = "delete"          // 删除配置
	CmdNameReplaceRes    = "replace_res"     // 替换游戏服务器资源
	CmdNameCheckVersion  = "version_check"   // 检查版本信息
)

// service 的命令
const (
	ServiceCmdNameStart     = "start"     // 启动
	ServiceCmdNameInstall   = "install"   // 安装
	ServiceCmdNameStop      = "stop"      // 停止
	ServiceCmdNameUnInstall = "uninstall" // 卸载
)
