package grouping

import (
	"github.com/televi-go/televi/util"
	"time"
)

type CanBeGroupedByDate interface {
	GetTime() time.Time
	InflateFormatted(formatted string)
}

type DateGroup[TData CanBeGroupedByDate] struct {
	Title                       string
	Count                       int
	MonthFactor                 float64
	PercentsComparedToReference float64
	Data                        []TData
}

type SplittingAlgorithm[T CanBeGroupedByDate] func(source []T, mi util.MonthIterator) (conforming, notConforming []T)

type AlgorithmChainNode[T CanBeGroupedByDate] struct {
	Algorithm   SplittingAlgorithm[T]
	Title       string
	TimeFormat  string
	MonthFactor float64
	IsRecurring bool
	Next        *AlgorithmChainNode[T]
}

type AlgorithmChainRoot[T CanBeGroupedByDate] struct {
	RootNode             *AlgorithmChainNode[T]
	ReferenceGroupFilter func(T) bool
	ReferenceVal         int
	refValEnsured        bool
}

func (root *AlgorithmChainRoot[T]) EnsureRefVal(source []T) int {
	if root.refValEnsured {
		return root.ReferenceVal
	}
	root.refValEnsured = true
	if root.ReferenceGroupFilter == nil {
		return root.ReferenceVal
	}
	for _, t := range source {
		if root.ReferenceGroupFilter(t) {
			root.ReferenceVal++
		}
	}
	return root.ReferenceVal
}

type AlgorithmChainBuilder[T CanBeGroupedByDate] struct {
	root    *AlgorithmChainNode[T]
	current *AlgorithmChainNode[T]
}

func (algo *AlgorithmChainBuilder[T]) Link(node AlgorithmChainNode[T]) *AlgorithmChainBuilder[T] {
	if algo.current == nil {
		algo.current = &node
		algo.root = &node
		return algo
	}

	algo.current.Next = &node
	algo.current = algo.current.Next
	return algo
}

func (algo *AlgorithmChainBuilder[T]) Build(refVal int) AlgorithmChainRoot[T] {
	return AlgorithmChainRoot[T]{
		RootNode:     algo.root,
		ReferenceVal: refVal,
	}
}

func NewChainBuilder[T CanBeGroupedByDate]() *AlgorithmChainBuilder[T] {
	return &AlgorithmChainBuilder[T]{}
}

func StandardChainNode[T CanBeGroupedByDate](refVal int) AlgorithmChainRoot[T] {
	return NewChainBuilder[T]().
		Link(AlgorithmChainNode[T]{
			Algorithm:   splitToday[T],
			Title:       "Today",
			TimeFormat:  "15:04",
			MonthFactor: 30,
		}).
		Link(AlgorithmChainNode[T]{
			Algorithm:   splitYesterday[T],
			Title:       "Yesterday",
			TimeFormat:  "15:04",
			MonthFactor: 30,
		}).
		Link(AlgorithmChainNode[T]{
			Algorithm:   splitThisWeek[T],
			Title:       "This week",
			TimeFormat:  "Monday 15:04",
			MonthFactor: 4,
		}).
		Link(AlgorithmChainNode[T]{
			Algorithm:   splitMonth[T],
			Title:       "",
			TimeFormat:  "02.01",
			MonthFactor: 1,
			IsRecurring: true,
			Next:        nil,
		}).
		Build(refVal)
}

func splitToday[T CanBeGroupedByDate](source []T, _ util.MonthIterator) (today []T, notToday []T) {
	return util.FilterOut[T](source, func(elem T) bool {
		return IsToday(elem.GetTime())
	})
}

func splitYesterday[T CanBeGroupedByDate](source []T, _ util.MonthIterator) (yesterday []T, earlier []T) {
	return util.FilterOut(source, func(elem T) bool {
		return IsYesterday(elem.GetTime())
	})
}

func splitThisWeek[T CanBeGroupedByDate](source []T, _ util.MonthIterator) (thisWeek []T, earlier []T) {
	return util.FilterOut(source, func(elem T) bool {
		return IsThisWeek(elem.GetTime())
	})
}

func splitMonth[T CanBeGroupedByDate](source []T, mi util.MonthIterator) (thisMonth []T, earlier []T) {
	return util.FilterOut(source, func(elem T) bool {
		return elem.GetTime().After(mi.Begin()) && elem.GetTime().Before(mi.End())
	})
}

func GroupByDate[T CanBeGroupedByDate](source []T) []DateGroup[T] {
	root := StandardChainNode[T](0)
	prevMi := util.MonthIterFrom(time.Now()).Prev()
	root.ReferenceGroupFilter = func(t T) bool {
		return prevMi.Begin().After(t.GetTime()) && prevMi.End().Before(t.GetTime())
	}
	return root.GroupByDate(source)
}

func (root *AlgorithmChainRoot[T]) GroupByDate(source []T) []DateGroup[T] {
	result := make([]DateGroup[T], 0, 100)
	refVal := float64(root.EnsureRefVal(source))

	remains := source
	var current []T
	algo := root.RootNode
	mi := util.MonthIterFrom(time.Now())
	for len(remains) > 0 && algo != nil {
		current, remains = algo.Algorithm(remains, mi)

		if len(current) > 0 {

			for _, t := range current {
				t.InflateFormatted(t.GetTime().Format(algo.TimeFormat))
			}

			result = append(result, DateGroup[T]{
				Title:                       algo.Title,
				Count:                       len(current),
				MonthFactor:                 algo.MonthFactor,
				PercentsComparedToReference: algo.MonthFactor * float64(len(current)) / refVal,
				Data:                        current,
			})
		}

		if !algo.IsRecurring {
			algo = algo.Next
		} else {
			mi = mi.Prev()
		}
	}

	return result
}
