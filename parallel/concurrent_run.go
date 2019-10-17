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
	finished := make([]bool, len(tasks))
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
		localI := i
		//run task
		ants.Submit(func() {
			tasks[localI].Func()
			if ctx.Err() == nil {
				finished[localI] = true
			}
			if !tasks[localI].Ignorable {
				wg.Done()
			}
		})
		//register cancel function
		if tasks[localI].CancelFunc != nil {
			ants.Submit(func(){
				<- ctx.Done()
				if !finished[localI] {
					tasks[localI].CancelFunc()
				}
			})

		}
	}

	ants.Submit( func() {
		wg.Wait()
		cancelFun()
	})

	<- ctx.Done()
	return finished
}
