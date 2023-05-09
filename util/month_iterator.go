package util

import "time"

type MonthIterator interface {
	Begin() time.Time
	End() time.Time
	Next() MonthIterator
	Prev() MonthIterator
}

type monthIteratorImpl struct {
	month int
	year  int
}

func (monthIterator monthIteratorImpl) Begin() time.Time {
	return time.Date(monthIterator.year, time.Month(monthIterator.month), 1, 0, 0, 0, 0, time.Local)
}

func (monthIterator monthIteratorImpl) End() time.Time {
	return monthIterator.Next().Begin().Add(-time.Nanosecond)
}

func (monthIterator monthIteratorImpl) Next() MonthIterator {
	if monthIterator.month == 12 {
		return monthIteratorImpl{month: 1, year: monthIterator.year + 1}
	}
	return monthIteratorImpl{month: monthIterator.month + 1, year: monthIterator.year}
}

func (monthIterator monthIteratorImpl) Prev() MonthIterator {
	if monthIterator.month == 1 {
		return monthIteratorImpl{
			month: 12,
			year:  monthIterator.year - 1,
		}
	}
	return monthIteratorImpl{month: monthIterator.month - 1, year: monthIterator.year}
}

func MonthIterFrom(t time.Time) MonthIterator {
	return monthIteratorImpl{month: int(t.Month()), year: t.Year()}
}
