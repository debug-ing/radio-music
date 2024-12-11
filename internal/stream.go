package internal

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

var (
	CurrentMusic string
)

func getPlaylist(folderPath string) ([]string, error) {
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}

	var playlist []string
	for _, file := range files {
		if !file.IsDir() {
			playlist = append(playlist, folderPath+"/"+file.Name())
		}
	}
	return playlist, nil
}

func StartStream(folderPath string, client *Client) {
	playlist, err := getPlaylist(folderPath)
	if err != nil {
		log.Fatal("Error reading playlist:", err)
	}

	if len(playlist) == 0 {
		log.Fatal("No files found in the folder")
	}

	for {
		for _, filePath := range playlist {
			CurrentMusic = strings.Split(filePath, "/")[1]
			file, err := os.Open(filePath)
			if err != nil {
				log.Println("Error opening file:", err)
				continue
			}

			buffer := make([]byte, 2048)
			for {
				n, err := file.Read(buffer)
				if err != nil {
					if err.Error() == "EOF" {
						break // به فایل بعدی بروید
					}
					log.Println("Error reading file:", err)
					break
				}
				client.Broadcast(buffer[:n])
				// broadcast(buffer[:n])
				time.Sleep(50 * time.Millisecond)
			}

			file.Close()
		}
	}
}
