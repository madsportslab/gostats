# gostats
Go backend RESTful API and websocket interface for madsportslab.

## Environment

### Development

(Runs on single server for now)
1.  Ubuntu 14.04.x server
  * sudo apt-get install git gcc make libssl-dev libreadline-dev zlib1g-dev g++ libsqlite3-dev sqlite3 libpcre3-dev
1.  golang 1.6.x
1.  redis 3.0.x
1.  sqlite3


### Production

(Runs on single server for now)
1.  Ubuntu 14.04.x server
  * sudo apt-get install git gcc make libssl-dev libreadline-dev zlib1g-dev g++ libsqlite3-dev sqlite3 libpcre3-dev
1.  golang 1.6.x
1.  redis 3.0.x
1.  postgresql
1.  nginx 1.8.x

### Setup

1.  clone code from github.com/madsportslab/gostats into $GOPATH/src/github.com/madsportslab
1.  install dependencies
  * go get github.com/eknkc/amber
  * go get github.com/garyburd/redigo/redis
  * go get github.com/gorilla/mux
  * go get github.com/gorilla/websocket
  * go get github.com/mattn/go-sqlite3
  * go get github.com/mattes/migrate 
1.  go build (from $GOPATH/src/github.com/madsportslab/gostats)
1.  migrate -url sqlite3://db/meta.db -path ./db/migrations up (from $GOPATH/src/github.com/madsportslab/gostats)
1.  ./gostats -redis-addr=X.X.X.X:6379


## Description

gostats connects to redis where all frequently changed game information is stored.  There
are 2 general use cases for gostats:

1.  Game Tracker for updating game logs
1.  Real Time statistics scores/stats

Data is updated in real time on clients and there could be multiple clients
accessing the socket.

## API

1.  connect to /ws/game/{gametoken}
1.  initialize game by sending an INIT cmd
```
{
  "cmd": "INIT",
  "leagueId": "38",
  "seasonId": "38",
  "gameId": "11",
  "homeId": "15",
  "awayId": "16",
  "periods": "4",
  "timestamp": ""
}
```

# Web Socket Interface

command | description | parameters | return value
--- | --- | --- | ---
INIT | Initializes game (scores), idempotent | Req structure |
PLAYCREATE | Makes a play | Req structure |
PLAYREMOVE  | Removes a play | Req structure |
PLAYCHANGE | Changes a play | Req structure |
PLAYS | Retrieves all plays | Req structure |
SCORE | Retrieves current score | Req structure |
FINAL | Ends game, updates record | Req structure |

# PROTOCOL

