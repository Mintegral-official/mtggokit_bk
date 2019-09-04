package streamer

import (
	"github.com/Mintegral-official/mtggokit/bifrost/container"
	"github.com/easierway/concurrent_map"
	"github.com/pkg/errors"
	"strings"
)

type DefaultTextParser struct {
}

func (*DefaultTextParser) Parse(data []byte, userData interface{}) (container.DataMode, container.MapKey, interface{}, error) {
	s := string(data)
	items := strings.SplitN(s, "\t", 2)
	if len(items) != 2 {
		return container.DataModeAdd, nil, nil, errors.New("items len is not 2, item[" + s + "]")
	}
	return container.DataModeAdd, concurrent_map.StrKey(items[0]), items[1], nil
}
