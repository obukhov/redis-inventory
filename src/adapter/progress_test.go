package adapter

import (
	"github.com/jedib0t/go-pretty/v6/progress"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type ProgressWriterTestSuite struct {
	suite.Suite
}

func (suite ProgressWriterTestSuite) TestStart() {
	pwMock := &ProgressMock{}
	prettyProgressWriter := &PrettyProgressWriter{
		pw:      pwMock,
		tracker: nil,
	}

	pwMock.On("AppendTracker", mock.Anything).Once().Run(func(args mock.Arguments) {
		tracker, ok := args.Get(0).(*progress.Tracker)

		suite.Assert().True(ok)
		suite.Assert().Equal(int64(42), tracker.Total)
	})

	pwMock.On("Render").Once()

	prettyProgressWriter.Start(42)

	time.Sleep(time.Millisecond)
	pwMock.AssertExpectations(suite.T())
}

func (suite ProgressWriterTestSuite) TestStop() {
	pwMock := &ProgressMock{}
	trackerMock := &TrackerMock{}
	prettyProgressWriter := &PrettyProgressWriter{
		pw:      pwMock,
		tracker: trackerMock,
	}

	pwMock.On("Stop").Once()
	trackerMock.On("MarkAsDone").Once()

	prettyProgressWriter.Stop()

	pwMock.AssertExpectations(suite.T())
	trackerMock.AssertExpectations(suite.T())
}

func (suite ProgressWriterTestSuite) TestIncrement() {
	pwMock := &ProgressMock{}
	trackerMock := &TrackerMock{}

	prettyProgressWriter := &PrettyProgressWriter{
		pw:      pwMock,
		tracker: trackerMock,
	}

	trackerMock.On("Increment", int64(1)).Once()

	prettyProgressWriter.Increment()

	pwMock.AssertExpectations(suite.T())
	trackerMock.AssertExpectations(suite.T())
}

func TestProgressWriterTestSuite(t *testing.T) {
	suite.Run(t, new(ProgressWriterTestSuite))
}
