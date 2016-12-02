package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

const (
	W       = "W"
	L       = "L"
	L10W    = "L10W"
	L10L    = "L10L"
	HOMEW   = "HOMEW"
	HOMEL   = "HOMEL"
	AWAYW   = "AWAYW"
	AWAYL   = "AWAYL"
	STREAKW = "STREAKW"
	STREAKL = "STREAKL"
)

func checkerr(e error) {

	if e != nil {
		log.Println(e)
	}

} // checkerr

func keyExists(key string) bool {

	res, err := redis.Bool(config.Redis.Do("EXISTS", key))

	if err != nil {
		log.Println(err)
		return false
	}

	if res {
		return true
	} else {
		return false
	}

} // keyExists

// TODO: initTeamScore

func initScore(req *Req) {

	key1 := fmt.Sprintf("score.%s.%s.%s.%s", req.LeagueId, req.SeasonId,
		req.GameId, req.HomeId)

	key2 := fmt.Sprintf("score.%s.%s.%s.%s", req.LeagueId, req.SeasonId,
		req.GameId, req.AwayId)

	periods, err := strconv.ParseInt(req.Periods, 10, 32)

	if err != nil {
		periods = 2 // default number of periods
	}

	if !keyExists(key1) {

		if periods == 2 {

			_, err := config.Redis.Do("HMSET", key1, "1", 0, "2", 0, "T", 0)

			if err != nil {
				log.Println(err)
			}

		} else {

			_, err := config.Redis.Do("HMSET", key1, "1", 0, "2", 0, "3", 0, "4", 0,
				"T", 0)

			if err != nil {
				log.Println(err)
			}

		}

	}

	if !keyExists(key2) {

		if periods == 2 {

			_, err = config.Redis.Do("HMSET", key2, "1", 0, "2", 0, "T", 0)

			if err != nil {
				log.Println(err)
			}

		} else {

			_, err = config.Redis.Do("HMSET", key2, "1", 0, "2", 0, "3", 0, "4", 0,
				"T", 0)

			if err != nil {
				log.Println(err)
			}

		}

	}

} // initScore

func initRecord(req *Req, teamid string) {

	key := fmt.Sprintf("record.%s.%s.%s", req.LeagueId, req.SeasonId,
		teamid)

	if !keyExists(key) {

		_, err := config.Redis.Do("HMSET", key, "W", 0, "L", 0, "PCT", "0.000",
			"HOMEW", 0, "HOMEL", 0, "AWAYW", 0, "AWAYL", 0, "L10W", 0, "L10L", 0,
			"STREAKW", 0, "STREAKL", 0, "STREAK", 0)

		if err != nil {
			log.Println(err)
		}

	}

} // initRecord

func initPlayerStat(key string) {

	if !keyExists(key) {

		_, err := config.Redis.Do("HMSET", key, "1PTA", 0, "1PTM", 0, "2PTA", 0, "2PTM",
			0, "3PTA", 0, "3PTM", 0, "OREB", 0, "TREB", 0, "AST", 0, "ST", 0,
			"BS", 0, "TO", 0, "PF", 0)

		if err != nil {
			log.Println(err)
		}

	}

} // initPlayerStat

func getPlay(req *Req, t string) []byte {

	key := fmt.Sprintf("game.%s.%s.%s", req.LeagueId, req.SeasonId, req.GameId)

	res, err := redis.String(config.Redis.Do("HGET", key, t))

	if err != nil {
		log.Println(err)
	}

	pcr := PlayCreateRes{
		ResponseType: PLAYCREATE,
		Timestamp:    t,
		Detail:       res,
	}

	j, _ := json.Marshal(pcr)

	return j

} // getPlay

func playRemove(req *Req) []byte {

	key := fmt.Sprintf("game.%s.%s.%s", req.LeagueId, req.SeasonId, req.GameId)

	_, err := config.Redis.Do("HDEL", key, req.Timestamp)

	if err != nil {
		log.Println(err)
	}

	rp := PlayRemoveRes{
		ResponseType: PLAYREMOVE,
		Timestamp:    req.Timestamp,
		PlayerId:     req.Data.PlayerId,
	}

	j, _ := json.Marshal(rp)

	return j

} // playRemove

// wrapper function
func update(req *Req) []byte {

	playerUpdate(req)
	scoreUpdate(req)
	//seasonUpdate(req)
	//careerUpdate(req)

	j := getScore(req)

	return j

} // update

