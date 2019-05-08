package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Coord struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Snake struct {
	Kind          string          `json:"kind"`
	Name          string          `json:"name"`
	Color         string          `json:"color"`
	State         string          `json:"state"`
	Coords        []Coord         `json:"body"`
	WS            *websocket.Conn `json:"-"`
	Used          bool            `json:"-"`
	Direction     string          `json:"-"`
	LastDirection string          `json:"-"`
}

type Update struct {
	Kind   string  `json:"kind"`
	Apples []Coord `json:"apples"`
	Snakes []Snake `json:"snakes"`
}

type Init struct {
	Kind        string `json:"kind"`
	PlayersSlot []int  `json:"players_slot"`
	StateGame   string `json:"state_game"`
	MapSize     int    `json:"map_size"`
}

type Restart struct {
	Kind string `json:"kind"`
}

type KindOnly struct {
	Kind string `json:"kind"`
}

type Winner struct {
	Kind   string `json:"kind"`
	Player string `json:"player"`
}

var GeneralMutex sync.Mutex

var RestartGame = Restart{
	Kind: "restart",
}

var StateGame = Init{
	Kind:        "init",
	StateGame:   "waiting",
	MapSize:     50,
	PlayersSlot: []int{1, 2, 3, 4},
}

var socketList []*websocket.Conn

var snakeList []*Snake

var applesList []Coord

var Player1 = Snake{
	Kind: "snake",
	Coords: []Coord{
		Coord{X: 1, Y: 3},
		Coord{X: 1, Y: 2},
		Coord{X: 1, Y: 1},
	},
	Direction: "down",
}

var Player2 = Snake{
	Kind: "snake",
	Coords: []Coord{
		Coord{X: 48, Y: 3},
		Coord{X: 48, Y: 2},
		Coord{X: 48, Y: 1},
	},
	Direction: "down",
}

var Player3 = Snake{
	Kind: "snake",
	Coords: []Coord{
		Coord{X: 1, Y: 46},
		Coord{X: 1, Y: 47},
		Coord{X: 1, Y: 48},
	},
	Direction: "up",
}

var Player4 = Snake{
	Kind: "snake",
	Coords: []Coord{
		Coord{X: 48, Y: 46},
		Coord{X: 48, Y: 47},
		Coord{X: 48, Y: 48},
	},
	Direction: "up",
}

var PlayerWinner = Winner{
	Kind: "won",
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	http.Handle("/", websocket.Handler(HandleClient))
	err := http.ListenAndServe(":8082", nil) // starts an HTTP server with a given address and handler
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}

}

