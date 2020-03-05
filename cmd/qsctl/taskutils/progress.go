package taskutils

import (
	"sync"
	"time"

	"github.com/qingstor/noah/pkg/progress"
	"github.com/vbauerster/mpb/v4"
	"github.com/vbauerster/mpb/v4/decor"
)

// pBar is the struct contains the mark of bar finished, and the bar.
// Add the finished flag for cannot mark a bar as finished only with
// the chan coming state.
// If the bar's done equals to the total, but the incoming state is
// the finish state, we cannot judge whether this state is the first
// finish state.
// After adding the finished flag, the first finish state will set this
// to true, so that next finish state will be ignored.
type pBar struct {
	finished bool
	bar      *mpb.Bar
}

// wg is the global wait group used for multi-progress bar
// use pointer to keep it not copied
var wg *sync.WaitGroup

// pbPool is the multi-progress bar pool
var pbPool *mpb.Progress

// pbGroup is the local pBar group
// It takes taskID as the key, pointer to the pBar as value.
// So that every state will modify its relevant pBar.
var pbGroup map[string]*pBar

// sigChan is the channel to notify data progress channel to close.
var sigChan chan struct{}

func init() {
	wg = new(sync.WaitGroup)
	pbPool = mpb.New(mpb.WithWaitGroup(wg))
	pbGroup = make(map[string]*pBar)
	sigChan = make(chan struct{})
}

// StartProgress start to get state from state center.
// Use progress.Start to start a dataChan to get stateCenter from noah.
// The stateCenter is a map with taskID as key and its state as value.
// So we range the stateCenter and update relevant bar's progress.
func StartProgress(d time.Duration) error {
	dataChan := progress.Start(d)
	startTime := time.Now()
readChannel:
	for {
		select {
		case stateCenter := <-dataChan:
			for taskID, state := range stateCenter {
				pbar, ok := pbGroup[taskID]
				// bar already exists
				if ok {
					bar := pbar.GetBar()
					// if bar already finished, jump over
					if pbar.Finished() {
						continue
					}
					// change the bar attr
					bar.SetTotal(state.Total, false)
					bar.SetCurrent(state.Done, time.Since(startTime))
					// if this state is finish state, mark bar as finished
					if state.Finished() {
						pbar.MarkFinished()
						wg.Done()
					}
				} else { // bar not exists, create a new one.
					wg.Add(1)
					bar := pbPool.AddBar(state.Total,
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
					bar.SetCurrent(state.Done, time.Since(startTime))
					pbGroup[taskID] = &pBar{bar: bar}
				}
			}
		case <-sigChan:
			progress.End()
			break readChannel
		}
	}
	return nil
}

// WaitProgress wait the progress bar to complete
func WaitProgress() {
	pbPool.Wait()
}

// FinishProgress finish the progress bar and close the progress center
func FinishProgress() {
	close(sigChan)
}

// Finished is the flag of whether a pBar is finished
func (b pBar) Finished() bool {
	return b.finished
}

// MarkFinished mark a pBar as Finished
func (b *pBar) MarkFinished() {
	b.finished = true
}

// GetBar get the surrounded pointer to mpb.Bar
func (b pBar) GetBar() *mpb.Bar {
	return b.bar
}