func playCreate(req *Req) string {

	key := fmt.Sprintf("game.%s.%s.%s", req.LeagueId, req.SeasonId, req.GameId)

	now := time.Now().Unix()

	j, _ := json.Marshal(req.Data)

	_, err := config.Redis.Do("HSETNX", key, now, j)

	if err != nil {
		log.Println(err)
	}

	return strconv.FormatInt(now, 10)

} // playCreate

func playerUpdate(req *Req) {

	key := fmt.Sprintf("player.%s.%s.%s.%s", req.LeagueId, req.SeasonId,
		req.GameId, req.Data.PlayerId)

	v := 1

	if req.Cmd == PLAYREMOVE {
		v = v * -1
	}

	res, err := config.Redis.Do("HINCRBY", key, req.Data.Play, v)
	log.Println(res)

	if err != nil {
		log.Println(err)
	}

} // playerUpdate

func seasonUpdate(req *Req) {

	key := fmt.Sprintf("season.player.%s.%s.%s", req.LeagueId, req.SeasonId,
		req.Data.PlayerId)

	v := 1

	if req.Cmd == PLAYREMOVE {
		v = v * -1
	}

	_, err := config.Redis.Do("HINCRBY", key, req.Data.Play, v)

	if err != nil {
		log.Println(err)
	}

} // seasonUpdate

func careerUpdate(req *Req) {

	key := fmt.Sprintf("career.player.%s.%s", req.LeagueId, req.Data.PlayerId)

	v := 1

	if req.Cmd == PLAYREMOVE {
		v = v * -1
	}

	res, err := config.Redis.Do("HINCRBY", key, req.Data.Play, v)

	log.Println(res)

	if err != nil {
		log.Println(err)
	}

} // careerUpdate

func scoreUpdate(req *Req) {

	key := fmt.Sprintf("score.%s.%s.%s.%s", req.LeagueId, req.SeasonId,
		req.GameId, req.Data.TeamId)

	v := 0

	switch req.Data.Play {
	case "1PTM":
		v = 1
	case "2PTM":
		v = 2
	case "3PTM":
		v = 3
	}

	if req.Cmd == PLAYREMOVE {
		v = v * -1
	}

	_, err1 := config.Redis.Do("HINCRBY", key, req.Data.Period, v)

	if err1 != nil {
		log.Println(err1)
	}

	_, err2 := config.Redis.Do("HINCRBY", key, "T", v)

	if err2 != nil {
		log.Println(err2)
	}

} // scoreUpdate

func getAllPlays(req *Req) []byte {

	key := fmt.Sprintf("game.%s.%s.%s", req.LeagueId, req.SeasonId, req.GameId)

	res, err := redis.StringMap(config.Redis.Do("HGETALL", key))

	if err != nil {
		log.Println(err)
	}

	ap := AllPlaysRes{
		ResponseType: "PLAYS",
		All:          res,
	}

	j, err := json.Marshal(ap)

	if err != nil {
		log.Println(err)
	}

	return j

} // getAllPlays

func getTeamScore(req *Req, teamid string) map[string]string {

	key := fmt.Sprintf("score.%s.%s.%s.%s", req.LeagueId, req.SeasonId,
		req.GameId, teamid)

	res, err := redis.StringMap(config.Redis.Do("HGETALL", key))

	if err != nil {
		log.Println(err)
	}

	return res

} // getTeamScore

func getTeamScoreEx(leagueid string, seasonid string, gameid string,
	teamid string) map[string]string {

	key := fmt.Sprintf("score.%s.%s.%s.%s", leagueid, seasonid,
		gameid, teamid)

	res, err := redis.StringMap(config.Redis.Do("HGETALL", key))

	if err != nil {
		log.Println(err)
	}

	return res

} // getTeamScoreEx

func getScore(req *Req) []byte {

	awayScore := getTeamScore(req, req.AwayId)
	homeScore := getTeamScore(req, req.HomeId)

	awayTeam := getTeam(req.AwayId)
	homeTeam := getTeam(req.HomeId)

	sb := ScoreboardRes{
		ResponseType: "SCORE",
		Away:         awayScore,
		Home:         homeScore,
		AwayName:     awayTeam.Name,
		HomeName:     homeTeam.Name,
	}

	j, err := json.Marshal(sb)

	if err != nil {
		log.Println(err)
	}

	return j

} // getScore

func getScores(leagueid string, gameDay string) []byte {

	games := getGameScores(leagueid, gameDay)

	sr := ScoreRes{
		ResponseType: ALLSCORES,
		Games:        games,
	}

	j, err := json.Marshal(sr)

	if err != nil {
		log.Println(err)
	}

	return j

} // getScores

