package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func doesLeagueExist(leagueID string, leagues []League) bool {

  for _, l := range leagues {

		if leagueID == l.ID {
			return true
		}

	}

	return false

} // doesLeagueExist

func mergeLeagues(src []League, dest []League) []League {

	for _, l := range src {

	  if !doesLeagueExist(l.ID, dest) {
			dest = append(dest, l)
		}

	}

	return dest

} // mergeLeagues

func getDefaultLeague(user *User) *League {

	if user.DefaultLeague == 0 {
		return getLatestLeague(user)
	}

	row := config.Database.QueryRow(
		LeagueGet, user.DefaultLeague,
	)

	l := League{}

	err := row.Scan(&l.ID, &l.Name, &l.Canonical, &l.Icon, &l.Visible,
		&l.Official, &l.Metric, &l.City, &l.Country, &l.Location)

	if err == sql.ErrNoRows {
		log.Println(err)
		return nil
	}

	return &l

} // getDefaultLeague

func getLatestLeague(user *User) *League {

	row := config.Database.QueryRow(
		UserLeagueGetLatest, user.ID,
	)

	l := League{}

	err := row.Scan(&l.ID, &l.Name, &l.Canonical, &l.Icon)

	if err == sql.ErrNoRows {
		log.Println(err)
		return nil
	}
 
	return &l

} // getLatestLeague

func getAllLeagues() []League {

	rows, err := config.Database.Query(
		LeagueGetAll,
	)

	if err != nil {
		log.Println("getAllLeagues: ", err)
		return nil
	}

	defer rows.Close()

	leagues := []League{}

	for rows.Next() {

		l := League{}

		err := rows.Scan(&l.ID, &l.Name, &l.Icon, &l.Location)

		if err == sql.ErrNoRows || err != nil {
			log.Println("getAllLeagues: ", err)
			return nil
		}

		leagues = append(leagues, l)

	}

	return leagues

} // getAllLeagues

func getFollowedLeagues(user *User) []League {

  rows, err := config.Database.Query(
		UserLeagueFollowGetAll, user.ID,
	)

	if err != nil {
		log.Println("getFollowedLeagues: ", err)
		return nil
	}

	defer rows.Close()

	leagues := []League{}

	for rows.Next() {

		l := League{}

		err := rows.Scan(&l.ID, &l.Name, &l.Canonical, &l.Icon)

		if err == sql.ErrNoRows || err != nil {
			log.Println("getFollowedLeagues: ", err)
			return nil
		}

		leagues = append(leagues, l)

	}

	return leagues

} // getFollowedLeagues

func getLeagues(user *User) []League {

	rows, err := config.Database.Query(
		UserLeagueGetAll, user.ID,
	)

	if err != nil {
		log.Println("getLeagues: ", err)
		return nil
	}

	defer rows.Close()

	leagues := []League{}

	for rows.Next() {

		l := League{}

		err := rows.Scan(&l.ID, &l.Name, &l.Canonical, &l.Icon)

		if err == sql.ErrNoRows || err != nil {
			log.Println("getLeagues: ", err)
			return nil
		}

		leagues = append(leagues, l)

	}

	return leagues

} // getLeagues

func getAllMyLeagues(user *User) []League {

	rows, err := config.Database.Query(
		UserAllLeagueGetAll, user.ID, user.ID,
	)

	if err != nil {
		log.Println("getAllMyLeagues: ", err)
		return nil
	}

	defer rows.Close()

	leagues := []League{}

	for rows.Next() {

		l := League{}

		err := rows.Scan(&l.ID, &l.Name, &l.Canonical, &l.Icon)

		if err == sql.ErrNoRows || err != nil {
			log.Println("getAllMyLeagues: ", err)
			return nil
		}

		leagues = append(leagues, l)

	}

	return leagues

} // getAllMyLeagues

