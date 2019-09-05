package container

import (
	"errors"
	"fmt"
)

// 双bufMap, 仅提供Get/LoadBase接口
type BufferedMapContainer struct {
	innerData *map[interface{}]interface{}
	errorNum  int64
	totalNum  int64
	Tolerate  float64
}

func (bm *BufferedMapContainer) Get(key MapKey) (interface{}, error) {
	data, in := (*bm.innerData)[key.Value()]
	if !in {
		return nil, errors.New("Not exist")
	}
	return data, nil
}

func (bm *BufferedMapContainer) LoadBase(iterator DataIterator) error {
	bm.errorNum = 0
	bm.totalNum = 0
	tmpM := make(map[interface{}]interface{})
	for iterator.HasNext() {
		_, k, v, e := iterator.Next()
		bm.totalNum++
		if e != nil {
			bm.errorNum++
			continue
		}
		tmpM[k.Value()] = v
	}
	if bm.totalNum == 0 {
		bm.totalNum = 1
	}
	f := float64(bm.errorNum) / float64(bm.totalNum)
	if f > bm.Tolerate {
		return errors.New(fmt.Sprintf("LoadBase error, tolerate[%f], err[%f]", bm.Tolerate, f))
	}
	bm.innerData = &tmpM
	return nil
}

func (bm *BufferedMapContainer) Set(key MapKey, value interface{}) error {
	return errors.New("not implement")
}

func (bm *BufferedMapContainer) Del(key MapKey, value interface{}) {
}

func (bm *BufferedMapContainer) LoadInc(iterator DataIterator) error {
	return errors.New("not implement")
}
