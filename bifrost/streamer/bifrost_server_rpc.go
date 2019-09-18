package streamer

import (
	"context"
	"github.com/pkg/errors"
)

type Status int

const (
	Ok Status = iota
	Error
)

type BaseReq struct {
	Name     string
	Space    string
	Progress int64
}

type BaseRes struct {
	Status   Status // streamer name
	BaseInfo *BaseInfo
}

type IncReq struct {
	Name     string
	Space    string
	Batch    int
	Progress int64
}

type IncRes struct {
	Status     Status // streamer name
	IncRecords []*IncRecord
}

type BifrostServer struct {
	StreamerManager *StreamerProviderManager
}

func NewBifrostServer(streamerManager *StreamerProviderManager) *BifrostServer {
	return &BifrostServer{StreamerManager: streamerManager}
}

func (bs *BifrostServer) GetBase(ctx context.Context, req *BaseReq, res *BaseRes) error {
	sp := bs.StreamerManager.GetProvider(req.Name, req.Progress)
	if sp == nil {
		return errors.New("Not found streamer[" + req.Name + "]")
	}
	base := sp.GetBase()
	if base == nil {
		return errors.New("Get baseInfo error, streamer[" + req.Name + "]")
	}
	res = &BaseRes{
		Status:   Ok,
		BaseInfo: base,
	}
	return nil
}

func (bs *BifrostServer) GetIncs(ctx context.Context, req *IncReq, res *IncRes) error {
	sp := bs.StreamerManager.GetProvider(req.Name, req.Progress)
	if sp == nil {
		return errors.New("Not found streamer[" + req.Name + "]")
	}
	inc, _ := sp.GetInc(req.Progress, req.Batch)
	if inc == nil {
		return errors.New("Get incInfo error, streamer[" + req.Name + "]")
	}
	res = &IncRes{
		Status:     Ok,
		IncRecords: inc,
	}
	return nil
}
