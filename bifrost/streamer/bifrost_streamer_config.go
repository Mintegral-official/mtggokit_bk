package streamer

type BiFrostStreamerCfg struct {
	Name         string
	Version      int //数据格式的版本
	Ip           string
	Port         int
	BaseFilePath string
	Interval     int
	IsSync       bool
	IsOnline     bool // 离线模式生效
	WriteFile    bool // 离线模式生效
	CacheSize    int
}
