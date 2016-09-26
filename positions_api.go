package main

import (
  "encoding/json"
  "log"
  "net/http"
  
  "github.com/gorilla/mux"
)


func positionAPIHandler(w http.ResponseWriter, r *http.Request) {

  vars      := mux.Vars(r)
  sport     := vars["sport"]
  
  log.Println(sport)
  // TODO: check invalid sport_id
  
  switch r.Method {
  case http.MethodGet:
  
    rows, err := config.Database.Query(
      PositionGetAll, sport,
    )

    if err != nil {
      log.Println(err)
    }

    defer rows.Close()

    positions := []Position{}

    for rows.Next() {
      
      p := Position{}
      
      err := rows.Scan(&p.ID, &p.Name, &p.Short)

      if err != nil {
        log.Println(err)
      }
      
      log.Println(p)
      
      positions = append(positions, p)

    }
    
    j, err := json.Marshal(positions)
    
    if err != nil {
      log.Println(err)
      w.WriteHeader(http.StatusInternalServerError)
      return
    }
    
    w.Write(j)
    
  case http.MethodPut:
  case http.MethodDelete:
  case http.MethodPost:
  default:
    w.WriteHeader(http.StatusMethodNotAllowed)
  }
    
} // positionAPIHandler
