package container

import (
	"errors"
	"sync"
)

// 双bufMap, 仅提供Get/LoadBase接口
type BufferedMap struct {
	innerData map[MapKey]interface{}
	mutex     sync.Mutex
	ErrorNum  int
}

func (bm *BufferedMap) Get(key MapKey) (interface{}, error) {
	data, in := bm.innerData[key]
	if !in {
		return nil, errors.New("Not exist")
	}
	return data, nil
}

func (bm *BufferedMap) LoadBase(iterator DataIterator) error {
	bm.ErrorNum = 0
	tmpM := make(map[MapKey]interface{})
	for iterator.HasNext() {
		_, k, v, e := iterator.Next()
		if e != nil {
			bm.ErrorNum++
			continue
		}
		tmpM[k] = v
	}
	bm.mutex.Lock()
	bm.innerData = tmpM
	defer bm.mutex.Lock()
	return nil
}
