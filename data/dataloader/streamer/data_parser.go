package streamer

import (
	"mtggokits/data/container"
)

type DataParser interface {
	Parse([]byte) (container.MapKey, interface{}, error)
}
