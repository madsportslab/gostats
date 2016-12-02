package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func parseName(name string) (string, string, string) {

	if name == "" {
		return "", "", ""
	}

	names := strings.Split(name, " ")

	l := len(names)

	if l == 1 {
		return names[0], "", ""
	} else if l == 2 {
		return names[1], names[0], ""
	} else if l == 3 {
		return names[1], names[0], names[2]
	} else {
		return "", "", ""
	}

} // parseName

func hashPassword(password string, salt string, length int) string {

	digest := hmac.New(sha256.New, config.HashKey)

	digest.Write([]byte(password + salt + config.Secret))

	hash := hex.EncodeToString(digest.Sum(nil))

	return hash[:length]

} // hashPassword

func generateRandomHex(length int) string {

	buf := make([]byte, length)

	_, err := rand.Read(buf)

	if err != nil {
		return ""
	} else {
		return hex.EncodeToString(buf)
	}

} // generateRandomHex

func generateToken(length int) string {

	block, err := aes.NewCipher(config.BlockKey)

	if err != nil {
		return ""
	} else {

		cfb := cipher.NewCFBEncrypter(block, config.IV)

		plaintext := []byte(time.Now().String() + config.Secret +
			generateRandomHex(32))

		ciphertext := make([]byte, len(plaintext))

		cfb.XORKeyStream(ciphertext, plaintext)

		return hex.EncodeToString(ciphertext)

	}

} // generateToken

func tokenToUser(token string) *User {

	u := User{}

	row := config.Database.QueryRow(
		"SELECT id FROM users WHERE token=?", token)

	err := row.Scan(&u.ID)

	if err != nil {
		log.Println("tokenToId: ", err)
		return nil
	} else {
		return &u
	}

} // tokenToUser

func authenticate(r *http.Request) *User {

	var token string

	cookie, err := r.Cookie("madsportslab")

	if err != nil {

		token = r.FormValue("token")

		if token == "" {
			return nil
		}

	} else {
		token = cookie.Value
	}

	if token == "" {
		log.Println("authenticate: Token is an empty string")
		return nil
	}

	u := User{}

	row := config.Database.QueryRow(
		"SELECT id, icon, city, country, location, defaultLeague from users WHERE token=?", token)

	err2 := row.Scan(&u.ID, &u.Icon, &u.City, &u.Country, &u.Location, &u.DefaultLeague)

	if err2 != nil {
		log.Println(err2)
		return nil
	} else {
		return &u
	}

} // authenticate

func authorizeLeague(user *User, leagueid string) (*League, error) {

	row := config.Database.QueryRow(
		UserLeagueGet, leagueid, user.ID,
	)

	l := League{}

	err := row.Scan(&l.ID, &l.Name, &l.Canonical, &l.Icon)

	if err == sql.ErrNoRows {
		msg := fmt.Sprintf("authorize: %s", err)
		log.Println(msg)
		return nil, errors.New(msg)
	}

	return &l, nil

} // authorizeLeague

func authorizeTeam(user *User, teamid string) (*Team, error) {

	row := config.Database.QueryRow(
		UserTeamGet, teamid, user.ID,
	)

	t := Team{}

	err := row.Scan(&t.ID, &t.Name, &t.Canonical, &t.Icon)

	if err == sql.ErrNoRows {
		msg := fmt.Sprintf("authorize: %s", err)
		log.Println(msg)
		return nil, errors.New(msg)
	}

	return &t, nil

} // authorizeTeam

func isLeagueAdmin(user *User, leagueid string) bool {

	row := config.Database.QueryRow(
		LeagueAdminCheck, leagueid, user.ID)

	la := LeagueAdmin{}

	err := row.Scan(&la.LeagueID, &la.UserID)

	if err == sql.ErrNoRows {
		log.Println("isLeagueAdmin: ", err)
		return false
	} else {
		return true
	}

} // isLeagueAdmin

func isTeamAdmin(user *User, teamid string) bool {

	row := config.Database.QueryRow(
		TeamAdminCheck, teamid, user.ID)

	ta := TeamAdmin{}

	err := row.Scan(&ta.TeamID, &ta.UserID)

	if err == sql.ErrNoRows {
		log.Println("isTeamAdmin: ", err)
		return false
	} else {
		return true
	}

} // isTeamAdmin

func generateCanonical(text string) string {

	s := strings.Replace(text, " ", "", -1)

	ls := strings.ToLower(s)

	return ls

} // generateCanonical

func getGameToken(gameid string) string {

	row := config.Database.QueryRow(
		ScheduleGet, gameid,
	)

	g := Game{}

	err := row.Scan(&g.ID, &g.HomeID, &g.AwayID, &g.LeagueID, &g.SeasonID,
		&g.Completed)

	if err != nil {
		log.Println("getGameToken: ", err)
	}

	token := fmt.Sprintf("%s%s%s%s%s", g.ID, g.HomeID, g.AwayID, g.LeagueID, g.SeasonID)

	digest := hmac.New(sha256.New, config.HashKey)

	digest.Write([]byte(token))

	hash := hex.EncodeToString(digest.Sum(nil))

	return hash[:15]

} // getGameToken

func createJSON(msg string) []byte {

	jr := JSONResponse{
		Msg: msg,
	}

	j, err := json.Marshal(jr)

	if err != nil {
		log.Println("createJSON: ", err)
		return nil
	} else {
		return j
	}

} // createJSON

func extractDate(val string) string {

	log.Println("wow")
  log.Println(val)

	//const format = "2006-01-02 03:04:05 -0700"
  t, err := time.Parse(ShortForm, val)
	
	if err != nil {
		log.Println(err)
		return ""
	}

  d := int(t.Day())
	m := int(t.Month())

	f := fmt.Sprintf("%d-%s-%s", t.Year(), lpad(m), lpad(d))
  
	return f

} // extractDate

func lpad(val int) string {

  if val < 10 {
		return fmt.Sprintf("0%d", val)
	} else {
		return fmt.Sprintf("%d", val)
	}

} // lpad
