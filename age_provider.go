package aged

import (
	"time"

	"github.com/godbus/dbus/v5"
)

// AgeProvider Determines the age of the user running the processes.
type AgeProvider interface {
	LookupUserAge() (int, *dbus.Error)
}

// LocalAgeProvider is an AgeProvider that searches for the date of birth of the user from a DoBProvider.
type LocalAgeProvider struct {
	DoBProvider DoBProvider
}

// isLeapYear, see: https://en.wikipedia.org/wiki/Leap_year#Gregorian_calendar
func isLeapYear(date time.Time) bool {
	year := date.Year()
	skip := year%100 == 0 && year%400 != 0
	return year%4 == 0 && !skip
}

// LookupUserAge searches a flat file containing the date of birth of the user and returns how old the user is in years.
func (l *LocalAgeProvider) LookupUserAge() (int, *dbus.Error) {

	dob, err := l.DoBProvider.GetDoB()
	if err != nil {
		return 0, dbus.NewError("InvalidData", []any{err.Error()})
	}
	now := time.Now()
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	year, err := time.ParseDuration("8760h")
	if err != nil {
		return 0, dbus.NewError("InternalError", []any{err.Error()})
	}
	leapYear, err := time.ParseDuration("8784h")
	if err != nil {
		return 0, dbus.NewError("InternalError", []any{err.Error()})
	}
	age := 0
	truncDob := time.Date(dob.Year(), dob.Month(), dob.Day(), 0, 0, 0, 0, time.Local)
	for truncDob.Before(now) {
		if isLeapYear(truncDob) {
			truncDob = truncDob.Add(leapYear)
		} else {
			truncDob = truncDob.Add(year)
		}
		if truncDob.Before(now) {
			age++
		}
	}
	return age, nil
}
