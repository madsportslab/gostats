package main

import (
	"crypto/cipher"
	"database/sql"
	"flag"
	"log"
	"net/http"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	Redis    redis.Conn
	Database *sql.DB
	HashKey  []byte
	BlockKey []byte
	IV       []byte
	Cipher   *cipher.Block
	Secret   string
}

var config = Config{}

var serverAddr = flag.String("server", ":9999", "Server address")
var databaseAddr = flag.String("database", "./db/meta.db", "Database address")
var redisAddr = flag.String("redis", ":6379", "Redis address")
var domain = flag.String("domain", "127.0.0.1", "Domain address")
var certFile = flag.String("cert", "ssl.crt", "SSL Certificate")
var keyFile = flag.String("key", "ssl.key", "SSL Key")
var useSSL = flag.Bool("use-ssl", false, "Use Encrypted Server (SSL)")

func initRedis() {

	r, err := redis.Dial("tcp", *redisAddr)

	if err != nil {
		log.Fatal("Redis connection error: ", err)
	}

	config.Redis = r

} // initRedis

func initDatabase() {

	db, err := sql.Open("sqlite3", *databaseAddr)

	if err != nil {
		log.Fatal("Database connection error: ", err)
	}

	config.Database = db

} // initDatabase

func initKeys() {

	//TODO: leave key
	config.HashKey = []byte("ABCDEF")
	config.BlockKey = []byte("1234567890123456")
	config.IV = []byte("Rubbish 16 long.")
	config.Cipher = nil
	config.Secret = "S3CR3T"

} // initKeys

func initRoutes() *mux.Router {

	router := mux.NewRouter()

	// api routes
	router.HandleFunc("/auth", authAPIHandler)

	// leagues
	router.HandleFunc("/api/leagues", leagueAPIHandler)
	router.HandleFunc("/api/leagues/{league:all}", leagueAPIHandler)
	router.HandleFunc("/api/leagues/{league:[0-9]+}", leagueAPIHandler)
	router.HandleFunc("/api/leagues/{league:[0-9]+}/players",
		playerAPIHandler)
	router.HandleFunc("/api/leagues/{league:[0-9]+}/games",
		scheduleAPIHandler)
	router.HandleFunc("/api/leagues/{league:[0-9]+}/games/{game:[0-9]+}",
		scheduleAPIHandler)
	router.HandleFunc("/api/leagues/{league:[0-9]+}/standings",
		standingsAPIHandler)
	router.HandleFunc("/api/leagues/{league:[0-9]+}/followers",
	  followerAPIHandler)
	router.HandleFunc("/api/leagues/{league:[0-9]+}/followers/{follower:[0-9]+}",
	  followerAPIHandler)
	router.HandleFunc("/api/leagues/{league:following}", leagueAPIHandler)

	// seasons
	router.HandleFunc(
		"/api/leagues/{league:[0-9]+}/seasons/{season:[0-9]+}",
		seasonAPIHandler)

	router.HandleFunc(
		"/api/leagues/{league:[0-9]+}/seasons/{season:[0-9]+}/standings",
		standingsAPIHandler)

	router.HandleFunc(
		"/api/leagues/{league:[0-9]+}/seasons/{season:[0-9]+}/games",
		scheduleAPIHandler)

	router.HandleFunc(
		"/api/leagues/{league:[0-9]+}/seasons/{season:[0-9]+}/games/{game:[0-9]+}",
		scheduleAPIHandler)

	// teams
	router.HandleFunc("/api/leagues/{league:[0-9]+}/teams",
		teamAPIHandler)

	router.HandleFunc("/api/leagues/{league:[0-9]+}/teams/{team:[0-9]+}",
		teamAPIHandler)

	router.HandleFunc("/api/leagues/{league:[0-9]+}/teams/{team:[0-9]+}/admins",
		teamAdminHandler)

	router.HandleFunc(
		"/api/leagues/{league:[0-9]+}/teams/{team:[0-9]+}/players",
		playerAPIHandler)

	router.HandleFunc(
		"/api/leagues/{league:[0-9]+}/teams/{team:[0-9]+}/players/{player:[0-9]+}",
		playerAPIHandler)

	router.HandleFunc(
		"/api/leagues/{league:[0-9]+}/teams/{team:[0-9]+}/games",
		scheduleAPIHandler)

	// users
	router.HandleFunc("/api/users", userAPIHandler)
	router.HandleFunc("/api/users/{user:[0-9]+}", userAPIHandler)

	// special
	router.HandleFunc("/api/games", specialAPIHandler)
	router.HandleFunc("/api/games/{id:[0-9]+}", specialAPIHandler)

	// betausers
	router.HandleFunc("/api/betausers", betaAPIHandler)

	// sports
	router.HandleFunc("/api/sports", sportAPIHandler)
	router.HandleFunc("/api/sports/{sport:[0-9]+}/positions", positionAPIHandler)

	// websocket routes
	router.HandleFunc("/ws/game/{id:[0-9a-f]+}", gameHandler)
	router.HandleFunc("/ws/game/me", gameHandler)
	router.HandleFunc("/ws/game/scorechannel/{id:[0-9]+}", gameHandler)

	//router.HandleFunc("/ws/clocks/{id:[0-9a-z]+}", clockHandler)

	return router

} // initRoutes

func main() {

	flag.Parse()

	initRedis()
	initDatabase()
	initKeys()

	defer config.Redis.Close()
	defer config.Database.Close()

	go writer()

	if *useSSL {

		log.Println("Go Stats started on port 9999 in secure mode...")
		log.Fatal(http.ListenAndServeTLS(*serverAddr, *certFile, *keyFile,
			initRoutes()))

	} else {

		log.Println("Go Stats started on port 9999 in insecure mode...")
		log.Fatal(http.ListenAndServe(*serverAddr, initRoutes()))

	}

} // main
