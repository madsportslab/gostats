package main

import (
  "log"

  "github.com/gorilla/websocket"
)

type Message struct {
  ID      string
  W       *websocket.Conn
}

type GameChannel struct {
  Controllers map[*websocket.Conn] *Connection
  Listeners   map[*websocket.Conn] *Connection
}

type Mgr struct {
  Games         map[string] *GameChannel
  Subscribe     chan *Message
  Unsubscribe   chan *Message
}

var m = Mgr{
  Games: make(map[string] *GameChannel),
  Subscribe: make(chan *Message),
  Unsubscribe: make(chan *Message),
}

func (m *Mgr) Run() {

  for {
    select {
      case c := <-m.Subscribe:

        _, ok  := m.Games[c.ID]
                
        if !ok {
          
          gc := GameChannel{
            Controllers:  make(map[*websocket.Conn] *Connection),
            Listeners:    make(map[*websocket.Conn] *Connection),
          }
                    
          m.Games[c.ID] = &gc
          
        }
          
        bcast := Broadcast{
          ID: c.ID,
          Buffer: make(chan []byte),
        }
        
        conn := Connection{
          W: c.W,
          Broadcast: &bcast,
        }
        
        m.Games[c.ID].Listeners[c.W] = &conn
        
      case c := <-m.Unsubscribe:
        log.Println("unsubscribe not yet implemented")
        log.Println(c)
          
    }
    
  } 

} // Run

