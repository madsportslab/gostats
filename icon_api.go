package main

import (
  "bufio"
  "log"
  "net/http"
  //"strings"

  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/s3/s3manager"
  "github.com/gorilla/mux"
)

func updateLeagueIcon() {


} // updateLeagueIcon

func s3Upload(buffer *bufio.Reader, filename string, content string) {

  uploader := s3manager.NewUploader(session.New())
  
  res, err := uploader.Upload(&s3manager.UploadInput{
    Body: buffer,
    Bucket: aws.String("madsportslab-dev1"),
    Key: aws.String(filename),
  })

  if err != nil {
    log.Println(err)
  }

  log.Println(res)

} // s3Upload

func iconAPIHandler(w http.ResponseWriter, r *http.Request) {

  u := authenticate(r)

  log.Println(u)
  
	vars := mux.Vars(r)

	switch r.Method {
	case http.MethodPost:

    var statement = ""

    id := vars["id"]
    log.Println(id)

    // seems order of multipart items is important

    file, header, err := r.FormFile("file")
    log.Println(file)
    log.Println(header.Filename)
    log.Println(header.Header)
    log.Println(header.Header.Get("Content-Type"))

    defer file.Close()

    buffer := bufio.NewReader(file)

    if err != nil {
      log.Println(err)
      w.WriteHeader(http.StatusConflict)
    }

    s3Upload(buffer, header.Filename, header.Header.Get("Content-Type"))
    //log.Println(file, header)

    if err != nil {
      log.Println(err)
      w.WriteHeader(409)
    }

    log.Println(statement)
    
    /*
    if strings.Contains(r.URL.String(), "leagues") {
      statement = LeagueIconUpdate
    } else if strings.Contains(r.URL.String(), "teams") {
      statement = TeamIconUpdate
    } else if strings.Contains(r.URL.String(), "players") {
      statement = PlayerIconUpdate
    } else {
      w.WriteHeader(http.StatusConflict)
      return
    }
    
    // check privilege

    _, err := config.Database.Exec(
      statement, icon, id,
    )

    if err != nil {

      log.Println(err)
      w.WriteHeader(http.StatusConflict)

    }
    */


  case http.MethodGet:
  case http.MethodPut:
  case http.MethodDelete:
  default:
    w.WriteHeader(http.StatusMethodNotAllowed)
  }

} // iconAPIHandler