func getLeague(league string) *League {

	row := config.Database.QueryRow(
		LeagueGet, league,
	)

	l := League{}

	err := row.Scan(&l.ID, &l.Name, &l.Canonical, &l.Icon, &l.Visible,
		&l.Official, &l.Metric, &l.City, &l.Country, &l.Location)

	if err == sql.ErrNoRows {
		log.Println("getLeague: ", err)
		return nil
	}

	l.URL = fmt.Sprintf("/leagues/%s", l.Canonical)

	return &l

} // getLeague

func getLeagueByCanonical(canonical string) *League {

	row := config.Database.QueryRow(
		LeagueGetByCanonical, canonical,
	)

	l := League{}

	err := row.Scan(&l.ID, &l.Name, &l.Canonical, &l.Icon, &l.Visible,
		&l.Official, &l.Metric, &l.City, &l.Country, &l.Location)

	if err == sql.ErrNoRows {
		log.Println(err)
		return nil
	}

	l.URL = fmt.Sprintf("/leagues/%s", l.Canonical)

	return &l

} // getLeagueByCanonical

func leagueAPIHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	league := vars["league"]

	u := authenticate(r)

	switch r.Method {
	case http.MethodPost:

		name := r.FormValue("name")
		sport := r.FormValue("sport")
		periods := r.FormValue("periods")
		duration := r.FormValue("duration")

		// TODO: transaction
		res, err := config.Database.Exec(
			LeagueCreate, name, generateCanonical(name), sport,
		)

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusConflict)
			w.Write(createJSON("League name already taken."))
			return
		}

		// this code is suspect, if multiple requests come at
		// the same time, you could
		// get league_id of another league potentially
		league_id, err2 := res.LastInsertId()

		if err2 != nil {
			log.Println(err2)
			w.WriteHeader(http.StatusConflict)
		}

		_, err3 := config.Database.Exec(
			LeagueAdminCreate, league_id, u.ID,
		)

		if err3 != nil {
			log.Println(err3)
			w.WriteHeader(http.StatusConflict)
		}

		_, err4 := config.Database.Exec(
			SeasonCreate, league_id, periods, duration,
		)

		if err4 != nil {
			log.Println(err4)
			w.WriteHeader(http.StatusConflict)
		}

		l := League{
			ID: fmt.Sprintf("%d", league_id),
		}

		j, _ := json.Marshal(l)

		w.Write(j)

	case http.MethodGet:

		if league == "" {

			leagues := getLeagues(u)
						
			if leagues == nil {
				w.WriteHeader(http.StatusNotFound)
				w.Write(createJSON("No leagues found for user."))
			} else {
				j, _ := json.Marshal(leagues)
				w.Write(j)
			}

		} else {

			if league == "all" {

				leagues := getAllLeagues()

				if leagues == nil {
					w.WriteHeader(http.StatusNotFound)
					w.Write(createJSON("No leagues found."))
				} else {
					j, _ := json.Marshal(leagues)
					w.Write(j)
				}

			} else if league == "following" {

				leagues := getFollowedLeagues(u)

				if leagues == nil {
					w.WriteHeader(http.StatusNotFound)
					w.Write(createJSON("No leagues found."))
				} else {
					j, _ := json.Marshal(leagues)
					w.Write(j)
				}

			} else {

				l := getLeague(league)

				if len(league) == 0 {
					w.WriteHeader(http.StatusNotFound)
				} else {
					j, _ := json.Marshal(l)
					w.Write(j)
				}

			}

		}

	case http.MethodPut:

		// authorization

		seasonid := r.FormValue("seasonid")
		duration := r.FormValue("duration")
		periods := r.FormValue("periods")

		if periods != "2" || periods != "4" {
			w.WriteHeader(http.StatusBadRequest)
		}

		res, err := config.Database.Exec(
			SeasonUpdate, periods, duration, seasonid,
		)

		log.Println(res)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}

	case http.MethodDelete:
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

} // leagueAPIHandler
