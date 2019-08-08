package parallel

import (
	"context"
	"sync"
	"time"
)


// Task run by ConcurrentGo
// ConcurrentGo will return immediately after all unignorable tasks done
// CancelFun will be invoked when ConcurrentGo return. It's always context's cancel function.
type Task struct {
	Func func()
	Ignorable bool
	CancelFunc func()
}

// ConcurrentGo run your function concurrently
// ConcurrentGo give up when ctx.Done() if ctx != nil
// timeout set timeout for run given task
// return done or timeout flags according to given tasks
func ConcurrentGo(ctx context.Context, timeout time.Duration, tasks ...Task) []bool{
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
		if tasks[i].CancelFunc != nil {
			defer tasks[i].CancelFunc()
		}

		go func(i int) {
			tasks[i].Func()
			if ctx.Err() == nil {
				taskFlags[i] = true
			}
			if !tasks[i].Ignorable {
				wg.Done()
			}
		}(i)
	}

	go func() {
		wg.Wait()
		cancelFun()
	}()

	<- ctx.Done()
	return taskFlags
}
