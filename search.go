package main

import (
  "database/sql"
  "encoding/json"
  "log"
  "net/http"
  
  "github.com/gorilla/mux"
)

func searchLeagues(user_id string) []League {

  rows, err := config.Database.Query(
    UserLeagueGetAll, user_id,
  )

  if err != nil {
    log.Println(err)
  }

  defer rows.Close()

  leagues := []League{}

  for rows.Next() {

    l := League{}
        
    err := rows.Scan(&l.ID, &l.Name,&l.Canonical, &l.Icon)

    if err == sql.ErrNoRows || err != nil {
      log.Println(err)
      return nil
    }
          
    leagues = append(leagues, l)

  }
  
  return leagues
  
} // searchLeagues


func searchHandler(w http.ResponseWriter, r *http.Request) {

  vars := mux.Vars(r)
  
  league    := vars["league"]
  
  u := authenticate(r)
  
  if u == nil {
    log.Println("leagueHandler: token not found")
  }
  
  switch r.Method {
  case http.MethodPost:
    
    name        := r.FormValue("name")
    sport       := r.FormValue("sport")
   
    res, err := config.Database.Exec(
      LeagueCreate, name, generateCanonical(name), sport,
    )
    
    if err != nil {
      log.Println(err)
      w.WriteHeader(http.StatusConflict)
    }

    league_id, err2 := res.LastInsertId()
    
    if err2 != nil {
      log.Println(err2)
      w.WriteHeader(http.StatusConflict)
    }
    
    _, err3 := config.Database.Exec(
      LeagueAdminCreate, league_id, u.ID,
    )

    if err3 != nil {
      log.Println(err3)
      w.WriteHeader(http.StatusConflict)
    }
    
  case http.MethodGet:
        
    if league != "" {

      l := getLeague(league)
      
      if l == nil {
        w.WriteHeader(http.StatusNotFound)
      } else {
        j, _ := json.Marshal(l)
        w.Write(j)        
      }

    } else {
    
      leagues := searchLeagues("")
      
      if len(leagues) == 0 {
        w.WriteHeader(http.StatusNotFound)
      } else {
        j, _ := json.Marshal(leagues)
        w.Write(j)        
      }
      
    }

  case http.MethodPut:
  case http.MethodDelete:
  default:
    w.WriteHeader(http.StatusMethodNotAllowed)
  }
    
} // searchHandler
