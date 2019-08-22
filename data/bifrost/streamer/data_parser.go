package streamer

import (
	"github.com/Mintegral-official/mtggokit/data/container"
)

type DataParser interface {
	Parse([]byte) (container.MapKey, interface{}, error)
}