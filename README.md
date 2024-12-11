# 🎵 Live Radio Streaming Server

This project is a simple online radio server that allows you to stream your own playlist live for listeners. People can join at any time to enjoy the continuous music stream without restarting tracks from the beginning.


## Features

1. Live Music Streaming: Play music continuously from your playlist as a live stream.
2. Custom Playlist Support: Add your favorite tracks to a folder, and the server will stream them for listeners.
2. Professional Monitoring: Built-in support for Prometheus to monitor requests and server performance.

## Running the Project
```bash
git clone https://github.com/debug-ing/radio-music
cd radio-music
go mod tidy
```

template config.toml 
```
port : set port
folder : name music folder in base project 
```

```bash
make build
./main --config=./config.toml 
```