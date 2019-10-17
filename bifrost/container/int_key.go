package container

type Int64Key struct {
	Data int64
}

func (i *Int64Key) PartitionKey() int64 {
	return i.Data
}

func (i *Int64Key) Value() interface{} {
	return i.Data
}

func I64Key(key int64) *Int64Key {
	return &Int64Key{key}
}
