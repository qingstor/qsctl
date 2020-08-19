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

// predefined width, see more detail in comments of calBarSize()
const (
	widStatus          = 20
	widStat            = 40
	widBarDefault      = 40
	widTerminalDefault = 120
)

// barGroup is the progress bar group struct
// it contains available props for bar display.
// bars is the local bar group,
// It takes taskID as the key, pointer to the mpb.Bar as value.
// So that every state will modify its relevant pBar.
type barGroup struct {
	sync.Mutex
	bars map[string]*mpb.Bar
}

// ClearFunc is the alias of func to clear PbHandler
type ClearFunc func()

// PbHandler is used to handle progress bar
type PbHandler struct {
	// pbPool is the multi-progress bar pool
	pbPool *mpb.Progress
	// bGroup contains every bar
	bGroup *barGroup
	// closeSig is the channel to notify data progress channel close.
	closeSig chan struct{}
	// nameFilter is the filter func to handle bar tip info
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
	closeChan := make(chan struct{})
	pbGroup := &barGroup{
		bars: make(map[string]*mpb.Bar),
	}

	ctx, cancel := context.WithCancel(c) // conduct new context to make sure cancel pbPool's context
	pbPool := mpb.NewWithContext(ctx, mpb.WithWidth(barWidth))
	ph := &PbHandler{
		pbPool:   pbPool,
		bGroup:   pbGroup,
		closeSig: closeChan,
		nameFilter: func(s string) string {
			return truncateBefore(s, nameWidth)
		},
	}
	return ph, func() {
		// notify to stop getting stat from noah
		close(closeChan)
		// make sure cancel the progress bar's context
		cancel()
		// free local bars and group
		pbGroup.Lock()
		defer pbGroup.Unlock()
		for id, bar := range pbGroup.bars {
			bar.Abort(true)
			delete(pbGroup.bars, id)
		}
		// clear noah stats to for re-using
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
			h.bGroup.Lock()
			for taskID, state := range progress.GetStates() {
				// if bar state is list type, skip directly
				// because we do not show list progress now
				if state.IsListType() {
					continue
				}

				bar := h.bGroup.GetBarByID(taskID)
				// if bar is showing, only update its status
				if bar != nil {
					bar.SetCurrent(state.Done)
					bar.DecoratorEwmaUpdate(d)
					continue
				}
				// if bar is not showing
				// if the stat is finished, which means bar already removed, just ignore
				if state.Finished() {
					continue
				}
				// add bar into group
				b := h.addBarByState(state)
				b.SetCurrent(state.Done)
				b.DecoratorEwmaUpdate(d)
				h.bGroup.SetBarByID(taskID, b)
			}
			h.bGroup.Unlock()
		case <-h.closeSig:
			tc.Stop()
			h.WaitProgress() // have to wait all bars shutdown
			return
		}
	}
}

// WaitProgress wait the progress bar to complete
func (h *PbHandler) WaitProgress() {
	h.pbPool.Wait()
}

// GetBarByID returns the bar's pointer with given taskID
// if not exist, return nil
func (pg *barGroup) GetBarByID(id string) *mpb.Bar {
	bar, ok := pg.bars[id]
	if !ok {
		return nil
	}
	return bar
}

// SetBarByID set the bar with given taskID into barGroup
func (pg *barGroup) SetBarByID(id string, bar *mpb.Bar) {
	pg.bars[id] = bar
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
// |   status    |  name  |     bar      |                stat                      |
// [copy parts: ][abc.jpg][[===>-------|][ 1020.01 KiB/s ] 1020.01 MiB / 1023.99 MiB]
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
