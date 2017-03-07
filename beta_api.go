package main

import (
  "encoding/json"
  "log"
  "net/http"
  
)

func betaAPIHandler(w http.ResponseWriter, r *http.Request) {
  
  switch r.Method {
  case http.MethodPost:

    email := r.FormValue("email")
    
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

      j, err1 := json.Marshal(nil)

      if err1 != nil {
        log.Println(err1)
      } else {
        w.Write(j)
      }

      
    }
    
  case http.MethodGet:
  case http.MethodPut:      
  case http.MethodDelete:
  default:
    w.WriteHeader(http.StatusMethodNotAllowed)
  
  }
    
} // betaAPIHandler
