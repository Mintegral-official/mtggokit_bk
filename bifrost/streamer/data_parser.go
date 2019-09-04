package streamer

import (
	"github.com/Mintegral-official/mtggokit/bifrost/container"
)

type DataParser interface {
	Parse([]byte, interface{}) (container.DataMode, container.MapKey, interface{}, error)
}