func getGameScores(leagueid string, gameDay string) map[string]GameScore {

	games := getGames(leagueid, gameDay)
	
	x := map[string]GameScore{}

	for _, g := range games {

		awayScore := getTeamScoreEx(g.LeagueID, g.SeasonID, g.ID, g.AwayID)
		homeScore := getTeamScoreEx(g.LeagueID, g.SeasonID, g.ID, g.HomeID)

		hs, hok := homeScore["T"]
		as, aok := awayScore["T"]

		if !hok {
			hs = "0"
		}

		if !aok {
			as = "0"
		}

		gs := GameScore{
			HomeScore: hs,
			AwayScore: as,
			HomeName:  g.HomeName,
			AwayName:  g.AwayName,
			LeagueID:  g.LeagueID,
			SeasonID:  g.SeasonID,
			GameID:    g.ID,
		}

		x[g.ID] = gs

	}

	return x

} // getGameScores

func getMyScores(user *User, gameDay string) []byte {

	leagues := getAllMyLeagues(user)

	x := map[string]map[string]GameScore{}

	for _, league := range leagues {

		games := getGameScores(league.ID, gameDay)

		x[league.Name] = games

	}

	sr := MyScoreRes{
		ResponseType: MYSCORES,
		Games:        x,
	}

	j, err := json.Marshal(sr)

	if err != nil {
		log.Println(err)
	}

	return j

} // getMyScores

func getScoreTotal(req *Req, teamid string) int64 {

	key := fmt.Sprintf("score.%s.%s.%s.%s", req.LeagueId, req.SeasonId,
		req.GameId, teamid)

	res, err := redis.StringMap(config.Redis.Do("HGETALL", key))

	if err != nil {
		log.Println(err)
		return -1
	}

	total, err2 := strconv.ParseInt(res["T"], 10, 64)

	if err2 != nil {
		log.Println(err2)
		return -1
	}

	return total

} // getScoreTotal

func getTeamRecord(req *Req, teamid string) map[string]string {

	key := fmt.Sprintf("record.%s.%s.%s", req.LeagueId, req.SeasonId,
		teamid)

	record, err := redis.StringMap(config.Redis.Do("HGETALL",
		key))

	if err != nil {
		log.Println(err)
		return nil
	}

	return record

} // getTeamRecord

func addTeamWin(req *Req, teamid string) {

	record := getTeamRecord(req, teamid)

	var iL10W, iL10L int64 = 0, 0
	var err, err2 error = nil, nil

	if record[L10W] != "" {

		iL10W, err = strconv.ParseInt(record[L10W], 10, 32)

		if err != nil {
			log.Println(err)
			return
		}

	}

	if record[L10L] != "" {

		iL10L, err2 = strconv.ParseInt(record[L10L], 10, 32)

		if err2 != nil {
			log.Println(err)
			return
		}

	}

	key := fmt.Sprintf("record.%s.%s.%s", req.LeagueId, req.SeasonId,
		teamid)

	config.Redis.Send("HINCRBY", key, "W", 1)
	config.Redis.Send("HINCRBY", key, "STREAKW", 1)
	config.Redis.Send("HMSET", key, "STREAKL", 0)

	if req.HomeId == teamid {
		config.Redis.Send("HINCRBY", key, "HOMEW", 1)
	} else {
		config.Redis.Send("HINCRBY", key, "AWAYW", 1)
	}

	if iL10W < 10 {
		config.Redis.Send("HINCRBY", key, L10W, 1)
	}

	if iL10L > 0 {
		config.Redis.Send("HINCRBY", key, L10L, -1)
	}

	config.Redis.Flush()

} // addTeamWin

func addTeamLoss(req *Req, teamid string) {

	record := getTeamRecord(req, teamid)

	var iL10W, iL10L int64 = 0, 0
	var err, err2 error = nil, nil

	if record[L10W] != "" {

		iL10W, err = strconv.ParseInt(record[L10W], 10, 32)

		if err != nil {
			log.Println(err)
			return
		}

	}

	if record[L10L] != "" {

		iL10L, err2 = strconv.ParseInt(record[L10L], 10, 32)

		if err2 != nil {
			log.Println(err)
			return
		}

	}

	key := fmt.Sprintf("record.%s.%s.%s", req.LeagueId, req.SeasonId,
		teamid)

	config.Redis.Send("HINCRBY", key, "L", 1)
	config.Redis.Send("HINCRBY", key, "STREAKL", 1)
	config.Redis.Send("HMSET", key, "STREAKW", 0)

	if req.HomeId == teamid {
		config.Redis.Send("HINCRBY", key, "HOMEL", 1)
	} else {
		config.Redis.Send("HINCRBY", key, "AWAYL", 1)
	}

	if iL10L < 10 {
		config.Redis.Send("HINCRBY", key, L10L, 1)
	}

	if iL10W > 0 {
		config.Redis.Send("HINCRBY", key, L10W, -1)
	}

	config.Redis.Flush()

} // addTeamLoss

