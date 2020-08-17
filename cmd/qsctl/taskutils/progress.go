package taskutils

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/qingstor/noah/pkg/progress"
	"github.com/vbauerster/mpb/v5"
	"github.com/vbauerster/mpb/v5/decor"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	widStatus          = 20
	widStat            = 20
	widBarDefault      = 40
	widTerminalDefault = 120
)

const (
	pbNotExist pBarStatus = iota
	pbNotShow
	pbShown
	pbFinished
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

// ClearFunc is the alias of func to clear PbHandler
type ClearFunc func()

// PbHandler is used to handle progress bar
type PbHandler struct {
	// wg is the wait group used for multi-progress bar
	// use pointer to keep it not copied
	wg *sync.WaitGroup
	// pbPool is the multi-progress bar pool
	pbPool *mpb.Progress
	// pbGroup is the local pBar group
	// It takes taskID as the key, pointer to the pBar as value.
	// So that every state will modify its relevant pBar.
	pbGroup *pBarGroup
	// closeSig is the channel to notify data progress channel close.
	closeSig chan struct{}
	// nameFilter if filter func to handle bar tip info
	// for now, it used to truncate file name if it is too long to display
	nameFilter func(string) string
}

// NewHandler new a PbHandler to handle progress bar, and a cancelFunc to cancel the handler
func NewHandler(c context.Context) (*PbHandler, ClearFunc) {
	terminalWidth, _, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		terminalWidth = widTerminalDefault
	}
	// nameWidth and barWidth is the style width for progress bar
	nameWidth, barWidth := calBarSize(terminalWidth)
	wg := new(sync.WaitGroup)
	// ctx, cancel := context.WithCancel(c)
	closeChan := make(chan struct{})
	pbGroup := &pBarGroup{
		bars: make(map[string]*pBar),
	}
	pbPool := mpb.NewWithContext(c, mpb.WithWaitGroup(wg), mpb.WithWidth(barWidth))
	return &PbHandler{
			wg:       wg,
			pbPool:   pbPool,
			pbGroup:  pbGroup,
			closeSig: closeChan,
			nameFilter: func(s string) string {
				return truncateBefore(s, nameWidth)
			},
		}, func() {
			close(closeChan)
			for _, pbar := range pbGroup.bars {
				pbar.GetBar().Abort(true)
			}
			progress.ClearStates()
		}
}

// StartProgress start to get state from state center.
// d is the duration time between two data,
// Start a ticker to get data from noah periodically.
// The data from noah is a map with taskID as key and its state as value.
// So we range the data and update relevant bar's progress.
func (h *PbHandler) StartProgress(d time.Duration) {
	tc := time.NewTicker(d)
	for {
		select {
		case <-tc.C:
			for taskID, state := range progress.GetStates() {
				pbar := h.pbGroup.GetPBarByID(taskID)
				// if pbar already finished, skip directly
				if pbar.Finished() {
					continue
				}
				// if bar state is list type, skip directly
				if state.IsListType() {
					continue
				}
				// if pbar is shown, update progress
				if pbar.Shown() {
					bar := pbar.GetBar()
					// bar.SetTotal(state.Total, false)
					bar.SetCurrent(state.Done)
					bar.DecoratorEwmaUpdate(d)
					// if this state is finish state, mark bar as finished, dec the active bar amount
					if state.Finished() {
						pbar.MarkFinished()
						// spinner (list type state) is not counted
						if !state.IsListType() {
							h.pbGroup.DecActive()
						}
						h.wg.Done()
					}
					continue
				}

				bar := h.addBarByState(state)
				bar.SetCurrent(state.Done)
				bar.DecoratorEwmaUpdate(d)
				h.wg.Add(1)
				h.pbGroup.SetPBarByID(taskID, &pBar{status: pbShown, bar: bar})
				h.pbGroup.IncActive()
			}
		case <-h.closeSig:
			tc.Stop()
			return
		}
	}
}

// WaitProgress wait the progress bar to complete
func (h *PbHandler) WaitProgress() {
	h.pbPool.Wait()
}

// GetPBarByID returns the pbar's pointer with given taskID
func (pg *pBarGroup) GetPBarByID(id string) *pBar {
	pbar, ok := pg.bars[id]
	if !ok {
		pbar = &pBar{status: pbNotExist}
	}
	return pbar
}

// SetPBarByID set the pBar with given taskID into pBarGroup
func (pg *pBarGroup) SetPBarByID(id string, pbar *pBar) {
	pg.bars[id] = pbar
}

// GetActiveCount returns how many active bars in the group
func (pg *pBarGroup) GetActiveCount() int {
	return pg.activeBarCount
}

// IncActive add the active bar count by one
func (pg *pBarGroup) IncActive() {
	pg.activeBarCount++
}

// DecActive minus the active bar count by one
func (pg *pBarGroup) DecActive() {
	pg.activeBarCount--
}

// addSpinnerByState add a spinner to pool, and return it for advanced operation
func (h *PbHandler) addSpinnerByState(state progress.State) (bar *mpb.Bar) {
	bar = h.pbPool.Add(state.Total, mpb.NewSpinnerFiller([]string{".", "..", "..."}, mpb.SpinnerOnLeft),
		mpb.PrependDecorators(
			decor.Name(state.Status, decor.WCSyncSpaceR),
			decor.Name(h.nameFilter(state.Name), decor.WCSyncSpaceR),
		),
		mpb.AppendDecorators(
			decor.OnComplete(decor.Name(""), "done"),
		),
		mpb.BarRemoveOnComplete(),
	)
	return
}

// addBarByState add a bar to pool, and return it for advanced operation
func (h *PbHandler) addBarByState(state progress.State) (bar *mpb.Bar) {
	bar = h.pbPool.AddBar(state.Total, mpb.BarStyle("[=>-|"),
		mpb.PrependDecorators(
			decor.Name(state.Status, decor.WCSyncSpaceR),
			decor.Name(h.nameFilter(state.Name), decor.WCSyncSpaceR),
		),
		mpb.AppendDecorators(
			decor.EwmaSpeed(decor.UnitKiB, "% .2f", 0, decor.WCSyncSpace),
			// decor.EwmaETA(decor.ET_STYLE_GO, 0, decor.WCSyncSpace),
			decor.Name(" ] "),
			decor.OnComplete(
				decor.CountersKibiByte("% .2f / % .2f", decor.WCSyncSpace), "done",
				// decor.Percentage(decor.WCSyncSpace),
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
	b.GetBar().Abort(true)
}

// GetBar get the surrounded pointer to mpb.Bar
func (b *pBar) GetBar() *mpb.Bar {
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
