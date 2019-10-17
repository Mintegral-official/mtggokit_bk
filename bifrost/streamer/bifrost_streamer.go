package streamer

import (
	"context"
	"errors"
	"fmt"
	"github.com/Mintegral-official/mtggokit/bifrost/container"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/protocol"
	"github.com/smallnest/rpcx/share"
	"time"
)

type BifrostStreamer struct {
	cfg        *BiFrostStreamerCfg
	container  container.Container
	client     client.XClient
	incRecords []*IncRecord
	curIdx     int
	progress   int64
	totalNum   int64
	errorNum   int64
	startTime  int64
	endTime    int64
	hasInit    bool
}

func NewBiFrostStreamer(cfg *BiFrostStreamerCfg) *BifrostStreamer {
	d := client.NewPeer2PeerDiscovery("tcp@"+cfg.URI, "")
	option := client.DefaultOption
	option.SerializeType = protocol.SerializeType(5)
	share.Codecs[protocol.SerializeType(5)] = &GobCodec{}
	c := client.NewXClient("BifrostService", client.Failtry, client.RandomSelect, d, option)
	return &BifrostStreamer{
		cfg:    cfg,
		client: c,
	}
}

func (bs *BifrostStreamer) SetContainer(container container.Container) {
	bs.container = container
}

func (bs *BifrostStreamer) GetContainer() container.Container {
	return bs.container
}

func (bs *BifrostStreamer) GetSchedInfo() *SchedInfo {
	return &SchedInfo{
		TimeInterval: bs.cfg.Interval,
	}
}

func (bs *BifrostStreamer) HasNext() bool {
	return bs.curIdx < len(bs.incRecords)
}

func (bs *BifrostStreamer) Next() (container.DataMode, container.MapKey, interface{}, error) {
	bs.totalNum++
	if bs.curIdx < 0 || bs.curIdx >= len(bs.incRecords) {
		bs.errorNum++
		return container.DataModeAdd, nil, nil, errors.New(fmt.Sprintf("invalid index[%d],cap[%d]", bs.curIdx, len(bs.incRecords)))
	}
	r := bs.incRecords[bs.curIdx]
	bs.curIdx++
	return r.DataMode, r.MapKey, r.Value, nil
}

func (bs *BifrostStreamer) UpdateData(ctx context.Context) error {
	if bs.cfg.IsSync {
		bs.startTime = time.Now().UnixNano()
		err := bs.loadBase(ctx)
		bs.endTime = time.Now().UnixNano()
		if err != nil {
			bs.WarnStatus("Sync LoadBase error: " + err.Error())
			return err
		}
		bs.InfoStatus("Sync LoadBase succ")
		bs.hasInit = true
	}

	go func() {
		if !bs.hasInit {
			bs.startTime = time.Now().UnixNano()
			err := bs.loadBase(ctx)
			bs.endTime = time.Now().UnixNano()
			if err != nil {
				bs.WarnStatus("Async LoadBase error: " + err.Error())
			}
			bs.InfoStatus("Async LoadBase succ")
			bs.hasInit = true
		}
		for {
			inc := time.After(time.Duration(bs.cfg.Interval) * time.Second)
			select {
			case <-ctx.Done():
				return
			case <-inc:
				bs.startTime = time.Now().UnixNano()
				err := bs.loadInc(ctx)
				bs.endTime = time.Now().UnixNano()
				if err != nil {
					bs.WarnStatus("LoadBase error: " + err.Error())
				} else {
					bs.InfoStatus("LoadBase succ")
				}
			}
		}
	}()

	return nil
}

func (bs *BifrostStreamer) loadBase(ctx context.Context) error {
	bs.totalNum = 0
	bs.errorNum = 0
	res := &BaseRes{}
	err := bs.client.Call(ctx, "GetBase", &BaseReq{
		Name:     bs.cfg.Name,
		Progress: bs.progress,
	}, res)
	if err != nil {
		return err
	}
	if res.Status != Ok {
		return errors.New("GetBase error, status[" + res.Status.toString() + "]")
	}
	if res.BaseInfo == nil {
		return errors.New("GetBase error, baseInfo is nil")
	}
	bs.progress = res.BaseInfo.Progress
	records := make([]*IncRecord, 0, len(res.BaseInfo.Data))
	fmt.Println("xxxxxxxxxxxxxxxxxxxxxx: ", bs.progress, len(res.BaseInfo.Data))
	for k, v := range res.BaseInfo.Data {
		records = append(records, &IncRecord{
			DataMode: container.DataModeAdd,
			MapKey:   k,
			Value:    v,
			Progress: bs.progress,
		})
	}
	bs.incRecords = records
	return bs.container.LoadBase(bs)
}

func (bs *BifrostStreamer) loadInc(ctx context.Context) error {
	res := &IncRes{}
	err := bs.client.Call(ctx, "GetInc", &IncReq{
		Name:  bs.cfg.Name,
		Batch: 1000,
	}, res)
	if err != nil {
		return err
	}
	if res.Status != Ok {
		return errors.New("GetInc error, status[" + res.Status.toString() + "]")
	}
	bs.incRecords = res.IncRecords
	return bs.container.LoadInc(bs)
}

func (fs *BifrostStreamer) InfoStatus(s string) {
	if fs.cfg.Logger != nil {
		fs.cfg.Logger.Infof("streamer[%s] %s, totalNum[%d], errorNum[%d], userData[%s], timeUsed[%d]", fs.cfg.Name, s, fs.totalNum, fs.errorNum, fs.cfg.UserData, (fs.endTime-fs.startTime)/10e6)
	}
}

func (fs *BifrostStreamer) WarnStatus(s string) {
	if fs.cfg.Logger != nil {
		fs.cfg.Logger.Warnf("streamer[%s] %s, totalNum[%d], errorNum[%d], userData[%s], timeUsed[%d]", fs.cfg.Name, s, fs.totalNum, fs.errorNum, fs.cfg.UserData, (fs.endTime-fs.startTime)/10e6)
	}
}
