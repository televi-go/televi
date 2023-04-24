package profiler

import (
	"fmt"
	"io"
	"sync"
	"time"
)

type Throughput struct {
	series       []Measure
	LabelPrefix  string
	seriesAccess sync.Mutex
}

type parameters struct {
	CalledTimes int
	AvgTime     int64
	MaxTime     int64
	MinTime     int64
}

func (params parameters) write(label string, writer io.Writer) {
	fmt.Fprintf(writer, "Measuring %s:\n", label)
	fmt.Fprintf(writer, "Called %d times\n", params.CalledTimes)
	fmt.Fprintf(writer, "Min execute time: %v\n", time.Duration(params.MinTime))
	fmt.Fprintf(writer, "Max execute time: %v\n", time.Duration(params.MaxTime))
	fmt.Fprintf(writer, "Avg execute time: %v\n", time.Duration(params.AvgTime))
}

func getParameters(source []time.Duration) parameters {
	var (
		min int64 = int64(source[0])
		avg float64
		max int64
	)

	for _, dur := range source {
		if int64(dur) > max {
			max = int64(dur)
		}

		if int64(dur) < min {
			min = int64(dur)
		}

		avg += float64(dur) / float64(len(source))
	}

	return parameters{
		CalledTimes: len(source),
		MaxTime:     max,
		MinTime:     min,
		AvgTime:     int64(avg),
	}
}

func (thr *Throughput) WriteStats(writer io.Writer) {
	thr.seriesAccess.Lock()
	defer thr.seriesAccess.Unlock()
	mapping := make(map[string][]time.Duration)
	for _, gauge := range thr.series {
		mapping[gauge.label] = append(mapping[gauge.label], gauge.duration)
	}

	for label, durations := range mapping {
		params := getParameters(durations)
		params.write(thr.LabelPrefix+label, writer)
	}
}

type Measure struct {
	label    string
	duration time.Duration
}

type Stopwatch struct {
	thr       *Throughput
	startedAt time.Time
	label     string
}

func (sw *Stopwatch) Record() {
	stopTime := time.Now()
	go func() {
		sw.thr.seriesAccess.Lock()
		defer sw.thr.seriesAccess.Unlock()
		sw.thr.series = append(sw.thr.series, Measure{
			label:    sw.label,
			duration: stopTime.Sub(sw.startedAt),
		})
	}()
}

func (thr *Throughput) NewStopwatch(label string) Stopwatch {
	return Stopwatch{
		thr:       thr,
		startedAt: time.Now(),
		label:     label,
	}
}
