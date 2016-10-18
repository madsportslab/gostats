package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func getTeamPlayers(leagueid string, teamid string) []Player {

	rows, err := config.Database.Query(
		TeamPlayersGetAll, leagueid, teamid,
	)

	if err != nil {
		log.Println(err)
	}

	defer rows.Close()

	players := []Player{}

	for rows.Next() {

		p := Player{}

		err := rows.Scan(&p.ID, &p.First, &p.Middle, &p.Last, &p.Canonical,
			&p.PlayerNumber, &p.PositionID, &p.TeamID, &p.LeagueID)

		if err == sql.ErrNoRows || err != nil {
			log.Println(err)
			return nil
		}

		//p.URL = fmt.Sprintf("/leagues/%s/teams/%s/players/%s", league.Canonical,
		//team.Canonical, p.Canonical)

		if p.First.String != "" {
			p.ShortName = fmt.Sprintf("%c. %s", p.First.String[0], p.Last)
		} else {
			p.ShortName = fmt.Sprintf("%s", p.Last)
		}

		p.FullName = fmt.Sprintf("%s %s", p.First.String, p.Last)

		players = append(players, p)

	}

	return players

} // getTeamPlayers

func getPlayers(league *League) []Player {

	rows, err := config.Database.Query(
		PlayerGetAll, league.ID,
	)

	if err != nil {
		log.Println(err)
	}

	defer rows.Close()

	players := []Player{}

	for rows.Next() {

		p := Player{}

		err := rows.Scan(&p.ID, &p.First, &p.Middle, &p.Last, &p.Canonical,
			&p.PlayerNumber, &p.PositionID, &p.TeamID, &p.LeagueID, &p.TeamName)

		if err == sql.ErrNoRows || err != nil {
			log.Println(err)
			return nil
		}

		p.URL = fmt.Sprintf("/leagues/%s/teams/%s/players/%s", league.Canonical, p.Canonical)

		if p.First.String != "" {
			p.ShortName = fmt.Sprintf("%c. %s", p.First.String[0], p.Last)
		} else {
			p.ShortName = fmt.Sprintf("%s", p.Last)
		}

		p.FullName = fmt.Sprintf("%s %s", p.First.String, p.Last)

		players = append(players, p)

	}

	return players

} // getPlayers

func getPlayer(id string) *Player {

	row := config.Database.QueryRow(
		PlayerGet, id,
	)

	p := Player{}

	err := row.Scan(&p.ID, &p.First, &p.Middle, &p.Last, &p.Canonical,
		&p.PlayerNumber, &p.PositionID, &p.TeamID, &p.LeagueID)

	if err != nil {
		log.Println(err)
	}

	// TODO: get canonical league name
	p.URL = fmt.Sprintf("/leagues/%s/teams/%s/players/%s", p.LeagueID, p.Canonical)

	if p.First.String != "" {
		p.ShortName = fmt.Sprintf("%c. %s", p.First.String[0], p.Last)
	} else {
		p.ShortName = fmt.Sprintf("%s", p.Last)
	}

	p.FullName = fmt.Sprintf("%s %s", p.First.String, p.Last)

	return &p

} // getPlayer

func getPlayerByCanonical(league *League, team *Team, player string) *Player {

	row := config.Database.QueryRow(
		TeamPlayerGet, team.ID, player,
	)

	p := Player{}

	err := row.Scan(&p.ID, &p.First, &p.Middle, &p.Last, &p.Canonical,
		&p.PlayerNumber, &p.PositionID, &p.TeamID, &p.LeagueID)

	if err != nil {
		log.Println(err)
	}

	p.URL = fmt.Sprintf("/leagues/%s/teams/%s/player/%s",
		league.Canonical, team.Canonical, p.Canonical)

	if p.First.String != "" {
		p.ShortName = fmt.Sprintf("%c. %s", p.First.String[0], p.Last)
	} else {
		p.ShortName = fmt.Sprintf("%s", p.Last)
	}

	p.FullName = fmt.Sprintf("%s %s", p.First.String, p.Last)

	return &p

} // getPlayerByCanonical

func playerNumberExists(league string, team string, number string) bool {

  players := getTeamPlayers(league, team)

	for _, p := range players {

		if p.PlayerNumber.String == number {
			return true
		}

	}

	return false

} // playerNumberExists

func playerAPIHandler(w http.ResponseWriter, r *http.Request) {

	u := authenticate(r)

	if u == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodPost:

		vars := mux.Vars(r)

		league := vars["league"]
		team := vars["team"]

		name := r.FormValue("name")

		last, first, middle := parseName(name)

		//height    := r.FormValue("height")
		//weight    := r.FormValue("weight")
		//hand      := r.FormValue("hand")
		playerNumber := r.FormValue("playerNumber")
		//birth     := r.FormValue("birth")
		position := r.FormValue("position")

		if playerNumberExists(league, team, playerNumber) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		if team != "" {

			_, err := config.Database.Exec(
				TeamPlayerCreate, first, middle, last, generateCanonical(name),
				playerNumber, position, league, team,
			)

			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusConflict)
			}

		} else {

			_, err := config.Database.Exec(
				PlayerCreate, first, middle, last, generateCanonical(name),
				playerNumber, position, league,
			)

			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusConflict)
			}

		}

	case http.MethodGet:

		// TODO: auth
		vars := mux.Vars(r)

		league := vars["league"]
		team := vars["team"]
		player := vars["player"]

		if player != "" {

			p := getPlayer(player)

			if p == nil {
				w.WriteHeader(http.StatusNotFound)
			} else {
				j, _ := json.Marshal(p)
				w.Write(j)
			}

		} else {

			if team != "" {

				l := getLeague(league)

				t := getTeam(team)

				players := getTeamPlayers(l.ID, t.ID)

				if len(players) == 0 {
					w.WriteHeader(http.StatusNotFound)
				} else {
					j, _ := json.Marshal(players)
					w.Write(j)
				}

			} else {

				l := getLeague(league)

				players := getPlayers(l)

				if len(players) == 0 {
					w.WriteHeader(http.StatusNotFound)
				} else {
					j, _ := json.Marshal(players)
					w.Write(j)
				}

			}

		}

	case http.MethodPut:

		// TODO: auth

		vars := mux.Vars(r)

		league := vars["league"]
		team := vars["team"]
		player := vars["player"]

		last, first, middle := parseName(r.FormValue("name"))

		height := r.FormValue("height")
		weight := r.FormValue("weight")
		hand := r.FormValue("hand")
		position := r.FormValue("position")
		jersey := r.FormValue("jersey")
		birth := r.FormValue("birth")

		_, err := config.Database.Exec(
			PlayerUpdate, first, middle, last, height, weight, hand, position,
			jersey, birth, league, team, player,
		)

		if err != nil {
			log.Println(err)
		}

	case http.MethodDelete:
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)

	}

} // playerAPIHandler
