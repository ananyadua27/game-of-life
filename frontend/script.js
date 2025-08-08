const canvas = document.getElementById('gameCanvas');
const ctx = canvas.getContext('2d');
const socket = new WebSocket('ws://localhost:8080/ws');

const rows = 200;
const cols = 200;
const cellSize = 40;
let grid = Array.from({ length: rows }, () => Array(cols).fill(0));

canvas.addEventListener("click", function (event) {
    const rect = canvas.getBoundingClientRect();
    const x = Math.floor((event.clientX - rect.left) / cellSize);
    const y = Math.floor((event.clientY - rect.top) / cellSize);

    socket.send(JSON.stringify({
        type: "toggle",
        data: { x, y }
    }));
});

socket.onmessage = function (event) {
    const msg = JSON.parse(event.data);
    if (msg.type === "grid") {
        grid = msg.data;
        drawGrid();
    }
};

function drawGrid() {
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    ctx.font = `${cellSize}px sans-serif`;
    ctx.textAlign = "center";
    ctx.textBaseline = "middle";

    for (let y = 0; y < rows; y++) {
        for (let x = 0; x < cols; x++) {
            const cx = x * cellSize + cellSize / 2;
            const cy = y * cellSize + cellSize / 2;
            if (grid[y][x]) {
                ctx.fillText("🌸", cx, cy); 
            } else {
                ctx.strokeStyle = '#eee';
                ctx.strokeRect(x * cellSize, y * cellSize, cellSize, cellSize);
            }
        }
    }
}

function send(type, data = null) {
    socket.send(JSON.stringify({ type, data }));
}

function toggleRunning() {
    send(running ? "stop" : "start");
    running = !running;
}

function clearGrid() {
    send("clear");
}

function randomize() {
    send("random");
}

function savePattern() {
    const name = prompt("Enter pattern name:");
    if (name) {
        send("save", { name });
    }
}

function loadPattern() {
    const name = prompt("Enter pattern name to load:");
    if (name) {
        send("load", { name });
    }
}

let running = false;
