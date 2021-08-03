package adapter

import (
	"github.com/jedib0t/go-pretty/v6/progress"
	"github.com/stretchr/testify/mock"
	"io"
	"time"
)

// Tracker mock

type TrackerMock struct {
	mock.Mock
}

func (t *TrackerMock) Increment(value int64) {
	t.Called(value)
}

func (t *TrackerMock) MarkAsDone() {
	t.Called()
}

// Progress mock

type ProgressMock struct {
	mock.Mock
}

func (p *ProgressMock) SetAutoStop(autoStop bool) {
	p.Called(autoStop)
}

func (p *ProgressMock) AppendTracker(tracker *progress.Tracker) {
	p.Called(tracker)
}

func (p *ProgressMock) AppendTrackers(trackers []*progress.Tracker) {
	p.Called(trackers)
}

func (p *ProgressMock) IsRenderInProgress() bool {
	arg := p.Called()
	return arg.Bool(0)
}

func (p *ProgressMock) Length() int {
	arg := p.Called()
	return arg.Int(0)
}

func (p *ProgressMock) LengthActive() int {
	arg := p.Called()
	return arg.Int(0)
}

func (p *ProgressMock) LengthDone() int {
	arg := p.Called()
	return arg.Int(0)
}

func (p *ProgressMock) LengthInQueue() int {
	arg := p.Called()
	return arg.Int(0)
}

func (p *ProgressMock) SetMessageWidth(width int) {
	p.Called(width)
}

func (p *ProgressMock) SetNumTrackersExpected(numTrackers int) {
	p.Called(numTrackers)
}

func (p *ProgressMock) SetOutputWriter(output io.Writer) {
	p.Called(output)
}

func (p *ProgressMock) SetSortBy(sortBy progress.SortBy) {
	p.Called(sortBy)
}

func (p *ProgressMock) SetStyle(style progress.Style) {
	p.Called(style)
}

func (p *ProgressMock) SetTrackerLength(length int) {
	p.Called(length)
}

func (p *ProgressMock) SetTrackerPosition(position progress.Position) {
	p.Called(position)
}

func (p *ProgressMock) ShowETA(show bool) {
	p.Called(show)
}

func (p *ProgressMock) ShowOverallTracker(show bool) {
	p.Called(show)
}

func (p *ProgressMock) ShowPercentage(show bool) {
	p.Called(show)
}

func (p *ProgressMock) ShowTime(show bool) {
	p.Called(show)
}

func (p *ProgressMock) ShowTracker(show bool) {
	p.Called(show)
}

func (p *ProgressMock) ShowValue(show bool) {
	p.Called(show)
}

func (p *ProgressMock) SetUpdateFrequency(frequency time.Duration) {
	p.Called(frequency)
}

func (p *ProgressMock) Stop() {
	p.Called()
}

func (p *ProgressMock) Style() *progress.Style {
	arg := p.Called()
	return arg.Get(0).(*progress.Style)
}

func (p *ProgressMock) Render() {
	p.Called()
}
