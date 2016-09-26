package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func getTeamAdmins(team *Team) []TeamAdmin {
	return nil
} // getTeamAdmins

func getAllTeamAdmins(leagueid string) []TeamAdmin {

	rows, err := config.Database.Query(
		TeamAdminGetAll, leagueid,
	)

	if err != nil {
		log.Println(err)
	}

	defer rows.Close()

	teamAdmins := []TeamAdmin{}

	for rows.Next() {

		ta := TeamAdmin{}

		err := rows.Scan(&ta.ID, &ta.LeagueID, &ta.TeamID, &ta.UserID, &ta.TeamName,
			&ta.UserEmail)

		if err == sql.ErrNoRows || err != nil {
			log.Println(err)
			return nil
		}

		//ta.URL = fmt.Sprintf("/leagues/%s/teams/%s", league.Canonical, t.Canonical)

		teamAdmins = append(teamAdmins, ta)

	}

	return teamAdmins

} // getAllTeamAdmins

func addTeamAdmin(leagueid string, teamid string, userid string) {

	_, err := config.Database.Exec(
		TeamAdminCreate, leagueid, teamid, userid,
	)

	if err != nil {
		log.Println(err)
	}

} // addTeamAdmin

func teamAdminHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:

		// only certain privileges can do this

		vars := mux.Vars(r)

		leagueid := vars["league"]
		teamid := vars["team"]

		email := r.FormValue("email")

		u := getUserByEmail(email)

		if u == nil {

			_, err := config.Database.Exec(
				UserCreateUnregistered, email,
			)

			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
			}

		} else {

			l := getLeague(leagueid)

			t := getTeam(teamid)

			addTeamAdmin(l.ID, t.ID, u.ID)

		}

	case http.MethodPut:
	case http.MethodDelete:
	case http.MethodGet:
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

} // teamAdminHandler
