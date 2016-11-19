package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

  "github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
)

func getTeamLeaders(req *Req, team string) []PlayerAverage {

	avgs := []PlayerAverage{}

  l := getLeague(req.LeagueId)

	players := getTeamPlayers(l.ID, team)

  for _, p := range players {

    key := fmt.Sprintf("season.player.%s.%s.%s", req.LeagueId, req.SeasonId,
		  p.ID)
    
    res, err := redis.StringMap(config.Redis.Do("HGETALL", key))

		if err != nil {
			log.Println(err)
			continue
		}

		avg := PlayerAverage{
			TeamID: p.TeamID,
			TeamName: p.TeamName,
			PlayerID: p.ID,
			PlayerName: p.ShortName,
			Points: calcPtsAvg(res["1PTM"], res["2PTM"], res["3PTM"], res["GP"]),
      Rebounds: calcRebAvg(res["OREB"], res["DREB"], res["GP"]),
      Assists: calcStatAvg(res["AST"], res["GP"]),
      Steals: calcStatAvg(res["ST"], res["GP"]),
      Blocks: calcStatAvg(res["BS"], res["GP"]),
      Turnovers: calcStatAvg(res["TO"], res["GP"]),
      Fouls: calcStatAvg(res["PF"], res["GP"]),
      FTPct: calcPctAvg(res["1PTA"], res["1PTM"]),
      FG2Pct: calcPctAvg(res["2PTA"], res["2PTM"]),
      FG3Pct: calcPctAvg(res["3PTA"], res["3PTM"]),
		}

    avgs = append(avgs, avg)

  }

  return avgs

} // getTeamLeaders

func teamLeaderAPIHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	league := vars["league"]
  team   := vars["team"]

	//u := authenticate(r)

	switch r.Method {
	case http.MethodGet:

    if league == "" {
      w.WriteHeader(http.StatusBadRequest)
    } else {

      l := getLeague(league)

      s := getLatestSeasonByID(l.ID)

      leaders := Leaders{
        LeagueID: l.ID,
        SeasonID: s.ID,
        Players: []PlayerAverage{},
      }

      req := Req{
        LeagueId: l.ID,
        SeasonId: s.ID,
      }

      players := getTeamLeaders(&req, team)

      leaders.Players = players

      j, _ := json.Marshal(leaders)
			w.Write(j)

    }

	case http.MethodPut:
  case http.MethodPost:
	case http.MethodDelete:
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

} // teamLeaderAPIHandler
