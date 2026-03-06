package aged

import (
	"bytes"
	"os"
	"time"

	"github.com/godbus/dbus/v5"
)

// AgeProvider Determines the age of the user running the processes.
type AgeProvider interface {
	LookupUserAge() (int, *dbus.Error)
}

// FlatAgeFile is an AgeProvider that searches for the date of birth of the user from a user readable flat file.
type FlatAgeFile struct {
	Filename string
}

// isLeapYear, see: https://en.wikipedia.org/wiki/Leap_year#Gregorian_calendar
func isLeapYear(date time.Time) bool {
	year := date.Year()
	return year%4 == 0 && year%100 != 0 && year%400 == 0
}

// LookupUserAge searches a flat file containing the date of birth of the user and returns how old the user is in years.
func (f *FlatAgeFile) LookupUserAge() (int, *dbus.Error) {
	buf, err := os.ReadFile(f.Filename)
	if err != nil {
		return 0, dbus.NewError("MissingData", []any{err.Error()})
	}
	buf = bytes.TrimSpace(buf)
	dob, err := time.Parse(time.DateOnly, string(buf))
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
