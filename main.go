package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	addr := flag.String("addr", ":8080", "The address used for the listening socket")
	useMock := flag.Bool("mock", false, "Controls whether this instance will use mocked data or real data")
	flag.Parse()

	if *useMock {
		log.Println("Using mock data")
		GetUpstreamSchedule = GetUpstreamScheduleMock
	}

	http.HandleFunc("/api/v1/schedule", func(w http.ResponseWriter, r *http.Request) {
		// take in two parameters: date and teamId
		date := r.URL.Query().Get("date")
		teamIDStr := r.URL.Query().Get("teamId")

		if date == "" || teamIDStr == "" {
			WriteErrResponse(w, http.StatusBadRequest, "must provide date and teamId parameters")
			return
		}

		teamID, err := strconv.ParseInt(teamIDStr, 10, 32)
		if err != nil {
			WriteErrResponse(w, http.StatusBadRequest, fmt.Sprintf("improperly formatted teamId: %v", err))
			return
		}

		schedule, header, err := GetUpstreamSchedule(r.Context(), date)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Unable to fetch upstream schedule: %v", err)
			return
		}

		for d, date := range schedule.Dates {
			schedule.Dates[d].Games = SortGamesWithPreferredTeam(int(teamID), date.Games)
		}

		marshalled, err := json.Marshal(schedule)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Unable to marshal response: %v", err)
			return
		}

		// TODO: more efficient way?
		writeHeader := w.Header()
		for k, v := range header {
			writeHeader[k] = v
		}

		_, err = w.Write(marshalled)
		if err != nil {
			log.Printf("Unable to send response: %v", err)
			return
		}
	})

	server := &http.Server{
		Addr:              *addr,
		Handler:           http.DefaultServeMux,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       300 * time.Second,
	}

	log.Println("Listening on", *addr)

	err := server.ListenAndServe()
	if err != nil {
		log.Panicln(err)
	}
}
