package main

import (
  "log"
  "net/http"
  
)

func betaAPIHandler(w http.ResponseWriter, r *http.Request) {
  
  switch r.Method {
  case http.MethodPost:

    email     := r.FormValue("email")
    
    if email == "" {
      w.WriteHeader(http.StatusBadRequest)        
    } else {

      _, err := config.Database.Exec(
        BetaUserCreate, email,
      )
      
      if err != nil {
        
        log.Println(err)
        w.WriteHeader(http.StatusConflict)
        
      }
      
    }
    
  case http.MethodGet:
  case http.MethodPut:      
  case http.MethodDelete:
  default:
    w.WriteHeader(http.StatusMethodNotAllowed)
  
  }
    
} // betaAPIHandler
