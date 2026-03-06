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
	dlt := time.Now().Sub(dob)
	year, err := time.ParseDuration("8760h")
	if err != nil {
		return 0, dbus.NewError("InternalError", []any{err.Error()})
	}
	age := dlt.Truncate(year).Hours() / 8760
	return int(age), nil
}
