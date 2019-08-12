package parallel

import (
	"context"
	"github.com/panjf2000/ants"
	"sync"
	"time"
)


// Task run by ConcurrentRun
// ConcurrentRun will return immediately after all unignorable tasks done
// CancelFun will be invoked when this task overtime. It's always context's cancel function.
type Task struct {
	Func func()
	Ignorable bool
	CancelFunc func()
}

// ConcurrentRun run your function concurrently
// ConcurrentRun give up when ctx.Done() if ctx != nil
// timeout set timeout for run given task
// return done or timeout flags according to given tasks
func ConcurrentRun(ctx context.Context, timeout time.Duration, tasks ...Task) []bool{
	taskFlags := make([]bool, len(tasks))
	if ctx == nil {
		ctx = context.Background()
	}

	ctx, cancelFun := context.WithTimeout(ctx, timeout)
	defer cancelFun()

	var wg sync.WaitGroup
	for i := range tasks {
		if !tasks[i].Ignorable {
			wg.Add(1)
		}
		uniq_i := i
		ants.Submit(func() {
			tasks[uniq_i].Func()
			if ctx.Err() == nil {
				taskFlags[uniq_i] = true
			} else {
				if tasks[uniq_i].CancelFunc != nil {
					tasks[uniq_i].CancelFunc()
				}
			}
			if !tasks[uniq_i].Ignorable {
				wg.Done()
			}
		})
	}

	ants.Submit( func() {
		wg.Wait()
		cancelFun()
	})

	<- ctx.Done()
	return taskFlags
}
