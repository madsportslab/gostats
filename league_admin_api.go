package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func leagueAdminHandler(w http.ResponseWriter, r *http.Request) {

  u := authenticate(r)

  vars := mux.Vars(r)
  
  league := vars["league"]

	switch r.Method {
	case http.MethodPost:
	  // add administrator
	case http.MethodPut:
	case http.MethodDelete:
    // remove administrator
	case http.MethodGet:
    // check if administrator

    if isLeagueAdmin(u, league) {
      w.WriteHeader(http.StatusOK)
    } else {
      w.WriteHeader(http.StatusForbidden)
    }

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

} // leagueAdminHandler