func gameFinal(req *Req) {

  // check for if game already completed
	g := getGame(req.GameId)

	if g.Completed {
		log.Println("game already completed")
	 	return
	}

	initRecord(req, req.AwayId)
	initRecord(req, req.HomeId)

	awayTotal := getScoreTotal(req, req.AwayId)
	homeTotal := getScoreTotal(req, req.HomeId)

	if awayTotal > homeTotal {

		addTeamWin(req, req.AwayId)
		addTeamLoss(req, req.HomeId)

	} else {

		addTeamWin(req, req.HomeId)
		addTeamLoss(req, req.AwayId)

	}

  updateGamePlayed(req, req.HomeId)
	updateGamePlayed(req, req.AwayId)

	updateLifetime(req, req.HomeId)
	updateLifetime(req, req.AwayId)

} // gameFinal

func updateGamePlayed(req *Req, teamid string) {

	// based on whether or not a statistic was tracked
	// in the future the seconds played will be tracked
	// which will be much more reliable

	players := getTeamPlayers(req.LeagueId, teamid)

	s := Stat{
		Play: "GP",
	}
	req.Data = &s

	for _, p := range players {

		key := fmt.Sprintf("player.%s.%s.%s.%s", req.LeagueId, req.SeasonId,
			req.GameId, p.ID)

		res, err := redis.StringMap(config.Redis.Do("HGETALL", key))

		if err != nil {
			log.Println("updateGamePlayed: ", err)
		}
		
		for _, v := range res {

			if v != "0" {

				req.Data.PlayerId = p.ID

				update(req)
				break

			}

		}	

	}

	//update(req)

} // updateGamePlayed

func getTeamStats(req *Req, teamid string) [][]string {

	players := getTeamPlayers(req.LeagueId, teamid)

	gs := [][]string{}

	for _, p := range players {

		key := fmt.Sprintf("player.%s.%s.%s.%s", req.LeagueId, req.SeasonId,
			req.GameId, p.ID)

		if !keyExists(key) {
			initPlayerStat(key)
		}

		res, err := redis.StringMap(config.Redis.Do("HGETALL", key))

		if err != nil {
			log.Println(err)
		}

		stats := []string{
			p.ShortName,
			totalPoints(res["1PTM"], res["2PTM"], res["3PTM"]),
			fieldGoals(res["2PTA"], res["2PTM"]),
			fieldGoals(res["3PTA"], res["3PTM"]),
			fieldGoals(res["1PTA"], res["1PTM"]),
			res["OREB"],
			rebounds(res["OREB"], res["DREB"]),
			res["AST"],
			res["ST"],
			res["BS"],
			res["TO"],
			res["PF"],
			p.PlayerNumber.String,
		}

		gs = append(gs, stats)

	}

	return gs

} // getTeamStats

func getGameStats(req *Req) []byte {

	homeStats := getTeamStats(req, req.HomeId)
	awayStats := getTeamStats(req, req.AwayId)

	sr := StatRes{
		ResponseType: "STATS",
		Home:         homeStats,
		Away:         awayStats,
	}

	j, err := json.Marshal(sr)

	if err != nil {
		log.Println(err)
	}

	return j

} // getGameStats 

func updateLifetime(req *Req, teamid string) {

	players := getTeamPlayers(req.LeagueId, teamid)

  for _, p := range players {

		playerKey := fmt.Sprintf("player.%s.%s.%s.%s", req.LeagueId, req.SeasonId,
		  req.GameId, p.ID)

    seasonKey := fmt.Sprintf("season.player.%s.%s.%s", req.LeagueId, req.SeasonId,
		  p.ID)
    
		careerKey := fmt.Sprintf("career.player.%s.%s", req.LeagueId, p.ID)
    
		accumulate(playerKey, seasonKey)
		accumulate(playerKey, careerKey)

  }

} // updateLifetime

func accumulate(src string, dest string) {

	res, err := redis.StringMap(config.Redis.Do("HGETALL", src))

	if err != nil {
		log.Println(err)
		return
	}

	for play := range res {

		val, err := strconv.ParseInt(res[play], 10, 64)

		if err != nil {
			log.Println(err)
			continue
		}

		_, err2 := config.Redis.Do("HINCRBY", dest, play, val)
		
		if err2 != nil {
			log.Println(err2)
			continue
		}

	}

} // accumulate
