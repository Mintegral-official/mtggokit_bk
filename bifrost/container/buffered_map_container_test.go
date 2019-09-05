package container

import (
	"github.com/pkg/errors"
	"github.com/smartystreets/goconvey/convey"
	"strconv"
	"strings"
	"testing"
)

type TestDataIter struct {
	current int
	data    []string
}

func NewTestDataIter(data []string) *TestDataIter {
	return &TestDataIter{
		current: 0,
		data:    data,
	}
}

func (this *TestDataIter) HasNext() bool {
	return this.current < len(this.data)
}

func (this *TestDataIter) Next() (DataMode, MapKey, interface{}, error) {
	defer func() { this.current++ }()
	if this.current >= len(this.data) {
		return DataModeAdd, nil, nil, errors.New("current index is error")
	}
	s := this.data[this.current]
	items := strings.SplitN(s, "\t", 2)
	if len(items) != 2 {
		return DataModeAdd, nil, nil, errors.New("items len is not 2, item[" + s + "]")
	}
	return DataModeAdd, StrKey(items[0]), items[1], nil
}

type TestIntDataIter struct {
	current int
	data    []string
}

func NewTestIntDataIter(data []string) *TestIntDataIter {
	return &TestIntDataIter{
		current: 0,
		data:    data,
	}
}

func (this *TestIntDataIter) HasNext() bool {
	return this.current < len(this.data)
}

func (this *TestIntDataIter) Next() (DataMode, MapKey, interface{}, error) {
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

func TestBufferedMap_Get(t *testing.T) {
	convey.Convey("Test BufferedMapContainer Get", t, func() {
		bm := BufferedMapContainer{}
		convey.So(bm.LoadBase(NewTestDataIter([]string{})), convey.ShouldBeNil)
		convey.So(bm.errorNum, convey.ShouldEqual, 0)
		convey.So(len(*bm.innerData), convey.ShouldEqual, 0)
	})

	convey.Convey("Test BufferedMapContainer Get", t, func() {
		bm := BufferedMapContainer{}
		convey.So(bm.LoadBase(NewTestDataIter([]string{
			"1\t2",
			"a\tb",
		})), convey.ShouldBeNil)
		convey.So(bm.errorNum, convey.ShouldEqual, 0)
		convey.So(len(*bm.innerData), convey.ShouldEqual, 2)
		v, e := bm.Get(StrKey("1"))
		convey.So(e, convey.ShouldBeNil)
		convey.So(v, convey.ShouldEqual, "2")

		v, e = bm.Get(StrKey("a"))
		convey.So(e, convey.ShouldBeNil)
		convey.So(v, convey.ShouldEqual, "b")
	})

	convey.Convey("Test BufferedMapContainer Get", t, func() {
		bm := BufferedMapContainer{}
		convey.So(bm.LoadBase(NewTestIntDataIter([]string{
			"1\t2",
			"4\tb",
		})), convey.ShouldBeNil)
		convey.So(bm.errorNum, convey.ShouldEqual, 0)
		convey.So(len(*bm.innerData), convey.ShouldEqual, 2)
		v, e := bm.Get(I64Key(1))
		convey.So(e, convey.ShouldBeNil)
		convey.So(v, convey.ShouldEqual, "2")

		v, e = bm.Get(I64Key(4))
		convey.So(e, convey.ShouldBeNil)
		convey.So(v, convey.ShouldEqual, "b")
	})
}
