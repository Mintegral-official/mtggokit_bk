package container

import (
	"errors"
	"sync"
)

// 双bufMap, 仅提供Get/LoadBase接口
type BufferedMap struct {
	innerData *map[interface{}]interface{}
	mutex     sync.Mutex
	ErrorNum  int
}

func (bm *BufferedMap) Get(key MapKey) (interface{}, error) {
	data, in := (*bm.innerData)[key.Value()]
	if !in {
		return nil, errors.New("Not exist")
	}
	return data, nil
}

func (bm *BufferedMap) LoadBase(iterator DataIterator) error {
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
	bm.mutex.Lock()
	bm.innerData = &tmpM
	defer bm.mutex.Unlock()
	return nil
}
