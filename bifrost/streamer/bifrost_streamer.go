package streamer

import (
	"github.com/Mintegral-official/mtggokit/bifrost/container"
)

type FileStruct struct {
	Name        string
	UpdateTime  int64
	DataVersion int
	Data        map[container.MapKey]interface{}
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
