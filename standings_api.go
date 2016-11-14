package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
)

type ByGB []map[string]string

func (s ByGB) Len() int {
	return len(s)
} // len

func (s ByGB) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
} // Swap

func (s ByGB) Less(i, j int) bool {
	return s[i]["GB"] < s[j]["GB"]
} // Less

func getTopRecord(record []map[string]string) (int64, int64) {

	var topScore, topW, topL int64 = 0, 0, 0

	for _, v := range record {

		w, okw := v["W"]
		l, okl := v["L"]

		if !okw || !okl {
			continue
		}

		wins, err := strconv.ParseInt(w, 10, 32)

		if err != nil {
			continue
		}

		losses, err2 := strconv.ParseInt(l, 10, 32)

		if err2 != nil {
			continue
		}

		score := wins - losses

		if score > topScore {

			topScore = score
			topW = wins
			topL = losses

		}

	}

	return topW, topL

} // getTopRecord

func calcGB(record []map[string]string) {

	topW, topL := getTopRecord(record)

	for i, v := range record {

		w, okw := v["W"]
		l, okl := v["L"]

		if !okw || !okl {

			w = "0"
			l = "0"
			record[i]["W"] = "0"
			record[i]["L"] = "0"

		}

		wins, err := strconv.ParseInt(w, 10, 32)

		if err != nil {
			return
		}

		losses, err2 := strconv.ParseInt(l, 10, 32)

		if err2 != nil {
			return
		}

		record[i]["GB"] = fmt.Sprintf("%.1f", ((float64(topW-wins) + float64(losses-topL)) / 2.0))

	}

} // calcGB

func calcWinPCT(record map[string]string) string {

	var iw, il int64 = 0, 0
	var err1, err2 error

	w, okw := record["W"]
	l, okl := record["L"]

	if okw {

		iw, err1 = strconv.ParseInt(w, 10, 64)

		if err1 != nil {
			log.Fatal(err1)
		}

	}

	if okl {

		il, err2 = strconv.ParseInt(l, 10, 64)

		if err2 != nil {
			log.Fatal(err2)
		}

	}

	total := iw + il

	if total == 0 {
		return "0.000"
	} else {
		fpct := float64(iw) / float64(total)
		//return strconv.FormatFloat(fpct, 'f', 3, 64)
		return fmt.Sprintf("%.3f", fpct)
	}

} // calcWinPCT

func getStandings(req *Req) {

} // getStandings

func standingsAPIHandler(w http.ResponseWriter, r *http.Request) {

  vars := mux.Vars(r)
	league := vars["league"]
	//season := vars["season"]

	switch r.Method {
	case http.MethodGet:

		l := getLeague(league)

		s := getLatestSeasonByID(league)

		teams := getTeams(l)

		//standings := map[string]map[string]string{}
		standings := []map[string]string{}

		for _, t := range teams {

			key := fmt.Sprintf("record.%s.%s.%s", l.ID, s.ID, t.ID)

			res, err := redis.StringMap(config.Redis.Do("HGETALL", key))

			if err != nil {
				log.Println(err)
			}

			log.Println("fuck")
			log.Println(res)

			pct := calcWinPCT(res)

			res["PCT"] = pct
			res["NAME"] = t.Name
			res["ID"] = t.ID

			standings = append(standings, res)

		}

		calcGB(standings)

		sort.Sort(ByGB(standings))

		j, _ := json.Marshal(standings)
		w.Write(j)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

} // standingsAPIHandler
