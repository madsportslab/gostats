package main

import (
  "encoding/json"
  "log"
  "net/http"
  
)

func sportAPIHandler(w http.ResponseWriter, r *http.Request) {
  
  switch r.Method {
  case http.MethodGet:
  
    rows, err := config.Database.Query(
      SportGet,
    )

    if err != nil {
      log.Println(err)
    }

    defer rows.Close()

    sports := []Sport{}

    for rows.Next() {
  
      s := Sport{}
          
      err := rows.Scan(&s.ID, &s.Name)

      if err != nil {
        log.Println(err)
      }
            
      sports = append(sports, s)

    }
    
    j, err := json.Marshal(sports)
    
    if err != nil {
      log.Println(err)
      w.WriteHeader(http.StatusInternalServerError)
      return
    }
    
    w.Write(j)
  
  case http.MethodDelete:
  case http.MethodPost:
  case http.MethodPut:
  default:
    w.WriteHeader(http.StatusMethodNotAllowed)
    
  }
    
} // sportAPIHandler
