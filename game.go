package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

const (
	PLAYCREATE = "PLAYCREATE"
	PLAYREMOVE = "PLAYREMOVE"
	PLAYCHANGE = "PLAYCHANGE"
	PLAYS      = "PLAYS"
	SCORE      = "SCORE"
	ALLSCORES  = "ALLSCORES"
	MYSCORES   = "MYSCORES"
	INIT       = "INIT"
	FINAL      = "FINAL"
	STATS      = "STATS"
	STATUPDATE = "STATUPDATE"
)

type PlayRemoveRes struct {
	ResponseType string
	Timestamp    string
	PlayerId     string
}

type PlayCreateRes struct {
	ResponseType string
	Timestamp    string
	Detail       string
}

type AllPlaysRes struct {
	ResponseType string
	All          map[string]string
}

type GenericSignal struct {
	W       *websocket.Conn
	GameKey string
	Buffer  []byte
}

type ScoreboardRes struct {
	ResponseType string
	Away         map[string]string
	Home         map[string]string
	AwayName     string
	HomeName     string
}

type GameScore struct {
	HomeName  string `json:"homeName"`
	AwayName  string `json:"awayName"`
	HomeScore string `json:"homeScore"`
	AwayScore string `json:"awayScore"`
	Period    string `json:"period"`
	LeagueID  string `json:"leagueId"`
	SeasonID  string `json:"seasonId"`
	GameID    string `json:"gameId"`
}

type ScoreRes struct {
	ResponseType string
	Games        map[string]GameScore
}

type MyScoreRes struct {
	ResponseType string
	Games        map[string]map[string]GameScore
}

type StatRes struct {
	ResponseType string
	Away         [][]string
	Home         [][]string
	AwayName     string
	HomeName     string
}

type Stat struct {
	PlayerId string `json:"playerId"`
	TeamId   string `json:"teamId"`
	Play     string `json:"play"`
	Period   string `json:"period"`
	Clock    string `json:"clock"`
	Team     string `json:"team"`
	Name     string `json:"name"`
}

type Req struct {
	Cmd       string   `json:"cmd"`
	Data      *Stat    `json:"data"`
	LeagueId  string   `json:"leagueId"`
	SeasonId  string   `json:"seasonId"`
	GameId    string   `json:"gameId"`
	HomeId    string   `json:"homeId"`
	AwayId    string   `json:"awayId"`
	Periods   string   `json:"periods"`
	Timestamp string   `json:"timestamp"`
	GameKey   string   `json:"gamekey"`
	UserToken string   `json:"userToken"`
	Games     []string `json:"games"`
}

type Broadcast struct {
	ResponseType string
	ID           string
	Buffer       chan []byte
}

type Connection struct {
	W         *websocket.Conn
	Broadcast *Broadcast
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var broadcast = make(chan GenericSignal)
var unicast = make(chan GenericSignal)

func writer() {

	for {
		select {
		case s := <-broadcast:

			for cnx := range m.Games[s.GameKey].Listeners {

				err := m.Games[s.GameKey].Listeners[cnx].W.WriteMessage(
					websocket.TextMessage, []byte(s.Buffer))

				if err != nil {
					log.Println(err)
				}

			}

		case s := <-unicast:

			err := s.W.WriteMessage(websocket.TextMessage, []byte(s.Buffer))

			if err != nil {
				log.Println(err)
			}

		}
	}

} // writer

func gameHandler(w http.ResponseWriter, r *http.Request) {

	c, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	go m.Run()

	m.Subscribe <- &Message{
		ID: id,
		W:  c,
	}

	defer c.Close()

	for {

		_, msg, err := c.ReadMessage()

		if err != nil || msg == nil {
			return
		}

		req := Req{}

		json.Unmarshal(msg, &req)

		switch req.Cmd {
		case INIT:
			initScore(&req)

		case PLAYREMOVE:
			j := playRemove(&req)

			gs := GenericSignal{
				GameKey: id,
				Buffer:  j,
			}

			broadcast <- gs

			j2 := update(&req)

			gs2 := GenericSignal{
				GameKey: id,
				Buffer:  j2,
			}

			broadcast <- gs2

		case PLAYCREATE:
			t := playCreate(&req)

			j := getPlay(&req, t)

			gs := GenericSignal{
				GameKey: id,
				Buffer:  j,
			}

			broadcast <- gs

			j2 := update(&req)

			gs2 := GenericSignal{
				GameKey: id,
				Buffer:  j2,
			}

			broadcast <- gs2

		case PLAYCHANGE:

		case PLAYS:
			j := getAllPlays(&req)

			gs := GenericSignal{
				W:       c,
				GameKey: id,
				Buffer:  j,
			}

			unicast <- gs

		case SCORE:
			j := getScore(&req)

			gs := GenericSignal{
				GameKey: id,
				Buffer:  j,
			}

			broadcast <- gs

		case ALLSCORES:

			j := getScores(id)

			gs := GenericSignal{
				GameKey: id,
				Buffer:  j,
			}

			broadcast <- gs

		case MYSCORES:

			u := tokenToUser(req.UserToken)

			j := getMyScores(u)

			gs := GenericSignal{
				GameKey: id,
				Buffer:  j,
			}

			broadcast <- gs

		case FINAL:
			gameFinal(&req)

		case STATS:

			j := getGameStats(&req)

			gs := GenericSignal{
				GameKey: id,
				Buffer:  j,
			}

			broadcast <- gs

		case STATUPDATE:
			//getStat(&req)
			log.Println("ph")

		default:
			log.Println(string(msg) + " is an invalid command")

		}

	}

} // gameHandler
