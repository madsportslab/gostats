package main

import (
	"database/sql"
)

const (
	Madsportslab = "madsportslab"
	ShortForm    = "2006-01-02 15:04:05 -0700"
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
	DefaultLeague int         	 `json:"defaultLeague"`
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
	IsFollower bool         `json:"isFollower"`
}

type Team struct {
	ID              string         `json:"id"`
	Name            string         `json:"name"`
	Canonical       string         `json:"canonical"`
	Icon            sql.NullString `json:"icon"`
	LeagueID        string         `json:"leagueId"`
	LeagueCanonical string         `json:"leagueCanonical"`
	LeagueName      string         `json:"leagueName"`
	URL             string         `json:"url"`
}

type LeagueAdmin struct {
	ID 			 string `json:"id"`
	Email    string `json:"email"`
	LeagueID string `json:"leagueId"`
	UserID   string `json:"userId"`
}

type TeamAdmin struct {
	ID        string `json:"id"`
	LeagueID  string `json:"leagueId"`
	TeamID    string `json:"teamId"`
	TeamName  string `json:"teamName"`
	UserID    string `json:"userId"`
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
	IsFollower bool       `json:"isFollower"`
}

type Season struct {
	ID       string         `json:"id"`
	Games    int            `json:"games"`
	Periods  string         `json:"periods"`
	Duration int            `json:"duration"`
	Start    sql.NullString `json:"start"`
	Finish   sql.NullString `json:"finish"`
	LeagueID string         `json:"leagueId"`
}

type Game struct {
	ID              string         `json:"id"`
	Scheduled       sql.NullString `json:"scheduled"`
	Scheduled2      sql.NullString `json:"scheduled2"`
	HomeID          string         `json:"homeId"`
	AwayID          string         `json:"awayId"`
	HomeName        string         `json:"homeName"`
	AwayName        string         `json:"awayName"`
	Opponent        string         `json:"opponent"`
	Completed       bool           `json:"completed"`
	HomeScore       int64       	 `json:"homeScore"`
	AwayScore       int64       	 `json:"awayScore"`
	SeasonID        string         `json:"seasonId"`
	LeagueID        string         `json:"leagueId"`
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
	PositionID   string         `json:"positionId"`
	LeagueID     string         `json:"leagueId"`
	TeamID       string         `json:"teamId"`
	TeamName     string         `json:"teamName"`
	PlayerNumber sql.NullString `json:"playerNumber"`
	URL          string         `json:"url"`
}

type Leaders struct {
  LeagueID			string			`json:"leagueId"`
	SeasonID      string			`json:"seasonId"`
	Players       []PlayerAverage	`json:"players"`
}

type PlayerAverage struct {
  TeamID				string			`json:"teamId"`
	TeamName			string			`json:"teamName"`
	PlayerID			string			`json:"playerId"`
	PlayerName		string			`json:"playerName"`
	Points       	string      `json:"points"`
	Rebounds      string      `json:"rebounds"`
	Assists       string      `json:"assists"`
	Steals       	string      `json:"steals"`
	Blocks       	string      `json:"blocks"`
	Turnovers    	string      `json:"turnovers"`
	Fouls      	  string      `json:"fouls"`
	FTPct       	string      `json:"ftPCT"`
	FG2Pct       	string      `json:"fg2PCT"`
	FG3Pct       	string      `json:"fg3PCT"`
	GP            string      `json:"gp"`
}

type PlayerStat struct {
  Me        	*Player
	Seasons     []PlayerAverage
	Career      *PlayerAverage
}

type Score struct {
	Away map[string]string `json:"away"`
	Home map[string]string `json:"home"`
}

type Forgot struct {
	ID				string
  UserID		string
	Token			string
}

type JSONResponse struct {
	Msg string `json:"msg"`
}
