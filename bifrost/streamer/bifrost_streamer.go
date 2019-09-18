package streamer

import (
	"context"
	"github.com/Mintegral-official/mtggokit/bifrost/container"
)

type BifrostStreamer struct {
	Cfg       *BiFrostStreamerCfg
	container container.Container
}

func NewBiFrostStreamer(cfg *BiFrostStreamerCfg) *BifrostStreamer {
	return &BifrostStreamer{Cfg: cfg}
}

func (bs *BifrostStreamer) SetContainer(container container.Container) {
	bs.container = container
}

func (bs *BifrostStreamer) GetContainer() container.Container {
	return bs.container
}

func (bs *BifrostStreamer) GetSchedInfo() *SchedInfo {
	return &SchedInfo{
		TimeInterval: bs.Cfg.Interval,
	}
}

func (bs *BifrostStreamer) HasNext() bool {
	return false
}

func (bs *BifrostStreamer) Next() (container.DataMode, container.MapKey, interface{}, error) {
	return container.DataModeAdd, nil, nil, nil
}

func (fs *BifrostStreamer) UpdateData(ctx context.Context) error {
	return nil
}
