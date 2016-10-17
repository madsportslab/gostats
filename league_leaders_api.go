package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
  "strconv"

  "github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
)

func calcStatAvg(stat string, gp string) float64 {

  res1, err1 := strconv.ParseInt(stat, 10, 64)

	if err1 != nil {
		log.Println(err1)
    return 0
	}

  res2, err2 := strconv.ParseInt(gp, 10, 64)

	if err2 != nil {
		log.Println(err2)
    return 0
	}

	if res2 != 0 {
    return float64(res1)/float64(res2)
  } else {
    return 0
  }

} // calcStatAvg

func calcPtsAvg(ftm string, fg2m string, fg3m string, gp string) float64 {

  res1, err1 := strconv.ParseInt(ftm, 10, 64)

	if err1 != nil {
		log.Println(err1)
    return 0
	}

  res2, err2 := strconv.ParseInt(fg2m, 10, 64)

	if err2 != nil {
		log.Println(err2)
    return 0
	}

  res3, err3 := strconv.ParseInt(fg3m, 10, 64)

	if err3 != nil {
		log.Println(err3)
    return 0
	}

  res4, err4 := strconv.ParseInt(gp, 10, 64)

	if err4 != nil {
		log.Println(err4)
    return 0
	}

  if res4 != 0 {
    return float64(res1+res2*2+res3*3)/float64(res4)
  } else {
    return 0
  }

} // calcPtsAvg

func calcRebAvg(oreb string, dreb string, gp string) float64 {

  res1, err1 := strconv.ParseInt(oreb, 10, 64)

	if err1 != nil {
		log.Println(err1)
    return 0
	}

  res2, err2 := strconv.ParseInt(dreb, 10, 64)

	if err2 != nil {
		log.Println(err2)
    return 0
	}

  res3, err3 := strconv.ParseInt(gp, 10, 64)

	if err3 != nil {
		log.Println(err3)
    return 0
	}

  if res3 != 0 {
    return float64(res1+res2)/float64(res3)
  } else {
    return 0
  }

} // calcRebAvg

func calcPctAvg(attempt string, made string, gp string) float64 {

  res1, err1 := strconv.ParseInt(attempt, 10, 64)

	if err1 != nil {
		log.Println(err1)
    return 0
	}

  res2, err2 := strconv.ParseInt(made, 10, 64)

	if err2 != nil {
		log.Println(err2)
    return 0
	}

  res3, err3 := strconv.ParseInt(gp, 10, 64)

	if err3 != nil {
		log.Println(err3)
    return 0
	}

  if res3 != 0 && (res1+res2) != 0 {
    return float64(res2/(res1+res2))/float64(res3)
  } else {
    return 0
  }

} // calcPctAvg

func getLeaders(req *Req) []PlayerAverage {

	avgs := []PlayerAverage{}

  l := getLeague(req.LeagueId)

	players := getPlayers(l)

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
      FTPct: calcPctAvg(res["1PTA"], res["1PTM"], res["GP"]),
      FG2Pct: calcPctAvg(res["2PTA"], res["2PTM"], res["GP"]),
      FG3Pct: calcPctAvg(res["3PTA"], res["3PTM"], res["GP"]),
		}

    avgs = append(avgs, avg)

  }

  return avgs

} // getLeaders

func leagueLeaderAPIHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	league := vars["league"]

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

      players := getLeaders(&req)

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

} // leagueLeaderAPIHandler
