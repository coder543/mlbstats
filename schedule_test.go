package main

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

// As I refactored the Schedule type into multiple smaller types, it's good to make sure
// that no data is being unexpectedly lost.
func TestScheduleRoundTrip(t *testing.T) {
	sample, err := os.ReadFile("schedule.json")
	if err != nil {
		panic(err)
	}

	var schedule Schedule
	err = json.Unmarshal(sample, &schedule)
	if err != nil {
		panic(err)
	}

	marshalled, err := json.Marshal(schedule)
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

func GenerateGame(teamID int, doubleHeader string, statusCode GameStatusCode, startTimeTBD bool, date time.Time) Game {
	return Game{
		GameDate: &date,
		Status: GameStatus{
			StartTimeTBD: startTimeTBD,
			StatusCode:   statusCode,
		},
		Teams: struct {
			Away TeamStatus `json:"away"`
			Home TeamStatus `json:"home"`
		}{
			Home: TeamStatus{
				Team: Identifiable{
					ID: teamID,
				},
			},
		},
		DoubleHeader: doubleHeader,
	}
}

func TestSortGamesWithPreferredTeam(t *testing.T) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	type args struct {
		teamID int
		games  []Game
	}
	tests := []struct {
		name string
		args args
		want []Game
	}{
		{
			name: "Simple reorder",
			args: args{
				teamID: 115,
				games: []Game{
					GenerateGame(3, "N", GameStatusFinal, false, today),
					GenerateGame(115, "N", GameStatusFinal, false, today),
				},
			},
			want: []Game{
				GenerateGame(115, "N", GameStatusFinal, false, today),
				GenerateGame(3, "N", GameStatusFinal, false, today),
			},
		},
		{
			name: "No relevant games, no reorder",
			args: args{
				teamID: 115,
				games: []Game{
					GenerateGame(1, "N", GameStatusFinal, false, today),
					GenerateGame(2, "N", GameStatusFinal, false, today),
				},
			},
			want: []Game{
				GenerateGame(1, "N", GameStatusFinal, false, today),
				GenerateGame(2, "N", GameStatusFinal, false, today),
			},
		},
		{
			name: "Traditional double header, TBD needs swap",
			args: args{
				teamID: 115,
				games: []Game{
					GenerateGame(1, "N", GameStatusFinal, false, today),
					GenerateGame(115, "Y", GameStatusFinal, true, today),
					GenerateGame(115, "Y", GameStatusFinal, false, today),
				},
			},
			want: []Game{
				GenerateGame(115, "Y", GameStatusFinal, false, today),
				GenerateGame(115, "Y", GameStatusFinal, true, today),
				GenerateGame(1, "N", GameStatusFinal, false, today),
			},
		},
		{
			name: "Traditional double header, in progress needs swap",
			args: args{
				teamID: 115,
				games: []Game{
					GenerateGame(1, "N", GameStatusFinal, false, today),
					GenerateGame(115, "Y", GameStatusFinal, false, today),
					GenerateGame(115, "Y", GameStatusInProgress, true, today),
				},
			},
			want: []Game{
				GenerateGame(115, "Y", GameStatusInProgress, true, today),
				GenerateGame(115, "Y", GameStatusFinal, false, today),
				GenerateGame(1, "N", GameStatusFinal, false, today),
			},
		},
		{
			name: "Traditional double header, no swap",
			args: args{
				teamID: 115,
				games: []Game{
					GenerateGame(1, "N", GameStatusFinal, false, today),
					GenerateGame(115, "Y", GameStatusFinal, false, today),
					GenerateGame(115, "Y", GameStatusFinal, true, today),
				},
			},
			want: []Game{
				GenerateGame(115, "Y", GameStatusFinal, false, today),
				GenerateGame(115, "Y", GameStatusFinal, true, today),
				GenerateGame(1, "N", GameStatusFinal, false, today),
			},
		},
		{
			name: "Split double header, GameDate needs swap",
			args: args{
				teamID: 115,
				games: []Game{
					GenerateGame(1, "N", GameStatusFinal, false, today),
					GenerateGame(115, "S", GameStatusFinal, false, today.Add(30*time.Minute)),
					GenerateGame(115, "S", GameStatusFinal, false, today),
				},
			},
			want: []Game{
				GenerateGame(115, "S", GameStatusFinal, false, today),
				GenerateGame(115, "S", GameStatusFinal, false, today.Add(30*time.Minute)),
				GenerateGame(1, "N", GameStatusFinal, false, today),
			},
		},
		{
			name: "Split double header, no swap",
			args: args{
				teamID: 115,
				games: []Game{
					GenerateGame(1, "N", GameStatusFinal, false, today),
					GenerateGame(115, "S", GameStatusFinal, false, today),
					GenerateGame(115, "S", GameStatusFinal, false, today.Add(30*time.Minute)),
				},
			},
			want: []Game{
				GenerateGame(115, "S", GameStatusFinal, false, today),
				GenerateGame(115, "S", GameStatusFinal, false, today.Add(30*time.Minute)),
				GenerateGame(1, "N", GameStatusFinal, false, today),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SortGamesWithPreferredTeam(tt.args.teamID, tt.args.games); !cmp.Equal(got, tt.want) {
				t.Errorf("SortGamesWithPreferredTeam() = %s", cmp.Diff(got, tt.want))
			}
		})
	}
}
