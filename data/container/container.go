package container

// key of the map, because of go-lang not support generic type，
// So, here defined an interface for int value or string value key
type MapKey interface {
	PartitionKey() int64
	Value() interface{}
}

type Container interface {
	Get(key MapKey) (interface{}, error)
	Set(key MapKey, value interface{}) error
	Del(key MapKey, value interface{})

	LoadBase(dataIter DataIterator) error
	LoadInc(dataIter DataIterator) error
}
