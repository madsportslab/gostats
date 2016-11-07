package main

import (
	"log"
	"net/http"
)

func authAPIHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:

		u := authenticate(r)

		if u == nil {
			w.WriteHeader(http.StatusUnauthorized)
		}
		

	case http.MethodPut:

		email := r.FormValue("email")
		password := r.FormValue("password")

		if email == "" || password == "" {
			w.WriteHeader(http.StatusBadRequest)
		} else {

			u := User{}

			row := config.Database.QueryRow(
				"SELECT email, password, salt from users WHERE email = $1", email)

			err := row.Scan(&u.Email, &u.Password, &u.Salt)

			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusUnauthorized)
			} else {

				if u.Password == hashPassword(password, u.Salt, 32) {

					token := generateToken(32)

					_, err := config.Database.Exec(
						UserUpdateToken, token, email,
					)

					if err != nil {
						log.Println(err)
						w.WriteHeader(http.StatusConflict)
						return
					}

					cookie := http.Cookie{
						Name:   Madsportslab,
						Value:  token,
						Domain: *domain,
						Path:   "/",
					}

					http.SetCookie(w, &cookie)

					log.Println(token)

					w.Write(createJSON(token))

				} else {
					w.WriteHeader(http.StatusUnauthorized)
				}

			}

		}

	case http.MethodDelete:

		cookie := http.Cookie{
			Name:   Madsportslab,
			Value:  "",
			Domain: *domain,
			Path:   "/",
			MaxAge: -1,
		}

		http.SetCookie(w, &cookie)

	case http.MethodPost:
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

} // authAPIHandler
