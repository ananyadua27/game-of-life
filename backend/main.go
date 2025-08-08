package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	rows    = 200
	cols    = 200
	tickDur = 100 * time.Millisecond
)

type Message struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data,omitempty"`
}

type CellToggle struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type SaveRequest struct {
	Name string `json:"name"`
}

type LoadRequest struct {
	Name string `json:"name"`
}

var (
	grid      = newGrid()
	clients   = make(map[*websocket.Conn]bool)
	gridMutex sync.Mutex
	run       = false
	upgrader  = websocket.Upgrader{}

	savedPatterns = make(map[string][][]int)
)

func newGrid() [][]int {
	g := make([][]int, rows)
	for i := range g {
		g[i] = make([]int, cols)
	}
	return g
}

func countNeighbors(g [][]int, x, y int) int {
	sum := 0
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			nx := (x + dx + cols) % cols
			ny := (y + dy + rows) % rows
			sum += g[ny][nx]
		}
	}
	return sum
}

func nextGen() [][]int {
	next := newGrid()
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			neighbors := countNeighbors(grid, x, y)
			if grid[y][x] == 1 && (neighbors == 2 || neighbors == 3) {
				next[y][x] = 1
			} else if grid[y][x] == 0 && neighbors == 3 {
				next[y][x] = 1
			}
		}
	}
	return next
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	clients[conn] = true

	// Send initial state
	sendGrid(conn)

	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			delete(clients, conn)
			return
		}

		var msg Message
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			continue
		}

		gridMutex.Lock()
		switch msg.Type {
		case "toggle":
			var toggle CellToggle
			json.Unmarshal(msg.Data, &toggle)
			grid[toggle.Y][toggle.X] = 1 - grid[toggle.Y][toggle.X]
		case "start":
			run = true
		case "stop":
			run = false
		case "clear":
			grid = newGrid()
		case "random":
			for y := range grid {
				for x := range grid[y] {
					if rand.Float64() < 0.3 {
						grid[y][x] = 1
					} else {
						grid[y][x] = 0
					}
				}
			}
		case "save":
			var req SaveRequest
			json.Unmarshal(msg.Data, &req)
			saved := newGrid()
			for y := range grid {
				copy(saved[y], grid[y])
			}
			savedPatterns[req.Name] = saved
		case "load":
			var req LoadRequest
			json.Unmarshal(msg.Data, &req)
			if pattern, ok := savedPatterns[req.Name]; ok {
				for y := range pattern {
					copy(grid[y], pattern[y])
				}
			}
		}
		gridMutex.Unlock()

		broadcastGrid()
	}
}

func sendGrid(conn *websocket.Conn) {
	gridMutex.Lock()
	defer gridMutex.Unlock()
	data, _ := json.Marshal(grid)
	msg := Message{Type: "grid", Data: data}
	conn.WriteJSON(msg)
}

func broadcastGrid() {
	gridMutex.Lock()
	defer gridMutex.Unlock()
	data, _ := json.Marshal(grid)
	msg := Message{Type: "grid", Data: data}

	for conn := range clients {
		if err := conn.WriteJSON(msg); err != nil {
			conn.Close()
			delete(clients, conn)
		}
	}
}

func tickLoop() {
	for {
		time.Sleep(tickDur)
		gridMutex.Lock()
		if run {
			grid = nextGen()
			gridMutex.Unlock()
			broadcastGrid()
		} else {
			gridMutex.Unlock()
		}
	}
}

func main() {
	fs := http.FileServer(http.Dir("../frontend"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", handleWS)

	go tickLoop()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
