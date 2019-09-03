package container

import (
	"github.com/easierway/concurrent_map"
	"github.com/pkg/errors"
)

// 多线程读写安全的container，支持增量
type BlockingMapContainer struct {
	innerData    *concurrent_map.ConcurrentMap
	numPartision int
	ErrorNum     int
}

func CreateBlockingMapContainer(numPartision int) *BlockingMapContainer {
	return &BlockingMapContainer{
		innerData:    concurrent_map.CreateConcurrentMap(numPartision),
		numPartision: numPartision,
	}
}

func (bm *BlockingMapContainer) Get(key MapKey) (interface{}, error) {
	data, in := bm.innerData.Get(key)
	if !in {
		return nil, errors.New("Not exist")
	}
	return data, nil
}

func (bm *BlockingMapContainer) Set(key MapKey, value interface{}) error {
	bm.innerData.Set(key, value)
	return nil
}

func (bm *BlockingMapContainer) Del(key MapKey, value interface{}) {
	bm.innerData.Del(key)
}

func (bm *BlockingMapContainer) LoadBase(iterator DataIterator) error {
	tmpM := concurrent_map.CreateConcurrentMap(bm.numPartision)
	bm.ErrorNum = 0
	for iterator.HasNext() {
		m, k, v, e := iterator.Next()
		if e != nil {
			bm.ErrorNum++
			continue
		}
		switch m {
		case DataModeAdd, DataModeUpdate:
			tmpM.Set(k, v)
		case DataModeDel:
			tmpM.Del(k)
		}
	}
	bm.innerData = tmpM
	return nil
}

func (bm *BlockingMapContainer) LoadInc(iterator DataIterator) error {
	for iterator.HasNext() {
		m, k, v, e := iterator.Next()
		if e != nil {
			bm.ErrorNum++
			continue
		}
		switch m {
		case DataModeAdd, DataModeUpdate:
			bm.innerData.Set(k, v)
		case DataModeDel:
			bm.Del(k, v)
		}
	}
	return nil
}
