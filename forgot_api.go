package main

import (
  "log"
	"net/http"

  //"github.com/sendgrid/sendgrid-go"
	"github.com/gorilla/mux"
)

func forgotAPIHandler(w http.ResponseWriter, r *http.Request) {

  u := authenticate(r)

	vars := mux.Vars(r)

	league := vars["league"]
	log.Println(league)

	switch r.Method {
	case http.MethodGet:

    if u != nil {
      w.WriteHeader(http.StatusForbidden)
      return
    }

    email := r.FormValue("email")

    if email == "" {
      w.WriteHeader(http.StatusConflict)
      return
    }



  case http.MethodPost:
	case http.MethodDelete:
	case http.MethodPut:
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)

	}

} // forgotAPIHandler
