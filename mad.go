package main

import (
	"database/sql"
)

const (
	MadSportsLab = "madsportslab"
)

// Database table structures

type User struct {
	ID            string         `json:"id"`
	Name          sql.NullString `json:"name"`
	Email         string         `json:"email"`
	Mobile        sql.NullString `json:"mobile"`
	Password      string         `json:"password"`
	Salt          string         `json:"salt"`
	Icon          sql.NullString `json:"icon"`
	Token         sql.NullString `json:"token"`
	Meta          sql.NullString `json:"meta"`
	DefaultLeague int            `json:"defaultLeague"`
	City          sql.NullString `json:"city"`
	Country       sql.NullString `json:"country"`
	Location      sql.NullString `json:"location"`
}

type Sport struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
}

type Position struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Short string `json:"short"`
}

type League struct {
	ID        string         `json:"id"`
	SportID   string         `json:"sport_id"`
	Name      string         `json:"name"`
	Canonical string         `json:"canonical"`
	Icon      sql.NullString `json:"icon"`
	Visible   bool           `json:"visible"`
	Official  bool           `json:"official"`
	Metric    bool           `json:"metric"`
	City      sql.NullString `json:"city"`
	Country   sql.NullString `json:"country"`
	Location  sql.NullString `json:"location"`
	URL       string         `json:"url"`
	IsAdmin   bool           `json:"isAdmin"`
}

type Team struct {
	ID              string         `json:"id"`
	Name            string         `json:"name"`
	Canonical       string         `json:"canonical"`
	Icon            sql.NullString `json:"icon"`
	LeagueID        string         `json:"leagueID"`
	LeagueCanonical string         `json:"leagueCanonical"`
	LeagueName      string         `json:"leagueName"`
	URL             string         `json:"url"`
}

type LeagueAdmin struct {
	LeagueID string `json:"leagueId"`
	UserID   string `json:"userId"`
}

type TeamAdmin struct {
	ID        string `json:"id"`
	LeagueID  string `json:"leagueId"`
	TeamID    string `json:"teamId"`
	TeamName  string `json:"teamName"`
	UserID    string `json:"userID"`
	UserEmail string `json:"userEmail"`
}

type Follower struct {
	ID				string `json:"id"`
	LeagueID	string `json:"leagueId"`
	UserID    string `json:"userId"`
}

type Follow struct {
	Followers []Follower	`json:"followers"`
	IsAdmin   bool				`json:"isAdmin"`
}

type Season struct {
	ID       string         `json:"id"`
	Games    int            `json:"games"`
	Periods  string         `json:"periods"`
	Duration int            `json:"duration"`
	Start    sql.NullString `json:"start"`
	Finish   sql.NullString `json:"finish"`
	LeagueID string         `json:"leagueID"`
}

type Game struct {
	ID              string         `json:"id"`
	Scheduled       sql.NullString `json:"scheduled"`
	HomeID          string         `json:"homeID"`
	AwayID          string         `json:"awayID"`
	HomeName        string         `json:"homeName"`
	AwayName        string         `json:"awayName"`
	Opponent        string         `json:"opponent"`
	Completed       bool           `json:"completed"`
	HomeScore       int64       	 `json:"homeScore"`
	AwayScore       int64       	 `json:"awayScore"`
	SeasonID        string         `json:"seasonID"`
	LeagueID        string         `json:"leagueID"`
	LeagueCanonical string         `json:"leagueCanonical"`
	URL             string         `json:"url"`
	Token           string         `json:"token"`
}

type Player struct {
	ID           string         `json:"id"`
	First        sql.NullString `json:"first"`
	Middle       sql.NullString `json:"middle"`
	Last         string         `json:"last"`
	FullName     string         `json:"fullName"`
	ShortName    string         `json:"shortName"`
	Canonical    string         `json:"canonical"`
	Height       float32        `json:"height"`
	Weight       float32        `json:"weight"`
	Hand         int            `json:"hand"`
	Birth        sql.NullString `json:"birth"`
	PositionID   string         `json:"positionID"`
	LeagueID     string         `json:"leagueID"`
	TeamID       string         `json:"teamID"`
	TeamName     string         `json:"teamName"`
	PlayerNumber sql.NullString `json:"playerNumber"`
	URL          string         `json:"url"`
}

type Score struct {
	Away map[string]string `json:"away"`
	Home map[string]string `json:"home"`
}

type JSONResponse struct {
	Msg string `json:"msg"`
}
