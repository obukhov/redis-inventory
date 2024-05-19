package adapter

import (
	"fmt"
	"io"
	"time"

	"github.com/jedib0t/go-pretty/v6/progress"
)

// ProgressWriter abstraction of progress writer
type ProgressWriter interface {
	// AddTracker adds a tracker to the progress writer
	// id is used to identify the tracker
	AddTracker(id string, total int64)
	// SetNumTrackersExpected sets the number of trackers expected
	SetNumTrackersExpected(numTrackers int)
	// Start initiates progress writing progress, if total is unknown should be zero
	Start()
	// Increment increments progress
	Increment(id string, incrementBy int64)
	// Stop labels progress as finished and stops updating progress
	Stop()
}

// NewPrettyProgressWriter creates PrettyProgressWriter
func NewPrettyProgressWriter(output io.Writer) *PrettyProgressWriter {
	p := &PrettyProgressWriter{pw: progress.NewWriter()}
	p.init(output)

	return p
}

// PrettyProgressWriter progress writer using go-pretty/progress library
type PrettyProgressWriter struct {
	pw       progress.Writer
	trackers map[string]Tracker
}

// Tracker is abstraction over libraries "Tracker" struct
type Tracker interface {
	Increment(value int64)
	MarkAsDone()
}

func (p *PrettyProgressWriter) init(output io.Writer) {
	p.pw.SetAutoStop(false)
	p.pw.SetTrackerLength(50)
	p.pw.ShowETA(true)
	p.pw.ShowOverallTracker(false)
	p.pw.ShowTime(true)
	p.pw.ShowTracker(true)
	p.pw.ShowValue(true)
	p.pw.SetMessageWidth(40)
	p.pw.SetNumTrackersExpected(1)
	p.pw.SetSortBy(progress.SortByPercentDsc)
	p.pw.SetStyle(progress.StyleDefault)
	p.pw.SetTrackerPosition(progress.PositionRight)
	p.pw.SetUpdateFrequency(time.Millisecond * 10)
	p.pw.Style().Colors = progress.StyleColorsExample
	p.pw.Style().Options.PercentFormat = "%4.1f%%"
	p.pw.SetOutputWriter(output)
	p.trackers = make(map[string]Tracker)
}

func (p *PrettyProgressWriter) SetNumTrackersExpected(numTrackers int) {
	if numTrackers > 1 {
		p.pw.ShowOverallTracker(true)
	}
	p.pw.SetNumTrackersExpected(numTrackers)
}

func (p *PrettyProgressWriter) AddTracker(id string, total int64) {
	scanningKeysTracker := &progress.Tracker{Message: fmt.Sprintf("[Shard %s] Scanning keys", id), Total: total, Units: progress.UnitsDefault}
	p.pw.AppendTracker(scanningKeysTracker)
	p.trackers[id] = scanningKeysTracker
}

// Start initiates progress writing progress, if total is unknown should be zero
func (p *PrettyProgressWriter) Start() {
	go p.pw.Render()
}

// Increment increments progress
func (p *PrettyProgressWriter) Increment(id string, incrementBy int64) {
	p.trackers[id].Increment(incrementBy)
}

// Stop labels progress as finished and stops updating progress
func (p *PrettyProgressWriter) Stop() {
	for _, tracker := range p.trackers {
		tracker.MarkAsDone()
	}
	p.pw.Stop()
}
