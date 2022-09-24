package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"tiktaktoe/pkg/tiktaktoe"

	"github.com/gorilla/websocket"
)

var done = make(chan struct{})

var addr = flag.String("addr", "localhost:3000", "http service address")
var roomId = flag.String("room", "", "room id")

func gameLoop(c *websocket.Conn) {
	gameStarted := false
	waitingMsgSended := false
	var game *tiktaktoe.Game
	var message tiktaktoe.Message
	for {
		err := c.ReadJSON(&message)
		if err != nil {
			log.Println("read:", err)
			return
		}
		if !gameStarted && !waitingMsgSended {
			waitingMsgSended = true
			fmt.Println("Waiting other player")
		}
		event := message.Event
		if event == "ERROR" || event == "ROOM_CONNECTION" {
			fmt.Println(message.Message)
		} else if event == "GAME_START" {
			fmt.Println("Game started")
			game = &tiktaktoe.Game{}
			gameStarted = true
			game.Draw()
		} else if event == "YOUR_TURN" {
			for {
				fmt.Println("Enter Row(0-2): ")
				var row int8
				fmt.Scanln(&row)
				fmt.Println("Enter Column(0-2): ")
				var column int8
				fmt.Scanln(&column)
				err := game.MakePlay(row, column)
				if err == nil {
					movement := tiktaktoe.Movement{Row: row, Column: column}
					c.WriteJSON(tiktaktoe.Message{Event: "PLAY", Message: "Movement", Movement: movement})
					game.Draw()
					break
				}
				fmt.Println(err)
			}
		} else if event == "OPPONENT_TURN" {
			fmt.Println("Opponent turn")
		} else if event == "MOVEMENT" {
			game.MakePlay(message.Movement.Row, message.Movement.Column)
			game.Draw()
		} else if event == "WIN" || event == "LOST" || event == "DRAW" {
			fmt.Println(message.Message)
		}
	}
}

func main() {
	flag.Parse()
	id := *roomId
	var u url.URL
	if id == "" {
		u = url.URL{Scheme: "ws", Host: *addr, Path: "/ws/create-room/"}
	} else {
		u = url.URL{Scheme: "ws", Host: *addr, Path: "/ws/join-room/" + id}
	}

	fmt.Printf("connecting to %s \n", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	gameLoop(c)
}
