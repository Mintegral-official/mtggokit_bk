package container

import (
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestBufferedMap_Get(t *testing.T) {
	convey.Convey("Test BufferedMap Get", t, func() {
		bm := BufferedMap{}
		bm.innerData[I64Key(1)] = 5
	})
}
