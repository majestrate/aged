package aged

import (
	"bytes"
	"os"
	"time"
)

// DoBProvider provides Date of Birth for the current user.
type DoBProvider interface {
	// GetDoB gets the date of birth of the current user.
	GetDoB() (time.Time, error)
}

// FlatFileDoBProvider provides Date of Birth from a flat file.
type FlatFileDoBProvider struct {
	Filename string
}

func (p *FlatFileDoBProvider) GetDoB() (time.Time, error) {
	buf, err := os.ReadFile(p.Filename)
	if err != nil {
		return time.Time{}, err
	}
	buf = bytes.TrimSpace(buf)
	return time.Parse(time.DateOnly, string(buf))
}
