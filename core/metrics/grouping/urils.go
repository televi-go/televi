package grouping

import "time"

func IsSameDay(val time.Time, point time.Time) bool {
	yearP, monthP, dayP := point.Date()
	yearV, monthV, dayV := val.Date()
	return yearP == yearV && monthP == monthV && dayP == dayV
}

func IsToday(val time.Time) bool {
	return IsSameDay(val, time.Now())
}

func IsYesterday(val time.Time) bool {
	return IsSameDay(val, time.Now().Add(-time.Hour*24)) && !IsToday(val)
}

func IsThisWeek(val time.Time) bool {
	currYear, currWeek := time.Now().ISOWeek()
	valYear, valWeek := val.ISOWeek()
	return currYear == valYear && currWeek == valWeek
}

func IsThisMonth(val time.Time) bool {
	yearP, monthP, _ := time.Now().Date()
	yearV, monthV, _ := val.Date()
	return yearP == yearV && monthP == monthV
}
