package streamer

import (
	"github.com/Mintegral-official/mtggokit/bifrost/container"
	"github.com/easierway/concurrent_map"
	"strings"
)

type DefaultTextParser struct {
}

func (*DefaultTextParser) Parse(data []byte, userData interface{}) []ParserResult {
	s := string(data)
	items := strings.SplitN(s, "\t", 2)
	if len(items) != 2 {
		return nil
	}
	return []ParserResult{{container.DataModeAdd, concurrent_map.StrKey(items[0]), items[1], nil}}
}
