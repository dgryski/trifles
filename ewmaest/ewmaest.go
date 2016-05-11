package ewmaest

import (
	"fmt"
	"time"
)

type EWMAEstimate struct {
	TotalItems int
	Alpha      float64

	completedItems int

	timePerItemEstimate time.Duration
	blockStartTime      time.Time
}

func New(total int) *EWMAEstimate {
	return &EWMAEstimate{
		TotalItems: total,
		Alpha:      0.25,
	}
}

func (e *EWMAEstimate) StartBlock() {
	e.blockStartTime = time.Now()
}

func (e *EWMAEstimate) CompleteItems(items int) time.Duration {
	lastBlockTime := time.Since(e.blockStartTime)
	lastItemEstimate := float64(lastBlockTime) / float64(items)
	e.timePerItemEstimate = time.Duration((e.Alpha * lastItemEstimate) + (1-e.Alpha)*float64(e.timePerItemEstimate))

	e.completedItems += items

	remainingItems := e.TotalItems - e.completedItems
	remainingTime := time.Duration(remainingItems) * e.timePerItemEstimate

	return remainingTime
}

func (e *EWMAEstimate) Progress(remainingTime time.Duration) string {
	return fmt.Sprintf("Completed %d/%d (%02d %%), ETA %v (%v)", e.completedItems, e.TotalItems, int(100*float64(e.completedItems)/float64(e.TotalItems)), time.Duration(remainingTime.Seconds())*time.Second, time.Now().Add(remainingTime).Format(time.Stamp))
}
