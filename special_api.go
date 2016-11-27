package main

import (
	"encoding/json"
	"net/http"
)

func specialAPIHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case http.MethodGet:

		//scheduled   := r.FormValue("scheduled")

		u := authenticate(r)

		if u == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		date := r.FormValue("gameDate")

		games := getGamesEx(u, date)

		if len(games) == 0 {
			w.WriteHeader(http.StatusNotFound)
		} else {
			j, _ := json.Marshal(games)
			w.Write(j)
		}

	case http.MethodPost:
	case http.MethodPut:
	case http.MethodDelete:
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)

	}

} // specialAPIHandler
