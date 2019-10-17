package container

import (
	"github.com/pkg/errors"
	"github.com/smartystreets/goconvey/convey"
	"strconv"
	"strings"
	"testing"
)

type TestIntIntDataIter struct {
	current int
	data    []string
}

func NewTestIntIntDataIter(data []string) *TestIntDataIter {
	return &TestIntDataIter{
		current: 0,
		data:    data,
	}
}

func (this *TestIntIntDataIter) HasNext() bool {
	return this.current < len(this.data)
}

func (this *TestIntIntDataIter) Next() (DataMode, MapKey, interface{}, error) {
	defer func() { this.current++ }()
	if this.current >= len(this.data) {
		return DataModeAdd, nil, nil, errors.New("current index is error")
	}
	s := this.data[this.current]
	items := strings.SplitN(s, "\t", 2)
	if len(items) != 2 {
		return DataModeAdd, nil, nil, errors.New("items len is not 2, item[" + s + "]")
	}
	n, e := strconv.ParseInt(items[0], 10, 64)
	if e != nil {
		return DataModeAdd, nil, nil, errors.New("parse key error, not an number")
	}
	return DataModeAdd, I64Key(n), items[1], nil
}

func TestBufferedKListContainer(t *testing.T) {
	convey.Convey("Test BufferedMapContainer Get", t, func() {
		bm := BufferedKListContainer{}
		convey.So(bm.LoadBase(NewTestDataIter([]string{})), convey.ShouldBeNil)
		convey.So(bm.ErrorNum, convey.ShouldEqual, 0)
		convey.So(len(*bm.innerData), convey.ShouldEqual, 0)
	})

	convey.Convey("Test BufferedMapContainer Get", t, func() {
		bm := BufferedKListContainer{}
		convey.So(bm.LoadBase(NewTestDataIter([]string{
			"1\t2",
			"a\tb",
			"a\tcc",
		})), convey.ShouldBeNil)
		convey.So(bm.ErrorNum, convey.ShouldEqual, 0)
		convey.So(len(*bm.innerData), convey.ShouldEqual, 2)
		{
			v, e := bm.Get(StrKey("1"))
			convey.So(e, convey.ShouldBeNil)
			arr, ok := v.([]interface{})
			convey.So(ok, convey.ShouldBeTrue)
			convey.So(len(arr), convey.ShouldEqual, 1)
			convey.So(arr[0], convey.ShouldEqual, "2")
		}
		{
			v, e := bm.Get(StrKey("a"))
			convey.So(e, convey.ShouldBeNil)
			arr, ok := v.([]interface{})
			convey.So(ok, convey.ShouldBeTrue)
			convey.So(len(arr), convey.ShouldEqual, 2)
			convey.So(arr[0], convey.ShouldEqual, "b")
			convey.So(arr[1], convey.ShouldEqual, "cc")
		}
	})
}
