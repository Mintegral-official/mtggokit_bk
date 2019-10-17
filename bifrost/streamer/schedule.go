package streamer

import (
	"container/heap"
	"context"
	"time"
)

type SchedInfo struct {
	TimeInterval int
}

type SchedUnit struct {
	name     string
	streamer Streamer
	deadline int
	index    int
}

type Sched []*SchedUnit

func (s Sched) Len() int { return len(s) }

func (s Sched) Less(i, j int) bool {
	return s[i].deadline < s[j].deadline
}

func (s Sched) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
	s[i].index = i
	s[j].index = j
}

func (s *Sched) Push(x interface{}) {
	item := x.(*SchedUnit)
	n := len(*s)
	item.index = n
	*s = append(*s, item)
}

func (s *Sched) Top() *SchedUnit {
	if len(*s) == 0 {
		return nil
	}
	return (*s)[0]
}

func (s *Sched) Pop() interface{} {
	old := *s
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*s = old[0 : n-1]
	return item
}

func (s *Sched) AddStreamer(name string, dataStreamer Streamer) {
	s.Push(&SchedUnit{
		name:     name,
		streamer: dataStreamer,
		deadline: int(time.Now().Unix()) + dataStreamer.GetSchedInfo().TimeInterval,
	})
}

func (s *Sched) Schedule(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if time.Now().Unix() < int64(s.Top().deadline) {
				time.Sleep(time.Second * time.Duration(int(time.Now().Unix())-s.Top().deadline))
				continue
			}
			x := heap.Pop(s)
			su := x.(*SchedUnit)
			//TODO handler error
			_ = su.streamer.UpdateData(context.Background())
			su.deadline += su.streamer.GetSchedInfo().TimeInterval
			heap.Push(s, su)
			time.Sleep(time.Second * time.Duration(s.Top().deadline-int(time.Now().Unix())))
		}
	}
}
