package discordgo

import (
	"testing"
	"time"
)

func TestTimestampParse(t *testing.T) {
	ts, err := Timestamp("2016-03-24T23:15:59.605000+00:00").Parse()
	if err != nil {
		t.Fatal(err)
	}
	if ts.Year() != 2016 || ts.Month() != time.March || ts.Day() != 24 {
		t.Error("Incorrect date")
	}
	if ts.Hour() != 23 || ts.Minute() != 15 || ts.Second() != 59 {
		t.Error("Incorrect time")
	}

	_, offset := ts.Zone()
	if offset != 0 {
		t.Error("Incorrect timezone")
	}
}
