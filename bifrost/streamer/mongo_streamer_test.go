package streamer

import (
	"github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

func TestNewMongoStreamer(t *testing.T) {
	convey.Convey("", t, func() {
		ms, err := NewMongoStreamer(&MongoStreamerCfg{
			Name:           "mongo_test",
			UpdatMode:      Dynamic,
			IncInterval:    60,
			IsSync:         true,
			URI:            "mongodb://127.0.0.1:21017",
			ConnectTimeout: 100,
			ReadTimeout:    20,
			BaseParser:     &DefaultTextParser{},
			IncParser:      &DefaultTextParser{},
			BaseQuery:      "mongo base query",
			IncQuery:       "mongo inc query",
			UserData:       "user defined data",
			OnBeforeInc: func(userData interface{}) interface{} {
				return "nfew inc base query"
			},
		})
		convey.So(err, convey.ShouldNotBeNil)
		convey.So(strings.Contains(err.Error(), "context deadline exceeded"), convey.ShouldBeTrue)
		convey.So(ms, convey.ShouldBeNil)
	})
}
