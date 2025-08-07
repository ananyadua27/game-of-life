# Conway’s Game of Life — Web-Based Cellular Automaton

An interactive, real-time implementation of Conway's Game of Life, built with Go for concurrency and performance, and WebSockets for low-latency updates. The frontend leverages the HTML5 Canvas API for rendering a scalable grid with rich UI interactions.

---

## Features

- **Live WebSocket Communication**: Real-time grid updates across multiple clients
- **Concurrent Game Loop**: Efficient tick-based simulation using Go’s goroutines
- **Pattern Persistence**: Save/load named patterns in memory for experimentation
- **Dynamic Grid Controls**: Toggle cells, randomize, clear, pause/play simulation
- **Responsive UI**: Canvas-based rendering with emoji-enhanced live cell display
- **Modular Design**: Easy to extend with additional features (e.g. pattern libraries, backend storage)

---

## Tech Stack

| Layer         | Technology                          |
|--------------|--------------------------------------|
| Backend       | **Go (Golang)** - concurrency, WebSocket server |
| Frontend      | **JavaScript (Vanilla)** - Canvas API rendering |
| Communication | **WebSockets** via Gorilla WebSocket |
| Web Server    | `net/http` in Go, static file serving |
| Rendering     | HTML5 `<canvas>` + Unicode emojis   |
| Architecture  | Event-driven with shared state and mutex locking |

---


