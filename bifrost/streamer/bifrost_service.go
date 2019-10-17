package streamer

import (
	"bytes"
	"context"
	"encoding/gob"
	"github.com/Mintegral-official/mtggokit/bifrost/container"
	"github.com/pkg/errors"
)

type Status int

const (
	Ok Status = iota
	Error
)

func init() {
	gob.Register(&container.Int64Key{})
	gob.Register(&container.StringKey{})
}

func (s *Status) toString() string {
	switch *s {
	case Ok:
		return "Ok"
	default:
		return "unknown"

	}
}

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

type BifrostService struct {
	StreamerManager *StreamerProviderManager
}

func NewBifrostServer(streamerManager *StreamerProviderManager) *BifrostService {
	return &BifrostService{StreamerManager: streamerManager}
}

func (bs *BifrostService) GetBase(ctx context.Context, req *BaseReq, res *BaseRes) error {
	sp := bs.StreamerManager.GetProvider(req.Name, req.Progress)
	if sp == nil {
		return errors.New("Not found streamer[" + req.Name + "]")
	}
	base := sp.GetBase()
	if base == nil {
		return errors.New("Get baseInfo error, streamer[" + req.Name + "]")
	}
	res.Status = Ok
	res.BaseInfo = base
	return nil
}

func (bs *BifrostService) GetInc(ctx context.Context, req *IncReq, res *IncRes) error {
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

type GobCodec struct {
}

func (c *GobCodec) Decode(data []byte, i interface{}) error {
	enc := gob.NewDecoder(bytes.NewBuffer(data))
	err := enc.Decode(i)
	return err
}

func (c *GobCodec) Encode(i interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(i)
	return buf.Bytes(), err
}
