package main

import (
  //"encoding/json"
  "log"
  "net/http"
  
  //"github.com/gorilla/mux"
)

func updatePassword(u *User, hash string) bool {

  _, err := config.Database.Exec(
    UserUpdatePassword, hash, u.ID, 
  )

  if err != nil {
    log.Println("updatePassword: ", err)
    return false
  } else {
    return true
  }

} // updatePassword

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


func getUser(id string) *User {
  
  row := config.Database.QueryRow(
    UserGet, id,
  )

  u := User{}
  
  err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Icon, &u.Salt)
  
  if err != nil {
    log.Println(err)
    return nil
  }
  
  return &u
  
} // getUser

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
  case http.MethodPut:  
  
    u := authenticate(r)
    
    if u == nil {
      w.WriteHeader(http.StatusUnauthorized)
      return
    }    
    
    leagueid := r.FormValue("leagueid")
    
    _, err2 := config.Database.Exec(
      UserUpdateDefaultLeague, leagueid, u.ID,
    )
    
    if err2 != nil {
      log.Println(err2)
      w.WriteHeader(http.StatusInternalServerError)
    }
    
  case http.MethodGet:
  case http.MethodDelete:
  default:
    w.WriteHeader(http.StatusMethodNotAllowed)
  
  }
    
} // userAPIHandler
