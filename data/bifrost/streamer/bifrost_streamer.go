package streamer

import "github.com/Mintegral-official/mtggokit/data/container"

type BiFrostStreamerCfg struct {
	Name         string
	Version      int //数据格式的版本
	IpPort       string
	BaseFilePath string
}

type BiFrostStreamer struct {
	Cfg       *BiFrostStreamerCfg
	container container.Container
}

func NewBiFrostStreamer(cfg *BiFrostStreamerCfg) *BiFrostStreamer {
	return &BiFrostStreamer{Cfg: cfg}
}

func UpdateData() {
	//先加载基准，然后后台增量
}
