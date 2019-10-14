package container

import (
	"fmt"
	"github.com/pkg/errors"
	"sync"
)

// 多线程读写安全的container，支持增量
type BlockingMapContainer struct {
	innerData *sync.Map
	errorNum  int64
	totalNum  int64
	Tolerate  float64
}

func CreateBlockingMapContainer(numPartision int, tolerate float64) *BlockingMapContainer {
	return &BlockingMapContainer{
		innerData: &sync.Map{},
		Tolerate:  tolerate,
	}
}

func (bm *BlockingMapContainer) Get(key MapKey) (interface{}, error) {
	data, in := bm.innerData.Load(key.Value())
	if !in {
		return nil, errors.New("Not exist")
	}
	return data, nil
}

func (bm *BlockingMapContainer) Set(key MapKey, value interface{}) error {
	bm.innerData.Store(key.Value(), value)
	return nil
}

func (bm *BlockingMapContainer) Del(key MapKey, value interface{}) {
	bm.innerData.Delete(key.Value())
}

func (bm *BlockingMapContainer) LoadBase(iterator DataIterator) error {
	tmpM := &sync.Map{}
	bm.errorNum = 0
	bm.totalNum = 0
	for iterator.HasNext() {
		m, k, v, e := iterator.Next()
		bm.totalNum++
		if e != nil {
			bm.errorNum++
			continue
		}
		switch m {
		case DataModeAdd, DataModeUpdate:
			tmpM.Store(k.Value(), v)
		case DataModeDel:
			tmpM.Delete(k.Value())
		}
	}
	if bm.totalNum == 0 {
		bm.totalNum = 1
	}
	f := float64(bm.errorNum) / float64(bm.totalNum)
	if f > bm.Tolerate {
		return errors.New(fmt.Sprintf("LoadBase error, tolerate[%f], err[%f]", bm.Tolerate, f))
	}
	bm.innerData = tmpM
	return nil
}

func (bm *BlockingMapContainer) LoadInc(iterator DataIterator) error {
	for iterator.HasNext() {
		m, k, v, e := iterator.Next()
		bm.totalNum++
		if e != nil {
			bm.errorNum++
			continue
		}
		switch m {
		case DataModeAdd, DataModeUpdate:
			bm.innerData.Store(k.Value(), v)
		case DataModeDel:
			bm.Del(k, v)
		}
	}
	if bm.totalNum == 0 {
		bm.totalNum = 1
	}
	f := float64(bm.errorNum) / float64(bm.totalNum)
	if f > bm.Tolerate {
		return errors.New(fmt.Sprintf("LoadInc error, tolerate[%f], err[%f]", bm.Tolerate, f))
	}
	return nil
}
