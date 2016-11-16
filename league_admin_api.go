package main

import (
  "database/sql"
  "encoding/json"
  "log"
	"net/http"

	"github.com/gorilla/mux"
)

func getLeagueAdmins(league string) []LeagueAdmin {

  rows, err := config.Database.Query(
		LeagueAdminGetAll, league,
	)

	if err != nil {
		log.Println("getLeagueAdmins: ", err)
		return nil
	}

	defer rows.Close()

	admins := []LeagueAdmin{}

	for rows.Next() {

		la := LeagueAdmin{}

		err := rows.Scan(&la.UserID, &la.Email, &la.LeagueID)

		if err == sql.ErrNoRows || err != nil {
			log.Println("getLeagueAdmins: ", err)
			return nil
		}

		admins = append(admins, la)

	}

	return admins

} // getLeagueAdmins

func addLeagueAdmin(league string, email string) bool {

  l := getLeague(league)
  u := getUserByEmail(email)

  if l == nil || u == nil {
    return false
  } else {

    row := config.Database.QueryRow(
      LeagueAdminGet, league, u.ID,
    )

    la := LeagueAdmin{}

    err := row.Scan(&la.LeagueID, &la.UserID)

    if err != sql.ErrNoRows {
      log.Println(err)
      return false
    }

    _, err2 := config.Database.Exec(
      LeagueAdminCreate, league, u.ID,
    )

    if err2 != nil {
      log.Println("addLeagueAdmin: ", err2)
      return false
    }

    return true
    
  }

} // addLeagueAdmin

func leagueAdminHandler(w http.ResponseWriter, r *http.Request) {

  u := authenticate(r)

  vars := mux.Vars(r)
  
  league := vars["league"]

  if !isLeagueAdmin(u, league) {
    w.WriteHeader(http.StatusForbidden)
    return
  }

	switch r.Method {
	case http.MethodPost:
	  
    email := r.FormValue("email")
    
    if email == "" {
      w.WriteHeader(http.StatusConflict)
      return
    }

    if !addLeagueAdmin(league, email) {
      w.WriteHeader(http.StatusConflict)
      return
    }

	case http.MethodPut:
	case http.MethodDelete:
    // remove administrator

    admin := vars["admin"]

    if admin == u.ID {
      w.WriteHeader(http.StatusForbidden)
      return
    }
    
    _, err := config.Database.Exec(
      LeagueAdminDelete, admin, league,
    )

    if err != nil {
      w.WriteHeader(http.StatusConflict)
    }

	case http.MethodGet:
    // check if administrator

    admins := getLeagueAdmins(league)

    j, _ := json.Marshal(admins)
      
    w.Write(j)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

} // leagueAdminHandler
