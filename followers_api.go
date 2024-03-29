package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func getLeagueFollowers(league string) []Follower {

  rows, err := config.Database.Query(
		LeagueFollowerGetAll, league,
	)

	if err != nil {
		log.Println("getLeagueFollowers: ", err)
		return nil
	}

	defer rows.Close()

	followers := []Follower{}

	for rows.Next() {

		f := Follower{}

		err := rows.Scan(&f.ID, &f.LeagueID, &f.UserID)

		if err == sql.ErrNoRows || err != nil {
			log.Println("getLeagueFollowers: ", err)
			return nil
		}

		followers = append(followers, f)

	}

	return followers

} // getLeagueFollowers

func isFollower(user *User, leagueid string) bool {

	row := config.Database.QueryRow(
		LeagueFollowerGet, leagueid, user.ID)

	f := Follower{}

	err := row.Scan(&f.LeagueID, &f.UserID)

	if err == sql.ErrNoRows {
		log.Println("isFollower: ", err)
		return false
	} else {
		return true
	}

} // isFollower

func followerAPIHandler(w http.ResponseWriter, r *http.Request) {

  u := authenticate(r)

	vars := mux.Vars(r)

	league := vars["league"]
	
	switch r.Method {
	case http.MethodPost:

	  // check if exist first

		_, err := config.Database.Exec(
			LeagueFollowerGet, league, u.ID,
		)

		if err == nil {

			log.Println("follower Add: err is nil")
			_, err2 := config.Database.Exec(
				LeagueFollowerCreate, league, u.ID,
			)

			if err2 != nil {
				log.Println("LeagueFollowerCreate: ", err2)
				w.WriteHeader(http.StatusConflict)
				return
			}
			
		} else {
			w.WriteHeader(http.StatusConflict)
      return
		}

	case http.MethodGet:

    followers := getLeagueFollowers(league)

    admin 		:= isLeagueAdmin(u, league)
		follower 	:= isFollower(u, league)

    f := Follow{
      Followers: followers,
      IsAdmin: admin,
			IsFollower: follower,
    }

		j, _ := json.Marshal(f)

		w.Write(j)

	case http.MethodDelete:

    _, err := config.Database.Exec(
			LeagueFollowerDelete, league, u.ID,
		)

		if err != nil {
			log.Println("LeagueFollowerDelete: ", err)
			w.WriteHeader(http.StatusConflict)
		}

	case http.MethodPut:
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)

	}

} // followerAPIHandler
