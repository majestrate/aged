package aged_test

import (
	"github.com/godbus/dbus/v5"
	"github.com/majestrate/aged"
	"testing"
)

type mockAgeProvider int

func (m mockAgeProvider) LookupUserAge() (int, *dbus.Error) {
	return int(m), nil
}

func TestGetAgeBracket(t *testing.T) {
	daemon, err := aged.NewDaemon()
	if err != nil {
		t.Fatal(err.Error())
	}
	assertBracket := func(expected string, minage, maxage int) {
		for n := minage; n <= maxage; n++ {
			daemon.AgeProvider = mockAgeProvider(n)
			bracket, err := daemon.GetAgeBracket()
			if err != nil {
				t.Fatal(err.Error())
			}
			if bracket != expected {
				t.Error("Expected ", n, " to be ", expected, "but was: ", bracket)
				t.FailNow()
			}
		}
	}
	assertBracket(aged.Under13, 0, 12)
	assertBracket(aged.Between13And16, 13, 16)
	assertBracket(aged.Seventeen, 17, 17)
	assertBracket(aged.Adult, 18, 20)
	assertBracket(aged.USAAdult, 21, 100)

}