func HandleClient(ws *websocket.Conn) {

	ws.Write(getInitMessage())
	ws.Write(getUpdateMessage())
	var checkSocket bool
	for _, socket := range socketList {
		if ws == socket {
			checkSocket = true
		}
	}
	if checkSocket == false {
		socketList = append(socketList, ws)
		ws.Write(getInitMessage())
		ws.Write(getUpdateMessage())
	}

	var info map[string]interface{}
	var data string
	for {
		err := websocket.Message.Receive(ws, &data)
		json.Unmarshal([]byte(data), &info)

		/*	if info["kind"] == "restart" {
			fmt.Println("in restart")
			StateGame = Init{
				Kind:        "init",
				StateGame:   "waiting",
				MapSize:     50,
				PlayersSlot: []int{1, 2, 3, 4},
			}

			Player1 = Snake{
				Kind: "snake",
				Coords: []Coord{
					Coord{X: 1, Y: 3},
					Coord{X: 1, Y: 2},
					Coord{X: 1, Y: 1},
				},
				Direction: "down",
			}

			Player2 = Snake{
				Kind: "snake",
				Coords: []Coord{
					Coord{X: 48, Y: 3},
					Coord{X: 48, Y: 2},
					Coord{X: 48, Y: 1},
				},
				Direction: "down",
			}

			Player3 = Snake{
				Kind: "snake",
				Coords: []Coord{
					Coord{X: 1, Y: 46},
					Coord{X: 1, Y: 47},
					Coord{X: 1, Y: 48},
				},
				Direction: "up",
			}

			Player4 = Snake{
				Kind: "snake",
				Coords: []Coord{
					Coord{X: 48, Y: 46},
					Coord{X: 48, Y: 47},
					Coord{X: 48, Y: 48},
				},
				Direction: "up",
			}

			for _, webs := range socketList {
				webs.Write(getRestartMessage())
			}

		} */

		if info["kind"] == "start" {
			var aliveSnake int
			aliveSnake = 0
			for _, snake := range snakeList {
				if snake.State == "alive" {
					aliveSnake++
				}
			}
			if aliveSnake > 1 {
				StateGame.StateGame = "playing"
				for _, webs := range socketList {
					webs.Write(getInitMessage())
				}
				go automaticMove()
			}

		}

		if info["kind"] == "move" {
			var current *Snake
			current = checkCurrentPlayer(ws)
			if current != nil {
				moveSnake(info, current)
			}
		}

		if info["kind"] == "connect" {
			if StateGame.StateGame == "waiting" {
				if info["slot"] != 0 {
					var slot string
					slot = strconv.FormatFloat(info["slot"].(float64), 'f', 0, 64)
					if slot == "1" {
						Player1.WS = ws
						Player1.State = "alive"
						Player1.Used = true
						Player1.Name = info["name"].(string)
						Player1.Color = info["color"].(string)
						snakeList = append(snakeList, &Player1)
						majSlot(info["slot"])
					} else if slot == "2" {
						Player2.WS = ws
						Player2.Used = true
						Player2.State = "alive"
						Player2.Name = info["name"].(string)
						Player2.Color = info["color"].(string)
						majSlot(info["slot"]) //Call to majSlot function who delete the Snake takken
						snakeList = append(snakeList, &Player2)
					} else if slot == "3" {
						Player3.WS = ws
						Player3.Used = true
						Player3.State = "alive"
						Player3.Name = info["name"].(string)
						Player3.Color = info["color"].(string)
						majSlot(info["slot"]) //Call to majSlot function who delete the Snake takken
						snakeList = append(snakeList, &Player3)
					} else if slot == "4" {
						Player4.WS = ws
						Player4.Used = true
						Player4.State = "alive"
						Player4.Name = info["name"].(string)
						Player4.Color = info["color"].(string)
						majSlot(info["slot"]) //Call to majSlot function who delete the Snake takken
						snakeList = append(snakeList, &Player4)
					}
				}

			}
			for _, webs := range socketList {
				webs.Write(getInitMessage())
				webs.Write(getUpdateMessage())
			}
		}

		if err != nil {
			for i, socket := range socketList {
				if ws == socket {
					socketList = append(socketList[:i], socketList[i+1:]...)
				}
			}
			return
		}

		for _, webs := range socketList {
			webs.Write(getUpdateMessage())
		}
	}
}

func automaticMove() {
	if StateGame.StateGame == "playing" {
		for len(applesList) < 2 {
			createApple()
		}
		NewCoord := []Coord{}
		if len(snakeList) != 0 {
			var a []Coord
			for _, snake := range snakeList {

				if snake.State == "alive" {

					Direction := snake.Direction
					if Direction == "right" {
						snake.LastDirection = "right"
						checkBorder(Direction, snake)
						NewCoord = []Coord{Coord{snake.Coords[0].X + 1, snake.Coords[0].Y}}
						a = append(NewCoord, snake.Coords...)
						test := checkApple(a)
						if test == false {
							snake.Coords = a[:len(a)-1]
						} else {
							snake.Coords = a
						}
					} else if Direction == "left" {
						snake.LastDirection = "left"
						checkBorder(Direction, snake)
						NewCoord = []Coord{Coord{snake.Coords[0].X - 1, snake.Coords[0].Y}}
						a = append(NewCoord, snake.Coords...)
						test := checkApple(a)
						if test == false {
							snake.Coords = a[:len(a)-1]
						} else {
							snake.Coords = a
						}
					} else if Direction == "down" {
						snake.LastDirection = "down"
						checkBorder(Direction, snake)
						NewCoord = []Coord{Coord{snake.Coords[0].X, snake.Coords[0].Y + 1}}
						a = append(NewCoord, snake.Coords...)
						test := checkApple(a)
						if test == false {
							snake.Coords = a[:len(a)-1]
						} else {
							snake.Coords = a
						}
					} else if Direction == "up" {
						snake.LastDirection = "up"
						checkBorder(Direction, snake)
						NewCoord = []Coord{Coord{snake.Coords[0].X, snake.Coords[0].Y - 1}}
						a = append(NewCoord, snake.Coords...)
						test := checkApple(a)
						if test == false {
							snake.Coords = a[:len(a)-1]
						} else {
							snake.Coords = a
						}
					}
					checkCollision(snake)
				}
				if snake.State == "dead" {
					snake.Coords = []Coord{}
				}
			}
			for _, webs := range socketList {
				webs.Write(getUpdateMessage())
			}
			time.Sleep(100 * time.Millisecond) // after 150 Ms, call again the func automaticMove
			automaticMove()                    // use of goroutines so event the player who
		}
	}
}

