package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func leagueTeamExists(leagueid string, name string) bool {

	rows, err := config.Database.Query(
		TeamGetAllNames, leagueid,
	)

	if err != nil {
		log.Println(err)
		return false
	}

	defer rows.Close()

	for rows.Next() {

		t := Team{}

		err := rows.Scan(&t.Name, &t.Canonical)

		if err == sql.ErrNoRows || err != nil {
			log.Println(err)
			return false
		}

		if name == t.Name || t.Canonical == generateCanonical(name) {
			return true
		}

	}

	return false

} // leagueTeamExists

func getTeams(league *League) []Team {

	rows, err := config.Database.Query(
		TeamGetAll, league.ID,
	)

	if err != nil {
		log.Println("getTeams: ", err)
	}

	defer rows.Close()

	teams := []Team{}

	for rows.Next() {

		t := Team{}

		err := rows.Scan(&t.ID, &t.Name, &t.Canonical, &t.Icon, &t.LeagueID)

		if err == sql.ErrNoRows || err != nil {
			log.Println("getTeams: ", err)
			return nil
		}

		t.URL = fmt.Sprintf("/leagues/%s/teams/%s", league.Canonical, t.Canonical)

		teams = append(teams, t)

	}

	return teams

} // getTeams

func getTeamsAsMap(league *League) map[string]Team {

	rows, err := config.Database.Query(
		TeamGetAll, league.ID,
	)

	if err != nil {
		log.Println("getTeamsAsMap", err)
	}

	defer rows.Close()

	teams := map[string]Team{}

	for rows.Next() {

		t := Team{}

		err := rows.Scan(&t.ID, &t.Name, &t.Canonical, &t.Icon, &t.LeagueID)

		if err == sql.ErrNoRows || err != nil {
			log.Println("getTeamsAsMap: ", err)
			return nil
		}

		t.URL = fmt.Sprintf("/leagues/%s/teams/%s", league.Canonical, t.Canonical)

		teams[t.ID] = t

	}

	return teams

} // getTeamsAsMap

func getTeamsEx(league *League, teamID string) []Team {

	rows, err := config.Database.Query(
		TeamGetAllBut, league.ID, teamID,
	)

	if err != nil {
		log.Println("getTeamsEx: ", err)
	}

	defer rows.Close()

	teams := []Team{}

	for rows.Next() {

		t := Team{}

		err := rows.Scan(&t.ID, &t.Name, &t.Canonical, &t.Icon, &t.LeagueID)

		if err == sql.ErrNoRows || err != nil {
			log.Println("getTeamsEx: ", err)
			return nil
		}

		t.URL = fmt.Sprintf("/leagues/%s/teams/%s", league.Canonical, t.Canonical)

		teams = append(teams, t)

	}

	return teams

} // getTeamsEx

func getTeam(id string) *Team {

	row := config.Database.QueryRow(
		TeamGet, id,
	)

	t := Team{}

	err := row.Scan(&t.ID, &t.Name, &t.Canonical, &t.Icon, &t.LeagueID,
		&t.LeagueCanonical, &t.LeagueName)

	if err != nil {
		log.Println("getTeam: ", err)
		return nil
	}

	t.URL = fmt.Sprintf("/leagues/%s/teams/%s", t.LeagueCanonical,
		t.Canonical)

	return &t

} // getTeam

func getTeamByCanonical(league *League, canonical string) *Team {

	row := config.Database.QueryRow(
		TeamGetByCanonical, canonical,
	)

	t := Team{}

	err := row.Scan(&t.ID, &t.Name, &t.Canonical, &t.Icon, &t.LeagueID)

	if err != nil {
		log.Println(err)
		return nil
	}

	t.URL = fmt.Sprintf("/leagues/%s/teams/%s", league.Canonical, canonical)

	return &t

} // getTeamByCanonical

func teamAPIHandler(w http.ResponseWriter, r *http.Request) {

	u := authenticate(r)

	vars := mux.Vars(r)

	leagueid := vars["league"]
	teamid := vars["team"]

	switch r.Method {
	case http.MethodPost:

		// check for nil

		name := strings.Title(r.FormValue("name"))
		icon := r.FormValue("icon")

		if leagueTeamExists(leagueid, name) {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("Team name already exists in league."))
			return
		}

		_, err := config.Database.Exec(
			TeamCreate, name, generateCanonical(name), icon, leagueid,
		)

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusConflict)
		}

	case http.MethodGet:

		// TODO: auth
		//league := r.FormValue("league")
		//team   := r.FormValue("team")

		l, err := authorizeLeague(u, leagueid)

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if teamid != "" {

			t := getTeam(teamid)

			if t == nil {
				w.WriteHeader(http.StatusNotFound)
			} else {
				j, _ := json.Marshal(t)
				w.Write(j)
			}

		} else {

			teams := getTeams(l)

			if len(teams) == 0 {
				w.WriteHeader(http.StatusNotFound)
			} else {
				j, _ := json.Marshal(teams)
				w.Write(j)
			}

		}

	case http.MethodPut:

		// TODO: auth

		vars := mux.Vars(r)

		//league  := vars["league"]
		team := vars["team"]

		name := r.FormValue("name")
		icon := r.FormValue("icon")

		_, err := config.Database.Exec(
			TeamUpdate, name, icon, team,
		)

		if err != nil {
			log.Println(err)
		}

	case http.MethodDelete:
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)

	}

} // teamAPIHandler
