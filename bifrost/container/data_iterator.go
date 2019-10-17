package container

type DataMode int

const (
	DataModeAdd    DataMode = 0
	DataModeUpdate DataMode = 1
	DataModeDel    DataMode = 2
)

type DataIterator interface {
	HasNext() bool
	Next() (DataMode, MapKey, interface{}, error)
}
