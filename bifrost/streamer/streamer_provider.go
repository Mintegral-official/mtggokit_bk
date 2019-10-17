package streamer

import (
	"fmt"
	"github.com/Mintegral-official/mtggokit/bifrost/container"
	"github.com/pkg/errors"
	"sync"
	"time"
)

type BaseInfo struct {
	Name        string
	Progress    int64
	DataVersion int
	Data        map[container.MapKey]interface{}
}

type IncRecord struct {
	DataMode container.DataMode
	MapKey   container.MapKey
	Progress int64 // 更新时间
	Value    interface{}
}

type StreamerProvider struct {
	Cfg      *StreamerProviderCfg
	BaseInfo *BaseInfo
	Cached   []*IncRecord
	lock     *sync.RWMutex
}

func NewStreamerProvider(cfg *StreamerProviderCfg) *StreamerProvider {
	return &StreamerProvider{
		Cfg:    cfg,
		Cached: make([]*IncRecord, 0),
		lock:   &sync.RWMutex{},
	}
}

func (sp *StreamerProvider) SetBase(baseInfo *BaseInfo) {
	sp.BaseInfo = baseInfo
}

func (sp *StreamerProvider) GetBase() *BaseInfo {
	return sp.BaseInfo
}

func (sp *StreamerProvider) AddInc(incs []*IncRecord) {
	idx := 0
	expireTime := time.Now().Unix() - sp.Cfg.ExpireTime
	for i, v := range sp.Cached {
		if v.Progress >= expireTime {
			idx = i
			break
		}
	}
	sp.lock.Lock()
	defer sp.lock.Unlock()
	if idx >= 0 {
		sp.Cached = sp.Cached[idx:]
	}
	for _, r := range incs {
		sp.Cached = append(sp.Cached, r)
	}
	if sp.Cfg.Logger != nil {
		sp.Cfg.Logger.Infof("StreamerProvider[%s] AddInc Succ, delete[%d], add[%d],total[%d]", sp.Cfg.Name, idx, len(incs), len(sp.Cached))
	}
}

func (sp *StreamerProvider) GetInc(progress int64, size int) ([]*IncRecord, error) {
	sp.lock.RLock()
	defer sp.lock.RUnlock()

	if sp.Cached == nil || len(sp.Cached) == 0 {
		return nil, errors.New(fmt.Sprintf("StreamerProvide[%s] LoadInc is nil", sp.Cfg.Name))
	}
	idx := BSearch(sp.Cached, progress)

	if idx == -1 {
		return nil, errors.New(fmt.Sprintf("StreamerProvide[%s] LoadInc error, Not found progress[%d], min[%d]", sp.Cfg.Name, progress, sp.Cached[0].Progress))
	}

	if idx+size > len(sp.Cached) {
		return sp.Cached[idx:], nil
	} else {
		return sp.Cached[idx:size], nil
	}
}

func BSearch(records []*IncRecord, progress int64) int {
	if len(records) == 0 || records[0].Progress > progress {
		return -1
	}
	idx := -1
	s, e := 0, len(records)
	for s > e {
		m := (s + e) / 2
		if records[m].Progress == progress {
			idx = m
			break
		} else if records[m].Progress > progress {
			e = m - 1
		} else {
			s = m + 1
		}
	}
	if idx == -1 {
		idx = s
	}
	return idx
}
