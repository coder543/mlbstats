package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"time"
)

type Schedule struct {
	Copyright            string `json:"copyright"`
	TotalItems           int    `json:"totalItems"`
	TotalEvents          int    `json:"totalEvents"`
	TotalGames           int    `json:"totalGames"`
	TotalGamesInProgress int    `json:"totalGamesInProgress"`
	Dates                []Date `json:"dates"`
}

type Date struct {
	Date                 string        `json:"date"`
	TotalItems           int           `json:"totalItems"`
	TotalEvents          int           `json:"totalEvents"`
	TotalGames           int           `json:"totalGames"`
	TotalGamesInProgress int           `json:"totalGamesInProgress"`
	Games                []Game        `json:"games"`
	Events               []interface{} `json:"events"`
}

type Game struct {
	GamePk       int        `json:"gamePk"`
	Link         string     `json:"link"`
	GameType     string     `json:"gameType"`
	Season       string     `json:"season"`
	GameDate     *time.Time `json:"gameDate"`
	OfficialDate string     `json:"officialDate"`
	Status       GameStatus `json:"status"`
	Teams        struct {
		Away TeamStatus `json:"away"`
		Home TeamStatus `json:"home"`
	} `json:"teams"`
	Venue   Identifiable `json:"venue"`
	Content struct {
		Link string `json:"link"`
	} `json:"content"`
	IsTie                  bool   `json:"isTie"`
	GameNumber             int    `json:"gameNumber"`
	PublicFacing           bool   `json:"publicFacing"`
	DoubleHeader           string `json:"doubleHeader"`
	GamedayType            string `json:"gamedayType"`
	Tiebreaker             string `json:"tiebreaker"`
	CalendarEventID        string `json:"calendarEventID"`
	SeasonDisplay          string `json:"seasonDisplay"`
	DayNight               string `json:"dayNight"`
	ScheduledInnings       int    `json:"scheduledInnings"`
	ReverseHomeAwayStatus  bool   `json:"reverseHomeAwayStatus"`
	InningBreakLength      int    `json:"inningBreakLength,omitempty"`
	GamesInSeries          int    `json:"gamesInSeries"`
	SeriesGameNumber       int    `json:"seriesGameNumber"`
	SeriesDescription      string `json:"seriesDescription"`
	RecordSource           string `json:"recordSource"`
	IfNecessary            string `json:"ifNecessary"`
	IfNecessaryDescription string `json:"ifNecessaryDescription"`
}

type GameStatus struct {
	AbstractGameState string         `json:"abstractGameState"`
	CodedGameState    string         `json:"codedGameState"`
	DetailedState     string         `json:"detailedState"`
	StatusCode        GameStatusCode `json:"statusCode"`
	StartTimeTBD      bool           `json:"startTimeTBD"`
	AbstractGameCode  string         `json:"abstractGameCode"`
}

type GameStatusCode string

const (
	GameStatusFinal      GameStatusCode = "F"
	GameStatusInProgress GameStatusCode = "P" // TODO: find the real code for an in-progress game
)

type TeamStatus struct {
	LeagueRecord struct {
		Wins   int    `json:"wins"`
		Losses int    `json:"losses"`
		Pct    string `json:"pct"`
	} `json:"leagueRecord"`
	Score        int          `json:"score"`
	Team         Identifiable `json:"team"`
	IsWinner     bool         `json:"isWinner"`
	SplitSquad   bool         `json:"splitSquad"`
	SeriesNumber int          `json:"seriesNumber"`
}

func SortGamesWithPreferredTeam(teamID int, games []Game) []Game {
	preferredGames := make([]Game, 0, len(games))
	otherGames := make([]Game, 0, len(games))

	for _, game := range games {
		if teamID == game.Teams.Home.Team.ID || teamID == game.Teams.Away.Team.ID {
			preferredGames = append(preferredGames, game)
			continue
		}

		otherGames = append(otherGames, game)
	}

	if len(preferredGames) > 1 {
		if len(preferredGames) != 2 {
			log.Println("Unexpected 3+ game day, sorting may be incorrect")
		}

		switch preferredGames[0].DoubleHeader {
		case "Y":
			if preferredGames[0].Status.StartTimeTBD == preferredGames[1].Status.StartTimeTBD {
				log.Println("Encountered traditional doubleheader with unexpected StartTimeTBD values, sorting may be incorrect")
			}

			// sort the StartTimeTBD game into the second slot
			if preferredGames[0].Status.StartTimeTBD {
				preferredGames[0], preferredGames[1] = preferredGames[1], preferredGames[0]
			}

		case "S":
			sort.Slice(preferredGames, func(i, j int) bool {
				return preferredGames[i].GameDate.Before(*preferredGames[j].GameDate)
			})

		default:
			log.Println("Unexpected multi-game day when DoubleHeader is neither Y nor S")
		}

		// If the second game is in progress, swap it to the front
		if preferredGames[1].Status.StatusCode == GameStatusInProgress {
			preferredGames[0], preferredGames[1] = preferredGames[1], preferredGames[0]
		}
	}

	return append(preferredGames, otherGames...)
}

var GetUpstreamSchedule = GetUpstreamScheduleReal

func GetUpstreamScheduleMock(_ context.Context, _ string) (Schedule, http.Header, error) {
	sample, err := os.ReadFile("schedule.json")
	if err != nil {
		return Schedule{}, nil, err
	}

	var schedule Schedule
	err = json.Unmarshal(sample, &schedule)
	if err != nil {
		return Schedule{}, nil, err
	}

	return schedule, http.Header{}, nil
}

func GetUpstreamScheduleReal(ctx context.Context, date string) (Schedule, http.Header, error) {
	client := http.DefaultClient

	req, err := http.NewRequest(
		http.MethodGet,
		"https://statsapi.mlb.com/api/v1/schedule?date="+url.QueryEscape(date)+"&sportId=1&language=en",
		nil,
	)
	if err != nil {
		return Schedule{}, nil, err
	}
	req = req.WithContext(ctx)

	resp, err := client.Do(req)
	if err != nil {
		return Schedule{}, nil, err
	}
	defer resp.Body.Close()

	var schedule Schedule
	err = json.NewDecoder(resp.Body).Decode(&schedule)
	return schedule, resp.Header, err
}
