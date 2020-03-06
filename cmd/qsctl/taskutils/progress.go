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

// pBarGroup is the progress bar group struct
// it contains available props for bar display.
// bars is the local pBar group,
// It takes taskID as the key, pointer to the pBar as value.
// So that every state will modify its relevant pBar.
type pBarGroup struct {
	sync.Mutex
	activeBarCount int // activeBarCount is the amount of not finished bar
	bars           map[string]*pBar
}

// wg is the global wait group used for multi-progress bar
// use pointer to keep it not copied
var wg *sync.WaitGroup

// pbPool is the multi-progress bar pool
var pbPool *mpb.Progress

// pbGroup is the local pBar group
// It takes taskID as the key, pointer to the pBar as value.
// So that every state will modify its relevant pBar.
var pbGroup *pBarGroup

// sigChan is the channel to notify data progress channel to close.
var sigChan chan struct{}

func init() {
	wg = new(sync.WaitGroup)
	pbPool = mpb.New(mpb.WithWaitGroup(wg))
	pbGroup = &pBarGroup{
		bars: make(map[string]*pBar),
	}
	sigChan = make(chan struct{})
}

// StartProgress start to get state from state center.
// d is the duration time between two data,
// maxBarCount is the max count of bar displayed.
// Use progress.Start to start a dataChan to get stateCenter from noah.
// The stateCenter is a map with taskID as key and its state as value.
// So we range the stateCenter and update relevant bar's progress.
func StartProgress(d time.Duration, maxBarCount int) error {
	dataChan := progress.Start(d)
	startTime := time.Now()
readChannel:
	for {
		select {
		case stateCenter := <-dataChan:
			for taskID, state := range stateCenter {
				pbar, ok := pbGroup.GetPBarByID(taskID)
				// bar already exists and pbar not nil
				// set pbar to nil means this state is received, but not add into pbPool
				if ok && pbar != nil {
					bar := pbar.GetBar()
					// if bar already finished, jump over
					if pbar.Finished() {
						continue
					}
					// change the bar attr
					bar.SetTotal(state.Total, false)
					bar.SetCurrent(state.Done, time.Since(startTime))
					// if this state is finish state, mark bar as finished, dec the active bar amount
					if state.Finished() {
						pbar.MarkFinished()
						pbGroup.DecActive()
						wg.Done()
					}
				} else {
					// if bar not exist, means the task state is the first time get
					// create a new pbar with nil in pbGroup to take the seat,
					// but do not add it into pbPool to display
					if !ok {
						wg.Add(1)
						pbGroup.SetPBarByID(taskID, nil)
					}
					// if active bar already beyond the max count, continue to next
					if pbGroup.GetActiveCount() >= maxBarCount {
						continue
					}
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
						mpb.BarClearOnComplete(),
					)
					bar.SetCurrent(state.Done, time.Since(startTime))
					pbGroup.SetPBarByID(taskID, &pBar{bar: bar})
					pbGroup.IncActive()
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

// GetPBarByID returns the pbar's pointer with given taskID
func (pg *pBarGroup) GetPBarByID(id string) (pbar *pBar, ok bool) {
	pg.Lock()
	defer pg.Unlock()
	pbar, ok = pg.bars[id]
	return
}

// SetPBarByID set the pBar with given taskID into pBarGroup
func (pg *pBarGroup) SetPBarByID(id string, pbar *pBar) {
	pg.Lock()
	defer pg.Unlock()
	pg.bars[id] = pbar
}

// GetActiveCount returns how many active bars in the group
func (pg *pBarGroup) GetActiveCount() int {
	pg.Lock()
	defer pg.Unlock()
	return pg.activeBarCount
}

// IncActive add the active bar count by one
func (pg *pBarGroup) IncActive() {
	pg.Lock()
	defer pg.Unlock()
	pg.activeBarCount++
}

// DecActive minus the active bar count by one
func (pg *pBarGroup) DecActive() {
	pg.Lock()
	defer pg.Unlock()
	pg.activeBarCount--
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
