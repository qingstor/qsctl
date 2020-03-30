package taskutils

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/qingstor/noah/pkg/progress"
	"github.com/vbauerster/mpb/v4"
	"github.com/vbauerster/mpb/v4/decor"
	"golang.org/x/crypto/ssh/terminal"
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
	status pBarStatus
	bar    *mpb.Bar
}

type pBarStatus int

const (
	pbNotExist pBarStatus = iota
	pbNotShow
	pbShown
	pbFinished
)

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
var wg = new(sync.WaitGroup)

// ctx is used to cancel pbPool
var ctx, cancel = context.WithCancel(context.Background())

// pbPool is the multi-progress bar pool
var pbPool *mpb.Progress

// pbGroup is the local pBar group
// It takes taskID as the key, pointer to the pBar as value.
// So that every state will modify its relevant pBar.
var pbGroup = &pBarGroup{
	bars: make(map[string]*pBar),
}

// sigChan is the channel to notify data progress channel to close.
var sigChan = make(chan struct{})

// nameWidth and barWidth is the style width for progress bar
var nameWidth, barWidth int

const (
	widStatus          = 20
	widStat            = 20
	widBarDefault      = 40
	widTerminalDefault = 120
)

func init() {
	terminalWidth, _, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		terminalWidth = widTerminalDefault
	}
	nameWidth, barWidth = calBarSize(terminalWidth)
	pbPool = mpb.NewWithContext(ctx, mpb.WithWaitGroup(wg), mpb.WithWidth(barWidth))
}

// StartProgress start to get state from state center.
// d is the duration time between two data,
// maxBarCount is the max count of bar displayed.
// Start a ticker to get data from noah periodically.
// The data from noah is a map with taskID as key and its state as value.
// So we range the data and update relevant bar's progress.
func StartProgress(d time.Duration, maxBarCount int) {
	tc := time.NewTicker(d)
	for {
		select {
		case <-tc.C:
			for taskID, state := range progress.GetStates() {
				pbar := pbGroup.GetPBarByID(taskID)
				// if pbar already finished, skip directly
				if pbar.Finished() {
					continue
				}
				// if pbar is shown, update progress
				if pbar.Shown() {
					bar := pbar.GetBar()
					bar.SetTotal(state.Total, false)
					bar.SetCurrent(state.Done, d)
					// if this state is finish state, mark bar as finished, dec the active bar amount
					if state.Finished() {
						pbar.MarkFinished()
						// spinner (list type state) is not counted
						if !state.IsListType() {
							pbGroup.DecActive()
						}
						wg.Done()
					}
					continue
				}

				// if bar not exist, means the task state is the first time get
				// create a new pbar with status "not show" in pbGroup to take the seat,
				// but do not add it into pbPool to display
				if pbar.NotExist() {
					wg.Add(1)
					pbGroup.SetPBarByID(taskID, &pBar{status: pbNotShow})
				}
				// if bar state is list type, always show spinner and not add into count
				if state.IsListType() {
					bar := addSpinnerByState(state)
					bar.SetCurrent(state.Done, d)
					pbGroup.SetPBarByID(taskID, &pBar{status: pbShown, bar: bar})
					continue
				}

				// if active bar already beyond the max count, continue to next
				if pbGroup.GetActiveCount() >= maxBarCount {
					continue
				}

				bar := addBarByState(state)
				bar.SetCurrent(state.Done, d)
				pbGroup.SetPBarByID(taskID, &pBar{status: pbShown, bar: bar})
				pbGroup.IncActive()
			}
		case <-sigChan:
			return
		}
	}
}

// WaitProgress wait the progress bar to complete
func WaitProgress() {
	pbPool.Wait()
}

// FinishProgress finish the progress bar and close the progress center
func FinishProgress() {
	close(sigChan)
	cancel()
}

// GetPBarByID returns the pbar's pointer with given taskID
func (pg *pBarGroup) GetPBarByID(id string) *pBar {
	pg.Lock()
	defer pg.Unlock()
	pbar, ok := pg.bars[id]
	if !ok {
		pbar = &pBar{status: pbNotExist}
	}
	return pbar
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

// addSpinnerByState add a spinner to pool, and return it for advanced operation
func addSpinnerByState(state progress.State) (bar *mpb.Bar) {
	bar = pbPool.Add(state.Total, mpb.NewSpinnerFiller([]string{".", "..", "..."}, mpb.SpinnerOnLeft),
		mpb.PrependDecorators(
			decor.Name(state.Status, decor.WCSyncSpaceR),
			decor.Name(truncateBefore(state.Name, nameWidth), decor.WCSyncSpaceR),
		),
		mpb.AppendDecorators(
			decor.OnComplete(decor.Name(""), "done"),
		),
		mpb.BarRemoveOnComplete(),
	)
	return
}

// addBarByState add a bar to pool, and return it for advanced operation
func addBarByState(state progress.State) (bar *mpb.Bar) {
	bar = pbPool.AddBar(state.Total, mpb.BarStyle("[=>-|"),
		mpb.PrependDecorators(
			decor.Name(state.Status, decor.WCSyncSpaceR),
			decor.Name(truncateBefore(state.Name, nameWidth), decor.WCSyncSpaceR),
		),
		mpb.AppendDecorators(
			decor.EwmaETA(decor.ET_STYLE_GO, 0, decor.WCSyncSpace),
			decor.Name(" ] "),
			decor.OnComplete(
				decor.Percentage(decor.WCSyncSpace), "done",
			),
		),
		mpb.BarRemoveOnComplete(),
	)
	return
}

// Finished is the flag of whether a pBar is finished
func (b pBar) Finished() bool {
	return b.status == pbFinished
}

// Shown is the flag of whether a pBar is shown
func (b pBar) Shown() bool {
	return b.status == pbShown
}

// NotExist means the pBar not added and doesn't contain a bar
func (b pBar) NotExist() bool {
	return b.status == pbNotExist
}

// MarkFinished mark a pBar as Finished
func (b *pBar) MarkFinished() {
	b.status = pbFinished
}

// GetBar get the surrounded pointer to mpb.Bar
func (b pBar) GetBar() *mpb.Bar {
	return b.bar
}

// truncateBefore keeps the last l chars of s, use ... before
func truncateBefore(s string, l int) string {
	if len(s) <= l {
		return s
	}
	return "..." + s[len(s)-l:]
}

// calBarSize calculate the bar size by given full width (usually terminal width)
// first return width for progress name, second return width for bar
// A progress bar consists of "status", "name", "bar" and "stat".
// |   status    |  name  |     bar      |    stat   |
// [copy parts: ][abc.jpg][[===>-------|][ 20s ] 25%]
func calBarSize(fullWid int) (nameWid, barWid int) {
	// aviWid equals full width minus status and stat reserved
	aviWid := fullWid - widStatus - widStat
	// if it is more than 2*widBarDefault, keep bar as default width
	if aviWid >= 2*widBarDefault {
		barWid = widBarDefault
		nameWid = aviWid - barWid
		return
	}
	// otherwise, divide it equally into bar and name
	barWid = aviWid / 2
	nameWid = barWid
	return
}
