package main

import (
  "bytes"
  //"bufio"
  "io"
  "log"
  "net/http"
  //"strings"

  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/s3/s3manager"
  //"github.com/gorilla/mux"
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

  //vars := mux.Vars(r)

	switch r.Method {
	case http.MethodPost:

    //id := vars["id"]
    
    // seems order of multipart items is important
    
    /*
    err := r.ParseMultipartForm(5 << 20)
    
    if err != nil {
      log.Println(err)
      w.WriteHeader(http.StatusInternalServerError)
      return
    }*/
    /*
    token := r.FormValue("token")

    log.Println(token)

    file, header, err2 := r.FormFile("file")

    if err2 != nil {
      log.Println(err2)
      w.WriteHeader(http.StatusConflict)
      return
    }

    log.Println(file)
    log.Println(header)
    */

    form, err := r.MultipartReader()

    if err != nil {
      log.Println(err)
      return
    }

    for {

      part, err := form.NextPart()

      if err == io.EOF {
        break
      }

      log.Println(part.FormName())

      buf := new(bytes.Buffer)
      buf.ReadFrom(part)
        
      if part.FormName() == "iconfile" {
        s3Upload(buf, "test")
      } else if part.FormName() == "token" {

        log.Println(buf.String())

      }

    }

    
/*
    for _, fhd := range r.MultipartForm.File {

      for _, hdr := range fhd {

        log.Println("shit")
        log.Println(hdr)
        //buffer := bufio.NewReader(hdr)

      }
    }
*/

/*
    if err != nil {
      log.Println("FormFile error: ", err)
      w.WriteHeader(http.StatusInternalServerError)
      return
    }

    defer file.Close()

    buffer := bufio.NewReader(file)

    if err != nil {
      log.Println(err)
      w.WriteHeader(http.StatusConflict)
    }
*/
/*
    s3Upload(buffer, header.Filename, header.Header.Get("Content-Type"))
    
    if err != nil {
      log.Println(err)
      w.WriteHeader(http.StatusConflict)
      return
    }

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
      q, header.Filename, id,
    )

    if err2 != nil {

      log.Println(err2)
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
