package main

import (
  "encoding/json"
  "log"
  "net/http"
  
  "github.com/gorilla/mux"
)


func getUserByEmail(email string) *User {
  
  row := config.Database.QueryRow(
    UserGetByEmail, email,
  )

  u := User{}
  
  err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Icon)
  
  if err != nil {
    log.Println(err)
    return nil
  }
  
  return &u
  
} // getUserByEmail


func userAPIHandler(w http.ResponseWriter, r *http.Request) {
  
  switch r.Method {
  case http.MethodPost:

    email     := r.FormValue("email")
    password  := r.FormValue("password")
    
    if email == "" || password == "" {
      w.WriteHeader(http.StatusBadRequest)        
    } else {

      salt  := generateRandomHex(16)
      
      hash  := hashPassword(password, salt, 32)
      
      token := generateToken(32)
                  
      _, err := config.Database.Exec(
        UserCreate, email, hash, salt, token,
      )
      
      if err != nil {
        
        log.Println(err)
        w.WriteHeader(http.StatusConflict)
        
      } else {

        cookie := http.Cookie{
          Name: Madsportslab,
          Value: token,
          Domain: "127.0.0.1",
          Path: "/",
        }
        
        http.SetCookie(w, &cookie)
                
        w.Write(createJSON("Account created successfully"))
                        
      }
      
    }
    
  case http.MethodGet:
  
    // TODO: auth
    vars := mux.Vars(r)
    
    user := vars["user"]
    
    row := config.Database.QueryRow(
      UserGet, user,
    )
  
    u := User{}
    
    err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Icon)
    
    if err != nil {
      log.Println(err)
    }
    
    j, _ := json.Marshal(u)
    
    w.Write(j)
    
  case http.MethodPut:  
  
    u := authenticate(r)
    
    if u == nil {
      log.Println(u)
      w.WriteHeader(http.StatusUnauthorized)
      w.Write([]byte("Must login as user."))
      return
    }    
    
    leagueid      := r.FormValue("leagueid")
    
    _, err2 := config.Database.Exec(
      UserUpdateDefaultLeague, leagueid, u.ID,
    )
    
    if err2 != nil {
      log.Println(err2)
      w.WriteHeader(http.StatusInternalServerError)
    }
    
  case http.MethodDelete:
  default:
    w.WriteHeader(http.StatusMethodNotAllowed)
  
  }
    
} // userAPIHandler
