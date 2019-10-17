package streamer

import (
	"github.com/Mintegral-official/mtggokit/bifrost/container"
)

type ParserResult struct {
	DataMode container.DataMode
	Key      container.MapKey
	Value    interface{}
	Err      error
}

type DataParser interface {
	Parse([]byte, interface{}) []ParserResult
}
