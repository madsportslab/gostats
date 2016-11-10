package main

import (
  "database/sql"
  "fmt"
  "log"
	"net/http"
  "os"

  "github.com/sendgrid/sendgrid-go"
  "github.com/sendgrid/sendgrid-go/helpers/mail"
)

func addForgot(id string, token string) bool {

  row := config.Database.QueryRow(
		ForgotGet, id)

	f := Forgot{}

	err := row.Scan(&f.ID, &f.UserID, &f.Token)

	if err == sql.ErrNoRows {
	
    _, err2 := config.Database.Exec(
      ForgotCreate, id, token)
    
    if err2 != nil {
      log.Println("addForgot:Create: ", err2)
      return false
    } else {
      return true
    }

  } else {

    _, err3 := config.Database.Exec(
      ForgotUpdate, token, id, token)
    
    if err3 != nil {
      log.Println("addForgot:Update: ", err3)
      return false
    } else {
      return true
    }

  }
  
} // addForgot

func getForgot(token string) *Forgot {

  row := config.Database.QueryRow(
		ForgotExists, token)

	f := Forgot{}

	err := row.Scan(&f.ID, &f.UserID, &f.Token)

	if err == sql.ErrNoRows {
    return nil
  } else {
    return &f
  }
  
} // getForgot

func resetPasswordContent(token string) string {

  link := fmt.Sprintf("https://madsportslab.com/forgot?token=%s", token)

  return fmt.Sprintf("Dear user,<br><br>You have requested to reset your password " +
    "for madsportslab.com, if you did not initiate this request then " +
    "we suggest that you forward this email to support@madsportslab.com.  " +
    "<br><br>Click the link below to reset your password: " +
    "<br><br><a href=\"%s\">%s</a><br><br>Sincerely,<br>Team madsportslab" +
    "<br><br><a href=\"https://madsportslab.com\">https://madsportslab.com</a>",
    link, link)

} // resetPasswordContent

// TODO: text/html or plaintext parameter
func sendEmail(email string, text string) {

  from := mail.NewEmail("no-reply@madsportslab.com", "no-reply@madsportslab.com")
  subject := "Reset Password madsportslab"
  to := mail.NewEmail(email, email)
  content := mail.NewContent("text/html", text)

  m := mail.NewV3MailInit(from, subject, to, content)

  request := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"),
    "/v3/mail/send", "https://api.sendgrid.com")
  
  request.Method = "POST"
  request.Body = mail.GetRequestBody(m)

  res, err := sendgrid.API(request)

  if err != nil {
    log.Println(err)
  } else {
    
    log.Println(res.StatusCode)
    log.Println(res.Body)
    log.Println(res.Headers)

  }

} // sendEmail

func forgotAPIHandler(w http.ResponseWriter, r *http.Request) {

  u := authenticate(r)

	switch r.Method {
	case http.MethodGet:

    log.Println(u)
    if u != nil {
      w.WriteHeader(http.StatusForbidden)
      return
    }

    email := r.FormValue("email")

    log.Println(email)

    u := getUserByEmail(email)

    if u == nil {
      w.WriteHeader(http.StatusConflict)
      return
    }

    token := generateToken(32)

    content  := resetPasswordContent(token)

    log.Println(token)

    sendEmail(email, content)

    if !addForgot(u.ID, token) {
      w.WriteHeader(http.StatusConflict)
    }

  case http.MethodPost:
	case http.MethodDelete:
	case http.MethodPut:

	  token := r.FormValue("token")

    log.Println(token)

    f := getForgot(token)

    if f == nil {
      w.WriteHeaeder(http.StatusForbidden)
      return 
    }

    u := getUser(f.UserID)

    if u == nil {
      w.WriteHeader(http.StatusForbidden)
    } else {
      
      password  := r.FormValue("password")

      if password == "" {
        w.WriteHeader(http.StatusConflict)
      } else {

        hash  := hashPassword(password, u.Salt, 32)

        if updatePassword(u, hash) {

          _, err := config.Database.Exec(
            ForgotDelete, u.ID, token,
          )

          if err != nil {
            log.Println("PutForgotAPI:ForgotDelete:", err)
            w.WriteHeader(http.StatusConflict)
          }

        } else {
          w.WriteHeader(http.StatusConflict)
        }
        
      }  

    }
	

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)

	}

} // forgotAPIHandler
