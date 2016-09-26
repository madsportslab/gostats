package main

import (
  "fmt"
  "net/http"
  "net/http/httptest"
  "strings"
  "testing"
  
)

var server *httptest.Server
var userEndpoint string

func TestCreateUser(t *testing.T) {
  
  req, err := http.NewRequest("POST", userEndpoint,
    strings.NewReader("email=test@test.com&password=test"))
  req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  
  res, err := http.DefaultClient.Do(req)
  
  if err != nil {
    t.Error(err)
  }
  
  if res.StatusCode != http.StatusOK {
    t.Errorf("Success expected: %d", res.StatusCode)
  }
  
} // TestCreateUser


func init() {

  server = httptest.NewServer(initRoutes())
  
  initDatabase()
  
  userEndpoint = fmt.Sprintf("%s/api/users", server.URL)
  
} // init
