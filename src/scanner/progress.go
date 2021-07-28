package scanner

import (
	"io"
	"time"

	"github.com/jedib0t/go-pretty/v6/progress"
)

// ProgressWriter abstraction of progress writer
type ProgressWriter interface {
	// Start initiates progress writing progress, if total is unknown should be zero
	Start(total int64)
	// Increment increments progress
	Increment()
	// Stop labels progress as finished and stops updating progress
	Stop()
}

// NewPrettyProgressWriter creates PrettyProgressWriter
func NewPrettyProgressWriter(output io.Writer) *PrettyProgressWriter {
	pw := progress.NewWriter()
	pw.SetAutoStop(false)
	pw.SetTrackerLength(50)
	pw.ShowETA(true)
	pw.ShowOverallTracker(false)
	pw.ShowTime(true)
	pw.ShowTracker(true)
	pw.ShowValue(true)
	pw.SetMessageWidth(13)
	pw.SetNumTrackersExpected(1)
	pw.SetSortBy(progress.SortByPercentDsc)
	pw.SetStyle(progress.StyleDefault)
	pw.SetTrackerPosition(progress.PositionRight)
	pw.SetUpdateFrequency(time.Millisecond * 10)
	pw.Style().Colors = progress.StyleColorsExample
	pw.Style().Options.PercentFormat = "%4.1f%%"
	pw.SetOutputWriter(output)

	return &PrettyProgressWriter{pw: pw}
}

// PrettyProgressWriter progress writer using go-pretty/progress library
type PrettyProgressWriter struct {
	pw      progress.Writer
	tracker *progress.Tracker
}

// Start initiates progress writing progress, if total is unknown should be zero
func (p *PrettyProgressWriter) Start(total int64) {
	p.tracker = &progress.Tracker{Message: "Scanning keys", Total: total, Units: progress.UnitsDefault}
	p.pw.AppendTracker(p.tracker)

	go p.pw.Render()
}

// Increment increments progress
func (p *PrettyProgressWriter) Increment() {
	p.tracker.Increment(1)
}

// Stop labels progress as finished and stops updating progress
func (p *PrettyProgressWriter) Stop() {
	p.tracker.MarkAsDone()
	p.tracker.PercentDone()
	p.pw.Stop()
}
