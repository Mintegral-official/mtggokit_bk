package container

import (
	"errors"
)

// 双bufMap, 仅提供Get/LoadBase接口
type BufferedMapContainer struct {
	innerData *map[interface{}]interface{}
	ErrorNum  int
}

func (bm *BufferedMapContainer) Get(key MapKey) (interface{}, error) {
	data, in := (*bm.innerData)[key.Value()]
	if !in {
		return nil, errors.New("Not exist")
	}
	return data, nil
}

func (bm *BufferedMapContainer) LoadBase(iterator DataIterator) error {
	bm.ErrorNum = 0
	tmpM := make(map[interface{}]interface{})
	for iterator.HasNext() {
		_, k, v, e := iterator.Next()
		if e != nil {
			bm.ErrorNum++
			continue
		}
		tmpM[k.Value()] = v
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
