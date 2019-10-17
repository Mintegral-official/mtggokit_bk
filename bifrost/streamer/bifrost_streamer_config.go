package streamer

import "github.com/Mintegral-official/mtggokit/bifrost/log"

type BiFrostStreamerCfg struct {
	Name         string      // streamer名字
	NameSpace    string      // streamer命名空间
	Version      int         // 数据格式的版本
	URI          string      //
	BaseFilePath string      // 基准文件路径
	Interval     int         // 增量更新时间间隔
	IsSync       bool        // 是否同步加载
	IsOnline     bool        // 离线模式生效
	WriteFile    bool        // 离线模式生效
	UserData     interface{} // 用户自定义数据
	Logger       log.BiLogger
}
