package main

import (
  "fmt"
  "log"
	"net/http"
  "os"

  "github.com/sendgrid/sendgrid-go"
  "github.com/sendgrid/sendgrid-go/helpers/mail"
)

func generateLink() string {

  token := generateRandomHex(32)

  return fmt.Sprintf("https://madsportslab.com/forgot?token=%s",
    token)

} // generateLink

func resetPasswordContent() string {

  link := generateLink()

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

    if u != nil {
      w.WriteHeader(http.StatusForbidden)
      return
    }

    email := r.FormValue("email")

    if email == "" {
      w.WriteHeader(http.StatusConflict)
      return
    }

    content := resetPasswordContent()

    sendEmail(email, content)

  case http.MethodPost:
	case http.MethodDelete:
	case http.MethodPut:
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)

	}

} // forgotAPIHandler
