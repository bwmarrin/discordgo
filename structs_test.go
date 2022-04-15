package discordgo

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestEmbeddingEpochMsTime(t *testing.T) {
	// example struct used in below test cases to show expected use of EpochMsTime in structs
	type someObject struct {
		ID    string       `json:"id"`
		Epoch *EpochMsTime `json:"epoch,omitempty"` // use a pointer here to support omitempty
	}

	// helper to get EpochMsTime pointer from a time.Time
	ep := func(in time.Time) *EpochMsTime {
		e := EpochMsTime(in)
		return &e
	}

	// list of test cases to run
	cases := []struct {
		ID       string
		Epoch    *EpochMsTime
		Expected string
	}{
		{"nil", nil, `{"id":"nil"}`},
		{"valid", ep(time.Date(2021, time.November, 19, 6, 0, 0, 0, time.UTC)), `{"id":"valid","epoch":1637301600000}`},
		{"epoch-zero", ep(time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)), `{"id":"epoch-zero","epoch":0}`},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d-%s", i+1, tc.ID), func(t *testing.T) {
			e1 := someObject{
				ID:    tc.ID,
				Epoch: tc.Epoch,
			}

			b, err := json.Marshal(e1)
			if err != nil {
				t.Fatalf("error in marshal: %v", err)
			}

			actual := string(b)
			if actual != tc.Expected {
				t.Fatalf("marshal mismatch:\n\tExpected: %v\n\tActual: %v", tc.Expected, actual)
			}

			var e2 someObject
			if err := json.Unmarshal(b, &e2); err != nil {
				t.Fatalf("error in unmarshal: %v", err)
			}

			switch {
			case e1.Epoch == nil && e2.Epoch == nil:
				// e1 == e2, so success (fallthrough)
			case e1.Epoch == nil:
				t.Fatalf("started as %v, but remarshalled as nil", e1.Epoch)
			case e2.Epoch == nil:
				t.Fatalf("started as nil, but remarshalled as %v", e2.Epoch)
			case !e1.Epoch.Equal(*e2.Epoch):
				t.Fatalf("epochs not equal after marshal/unmarshal:\n\tBefore: %v\n\tAfter: %v", e1.Epoch.Time(), e2.Epoch.Time())
			}
		})
	}
}
