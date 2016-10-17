package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func getGames(league_id string, scheduled string) []Game {

	league := getLeague(league_id)

	season := getLatestSeasonByID(league_id)

	teams := getTeamsAsMap(league)

	rows, err := config.Database.Query(
		ScheduleGetAll, league_id, season.ID,
	)

	if err != nil {
		log.Println("getGames: ", err)
		return nil
	}

	defer rows.Close()

	games := []Game{}

	for rows.Next() {

		g := Game{}

		err := rows.Scan(&g.ID, &g.Scheduled, &g.HomeID, &g.AwayID,
			&g.LeagueID, &g.SeasonID, &g.Completed)

		if err == sql.ErrNoRows || err != nil {
			log.Println("getGames: ", err)
			return nil
		}

		//g.URL = fmt.Sprintf("/viewgame?leagueid=%s&seasonid=%s&gameid=%s",
		//	g.LeagueID, g.SeasonID, g.ID)

		g.AwayName = teams[g.AwayID].Name
		g.HomeName = teams[g.HomeID].Name
		g.Token = getGameToken(g.ID)

		games = append(games, g)

	}

	return games

} // getGames

func getGamesEx(user *User) map[string][]Game {

	ret := map[string][]Game{}

	leagues := getLeagues(user)

	for _, league := range leagues {

		season := getLatestSeasonByID(league.ID)

		teams := getTeamsAsMap(&league)

		rows, err := config.Database.Query(
			ScheduleGetAll, league.ID, season.ID,
		)

		if err != nil {
			log.Println("getGamesEx: ", err)
			return nil
		}

		defer rows.Close()

		games := []Game{}

		for rows.Next() {

			g := Game{}

			err := rows.Scan(&g.ID, &g.Scheduled, &g.HomeID, &g.AwayID,
				&g.LeagueID, &g.SeasonID, &g.Completed)

			if err == sql.ErrNoRows || err != nil {
				log.Println("getGamesEx: ", err)
				return nil
			}

			//g.URL = fmt.Sprintf("/viewgame?leagueid=%s&seasonid=%s&gameid=%s",
			//	g.LeagueID, g.SeasonID, g.ID)

			g.AwayName = teams[g.AwayID].Name
			g.HomeName = teams[g.HomeID].Name
			g.Token = getGameToken(g.ID)

			games = append(games, g)

		}

		ret[league.Name] = games

	}

	return ret

} // getGamesEx

func getGamesBySeason(leagueid string, seasonid string, teamid string) []Game {

	rows, err := config.Database.Query(
		TeamScheduleGetAll, leagueid, seasonid, teamid, teamid, teamid,
	)

	if err != nil {
		log.Println("getGamesBySeason: ", err)
		return nil
	}

	defer rows.Close()

	games := []Game{}

	for rows.Next() {

		g := Game{}

		err := rows.Scan(&g.ID, &g.HomeID, &g.AwayID, &g.LeagueID, &g.SeasonID,
			&g.Completed, &g.Opponent)

		if err == sql.ErrNoRows || err != nil {
			log.Println("getGamesBySeason: ", err)
			return nil
		}

  req := Req{
		LeagueId: g.LeagueID,
		SeasonId: g.SeasonID,
		GameId: g.ID,
	}

  score1 := getTeamScore(&req, g.HomeID)
	score2 := getTeamScore(&req, g.AwayID)

	log.Println(score1)
	log.Println(score2)

	g.HomeScore, _ = strconv.ParseInt(score1["T"], 0, 64) 
	g.AwayScore, _ = strconv.ParseInt(score2["T"], 0, 64)

		g.URL = fmt.Sprintf("/viewgame?leagueid=%s&seasonid=%s&gameid=%s&myteamid=%s",
			g.LeagueID, g.SeasonID, g.ID, teamid)

		g.Token = getGameToken(g.ID)

		games = append(games, g)

	}

	return games

} // getGamesBySeason

func getGame(id string) *Game {

	row := config.Database.QueryRow(
		ScheduleGet, id,
	)

	g := Game{}

	err := row.Scan(&g.ID, &g.HomeID, &g.AwayID, &g.LeagueID, &g.SeasonID,
		&g.Completed)

	if err == sql.ErrNoRows {
		log.Println("getGame: ", err)
		return nil
	}

  req := Req{
		LeagueId: g.LeagueID,
		SeasonId: g.SeasonID,
		GameId: g.ID,
	}

  score1 := getTeamScore(&req, g.HomeID)
	score2 := getTeamScore(&req, g.AwayID)

	log.Println(score1)
	log.Println(score2)

	g.HomeScore, _ = strconv.ParseInt(score1["T"], 0, 64) 
	g.AwayScore, _ = strconv.ParseInt(score2["T"], 0, 64)

	g.URL = fmt.Sprintf("/viewgame?leagueid=%s&seasonid=%s&gameid=%s",
		g.LeagueID, g.SeasonID, g.ID)

	g.Token = getGameToken(g.ID)

	return &g

} // getGame

func endGame(league string, game string) {

		_, err := config.Database.Exec(
			ScheduleFinal, league, game,
		)

		if err != nil {
			log.Println("endGame: ", err)
		}

} // endGame

func scheduleAPIHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

  league := vars["league"]

	switch r.Method {
	case http.MethodPost:

		home := r.FormValue("home")
		away := r.FormValue("away")

		s := getLatestSeasonByID(league)

		if s == nil {
			return
		}

		_, err := config.Database.Exec(
			ScheduleCreate, home, away, s.ID, league,
		)

		if err != nil {

			log.Println("post scheduleAPIHandler: ", err)
			w.WriteHeader(http.StatusConflict)
			return

		}

	case http.MethodGet:

		// TODO: auth
		team := vars["team"]

		game := vars["game"]

		if game != "" {

			g := getGame(game)
			
			if g == nil {
				w.WriteHeader(http.StatusNotFound)
				return
			} else {
				j, _ := json.Marshal(g)
				w.Write(j)
			}

		} else {

			if team != "" {

				t := getTeam(team)

				s := getLatestSeasonByID(league)

				if s == nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				games := getGamesBySeason(league, s.ID, t.ID)

				log.Println(games)

				if len(games) == 0 {
					w.WriteHeader(http.StatusNotFound)
				} else {
					j, _ := json.Marshal(games)
					w.Write(j)
				}

			} else {

				u := authenticate(r)

				if u == nil {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				log.Println(u.ID)

				games := getGamesEx(u)

				if len(games) == 0 {
					w.WriteHeader(http.StatusNotFound)
				} else {
					j, _ := json.Marshal(games)
					w.Write(j)
				}

			}

		}

	case http.MethodPut:

		// TODO: auth

		game := vars["game"]

		g := getGame(game)

		req := Req{
			LeagueId: g.LeagueID,
			SeasonId: g.SeasonID,
			GameId: g.ID,
			HomeId: g.HomeID,
			AwayId: g.AwayID,

		}

    gameFinal(&req)

		endGame(league, game)

	case http.MethodDelete:
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)

	}

} // scheduleAPIHandler
