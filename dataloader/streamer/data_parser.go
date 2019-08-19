package streamer

import (
	"mtggokits/datacontainer"
)

type DataParser interface {
	Parse([]byte) (datacontainer.MapKey, interface{}, error)
}
