package util

import "time"

type StopWatch struct {
	start time.Time
}

func NewStopWatch() StopWatch {
	return StopWatch{
		time.Now(),
	}
}

func (stopWatch *StopWatch) Elapsed() time.Duration {
	return time.Since(stopWatch.start)
}

func (stopWatch *StopWatch) Reset() {
	stopWatch.start = time.Now()
}
