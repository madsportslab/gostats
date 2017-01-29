package main

import (
  "bytes"
  "log"
  "net/http"
  "strings"

  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/s3/s3manager"
  "github.com/gorilla/mux"
)

func updateLeagueIcon() {


} // updateLeagueIcon

func s3Upload(buf *bytes.Buffer, filename string) {

  config := &aws.Config{
    Region: aws.String("us-west-1"),
  }

  uploader := s3manager.NewUploader(session.New(config))
  
  res, err := uploader.Upload(&s3manager.UploadInput{
    Body: buf,
    Bucket: aws.String("madsportslab-dev1"),
    Key: aws.String(filename),
    ACL: aws.String("public-read"),
  })

  if err != nil {
    log.Println("s3Upload: ", err)
  }

  log.Println(res)

} // s3Upload

func iconAPIHandler(w http.ResponseWriter, r *http.Request) {

  //u := authenticate(r)

  vars := mux.Vars(r)

	switch r.Method {
	case http.MethodPost:

    id := vars["id"]
    
    // seems order of multipart items is important

    form, err := r.MultipartReader()

    if err != nil {
      log.Println(err)
      return
    }

    p1, err := form.NextPart()

    b1 := new(bytes.Buffer)
    b1.ReadFrom(p1)

    token := b1.String()

    log.Println(token)

    p2, err := form.NextPart()

    b2 := new(bytes.Buffer)
    b2.ReadFrom(p2)

    filename := b2.String()

    p3, err := form.NextPart()

    b3 := new(bytes.Buffer)
    b3.ReadFrom(p3)

    s3Upload(b3, filename)

    var q = ""
    
    if strings.Contains(r.URL.String(), "leagues") {
      q = LeagueIconUpdate
    } else if strings.Contains(r.URL.String(), "teams") {
      q = TeamIconUpdate
    } else if strings.Contains(r.URL.String(), "players") {
      q = PlayerIconUpdate
    } else {
      w.WriteHeader(http.StatusConflict)
      return
    }
    
    // check privilege
    
    _, err2 := config.Database.Exec(
      q, filename, id,
    )

    if err2 != nil {

      log.Println(err2)
      w.WriteHeader(http.StatusConflict)

    }

  case http.MethodGet:
  case http.MethodPut:
  case http.MethodDelete:
  default:
    w.WriteHeader(http.StatusMethodNotAllowed)
  }

} // iconAPIHandler
