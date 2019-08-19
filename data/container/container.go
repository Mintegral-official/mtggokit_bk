package container

type MapKey interface {
	PartitionKey() int64
	Value() interface{}
}

type Container interface {
	Get(key MapKey) (interface{}, error)
	Set(key MapKey, value interface{}) error
	Del(key MapKey) error

	LoadBase(dataIter DataIterator) error
	LoadInc(dataIter DataIterator) error
}
