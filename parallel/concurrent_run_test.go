package parallel

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestConcurrentRun(t *testing.T) {

	Convey("任务超时", t, func() {
		hasDone := ConcurrentRun(nil, time.Millisecond*10, Task{Ignorable: false, Func: func() { time.Sleep(time.Second) }})
		So(hasDone[0], ShouldEqual, false)
	})

	Convey("任务完成", t, func() {
		hasDone := ConcurrentRun(nil, time.Millisecond*10, Task{Ignorable: false, Func: func() {}})
		So(hasDone[0], ShouldEqual, true)
	})

	Convey("ignorable任务被取消, 对应cancelFunc被执行", t, func() {
		var canceledOne, canceledTwo bool
		hasDone := ConcurrentRun(nil, time.Second, Task{Ignorable: false, Func: func() {}, CancelFunc: func() { canceledOne = true }},
			Task{Ignorable: true, Func: func() { time.Sleep(time.Second * 2) }, CancelFunc: func() { canceledTwo = true }})
		So(hasDone[0], ShouldEqual, true)
		So(hasDone[1], ShouldEqual, false)
		time.Sleep(time.Second * 1)
		So(canceledOne, ShouldEqual, false)
		So(canceledTwo, ShouldEqual, true)
	})
}

func BenchmarkConcurrentRun(b *testing.B) {
	for i := 0; i < b.N; i++ { //use b.N for looping
		ConcurrentRun(nil, time.Millisecond*10, Task{Ignorable: false, Func: func() {
			var sum int
			for j := 0; j < 10000; j++ {
				sum += j
			}
		}})
	}
}

func BenchmarkNormalGo(b *testing.B) {
	f := func() {
		var sum int
		for j := 0; j < 10000; j++ {
			sum += j
		}
	}
	for i := 0; i < b.N; i++ { //use b.N for looping
		f()
	}
}