func checkCollision(currentSnake *Snake) {
	for _, snake := range snakeList {
		if snake.WS != currentSnake.WS {
			for _, coord := range snake.Coords {
				if coord == currentSnake.Coords[0] {
					currentSnake.State = "dead"
					checkEnded()
				}
			}
		} else {
			for _, coord := range snake.Coords[1:] {
				if coord == snake.Coords[0] {
					currentSnake.State = "dead"
					checkEnded()
				}
			}
		}
	}
}

func checkApple(coords []Coord) bool {
	for i, apple := range applesList {
		if apple == coords[0] {
			applesList = append(applesList[:i], applesList[i+1:]...)
			return true
		}
	}
	return false
}

func checkBorder(keyMove string, snake *Snake) {
	if keyMove == "right" && snake.Coords[0].X+1 >= StateGame.MapSize {
		snake.State = "dead"
		checkEnded()
	}
	if keyMove == "left" && snake.Coords[0].X-1 < 0 {
		snake.State = "dead"
		checkEnded()
	}
	if keyMove == "up" && snake.Coords[0].Y-1 < 0 {
		snake.State = "dead"
		checkEnded()
	}
	if keyMove == "down" && snake.Coords[0].Y+1 >= StateGame.MapSize {
		snake.State = "dead"
		checkEnded()
	}
}
func checkEnded() {
	var aliveSnake []*Snake
	for _, snake := range snakeList {
		if snake.State == "alive" {
			aliveSnake = append(aliveSnake, snake)
		}
	}
	if len(aliveSnake) == 1 {
		StateGame.StateGame = "ended"
		PlayerWinner.Player = aliveSnake[0].Name
		for _, webs := range socketList {
			message, _ := json.Marshal(PlayerWinner)
			webs.Write(message)
			webs.Write(getUpdateMessage())
			webs.Write(getInitMessage())
		}
	}
}

func moveSnake(info map[string]interface{}, CurrentPlayer *Snake) {
	Direction := info["key"]
	CurrentDirection := CurrentPlayer.LastDirection
	if Direction == "right" {
		if CurrentDirection != "left" {
			CurrentPlayer.Direction = "right"
		}
	} else if Direction == "left" {
		if CurrentDirection != "right" {
			CurrentPlayer.Direction = "left"
		}
	} else if Direction == "down" {
		if CurrentDirection != "up" {
			CurrentPlayer.Direction = "down"
		}
	} else if Direction == "up" {
		if CurrentDirection != "down" {
			CurrentPlayer.Direction = "up"
		}
	}

}

func checkCurrentPlayer(ws *websocket.Conn) *Snake {
	if ws == Player1.WS {
		return &Player1
	} else if ws == Player2.WS {
		return &Player2
	} else if ws == Player3.WS {
		return &Player3
	} else if ws == Player4.WS {
		return &Player4
	}

	return nil
}

func majSlot(slot interface{}) {
	for i, number := range StateGame.PlayersSlot {
		if float64(number) == slot.(float64) {
			StateGame.PlayersSlot = append(StateGame.PlayersSlot[:i], StateGame.PlayersSlot[i+1:]...)
		}
	}
}

func createApple() {
	var newApple Coord
	newApple = Coord{X: randInt(0, StateGame.MapSize), Y: randInt(0, StateGame.MapSize)}
	for _, snake := range snakeList {
		for _, coord := range snake.Coords {
			if coord == newApple {
				createApple()
			} else {
				applesList = append(applesList, newApple)
				return
			}
		}
	}
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func getUpdateMessage() []byte {
	var m Update

	m.Kind = "update"
	m.Snakes = []Snake{Player1, Player2, Player3, Player4}
	m.Apples = applesList

	message, err := json.Marshal(m) // Transformation de l'objet "Update" en JSON
	if err != nil {
		fmt.Println("Something wrong with JSON Marshal map")
	}
	return message // (Json)
}

// "init" dans le protocole
func getInitMessage() []byte {
	// Transformation de l'objet "Init" en JSON
	message, err := json.Marshal(StateGame)
	if err != nil {
		fmt.Println("Something wrong with JSON Marshal init")
	}
	return message
}

/*
func getRestartMessage() []byte {
	message, err := json.Marshal(RestartGame)
	if err != nil {
		fmt.Println("Something wrong with JSON Marshal restart")
	}
	return message
} */
