package main

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// As I refactored the Teams type into multiple smaller types, it's good to make sure
// that no data is being unexpectedly lost.
func TestTeamsRoundTrip(t *testing.T) {
	sample, err := os.ReadFile("teams.json")
	if err != nil {
		panic(err)
	}

	var team Teams
	err = json.Unmarshal(sample, &team)
	if err != nil {
		panic(err)
	}

	marshalled, err := json.Marshal(team)
	if err != nil {
		panic(err)
	}

	var rawSample, rawMarshalled any
	err = json.Unmarshal(sample, &rawSample)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(marshalled, &rawMarshalled)
	if err != nil {
		panic(err)
	}

	if !cmp.Equal(rawSample, rawMarshalled) {
		t.Error("mismatch:", cmp.Diff(sample, marshalled))
	}
}
