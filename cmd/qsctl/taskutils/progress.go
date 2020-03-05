package taskutils

import (
	"sync"
	"time"

	"github.com/qingstor/noah/pkg/progress"
	"github.com/vbauerster/mpb/v4"
	"github.com/vbauerster/mpb/v4/decor"
)

var pbPool *mpb.Progress
var pbGroup sync.Map
var wg sync.WaitGroup
var sigChan chan struct{}

func init() {
	sigChan = make(chan struct{})
	pbPool = mpb.New(mpb.WithWaitGroup(&wg))
}

// StartProgress start to get state from state center
func StartProgress(d time.Duration) error {
	data := progress.Start(d)
	startTime := time.Now()
readChannel:
	for {
		select {
		case stateCenter := <-data:
			stateCenter.Range(func(taskID, v interface{}) bool {
				var pbar *mpb.Bar
				state := v.(progress.State)
				bar, ok := pbGroup.Load(taskID)
				// bar already exists
				if ok {
					pbar = bar.(*mpb.Bar)
					// fmt.Println(pbar.Get(), "load, id:", taskID, "state:", state)
					pbar.SetTotal(state.Total, false)
					pbar.SetCurrent(state.Done, time.Since(startTime))
				} else {
					wg.Add(1)
					pbar = pbPool.AddBar(state.Total,
						mpb.PrependDecorators(
							decor.Name(state.TaskName, decor.WCSyncSpaceR),
							decor.NewElapsed(decor.ET_STYLE_HHMMSS, startTime),
						),
						mpb.AppendDecorators(
							decor.OnComplete(
								decor.Percentage(decor.WCSyncSpace), "done",
							),
						),
					)
					pbar.SetCurrent(state.Done, time.Since(startTime))
					// fmt.Println("add, id:", taskID, "state:", state)
					pbGroup.Store(taskID, pbar)
				}

				if state.Finished() {
					// pbar.SetTotal(state.Total, true)
					wg.Done()
				}
				return true
			})
		case <-sigChan:
			progress.End()
			break readChannel
		}
	}
	return nil
}

// WaitProgress wait progress
func WaitProgress() {
	pbPool.Wait()
}

// FinishProgress finish the progress bar and close the progress center
func FinishProgress() {
	close(sigChan)
}
