package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func getLatestSeasonByID(league string) *Season {

	row := config.Database.QueryRow(
		SeasonGetLatest, league,
	)

	s := Season{}

	err := row.Scan(&s.ID, &s.Periods, &s.Duration, &s.LeagueID)

	if err != nil {
		log.Println("getLatestSeasonByID: ", err)
		return nil
	}

	return &s

} // getLatestSeasonByID

func seasonAPIHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:

		vars := mux.Vars(r)
		league := vars["league"]

		periods := r.FormValue("periods")
		duration := r.FormValue("duration")
		start := r.FormValue("start")
		finish := r.FormValue("finish")

		_, err := config.Database.Exec(
			SeasonCreate, league, periods, duration, start, finish,
		)

		if err != nil {

			log.Println(err)
			w.WriteHeader(http.StatusConflict)

		}

	case http.MethodGet:

		// TODO: auth
		vars := mux.Vars(r)

		league := vars["league"]
		//season := vars["season"]

		row := config.Database.QueryRow(
			SeasonGetAll, league,
		)

		s := Season{}

		err := row.Scan(&s.ID, &s.Periods, &s.Duration, &s.LeagueID)

		if err != nil {
			log.Println(err)
			return
		}

		j, _ := json.Marshal(s)

		w.Write(j)

	case http.MethodPut:

		// TODO: auth
		vars := mux.Vars(r)
		seasonid := vars["season"]

		duration := r.FormValue("duration")
		periods := r.FormValue("periods")

		if !(periods == "2" || periods == "4") {
			w.WriteHeader(http.StatusBadRequest)
		}

		_, err := config.Database.Exec(
			SeasonUpdate, periods, duration, seasonid,
		)

		if err != nil {
			log.Println("put seasonsAPIHandler: ", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

	case http.MethodDelete:
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)

	}

} // seasonAPIHandler
