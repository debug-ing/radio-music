package internal

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dhowden/tag"
)

type InfoMusic struct {
	Title    string
	Artist   string
	Album    string
	Year     int
	Composer string
	Genre    string
	Lyrics   string
}

var (
	CurrentMusic string
	Info         *InfoMusic
)

func getPlaylist(folderPath string) ([]string, error) {
	files, err := os.ReadDir(folderPath)
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
			CurrentMusic = strings.Split(filePath, "/")[len(strings.Split(filePath, "/"))-1]
			file, err := os.Open(filePath)
			if err != nil {
				log.Println("Error opening file:", err)
				continue
			}
			defer file.Close()
			metadata, err := tag.ReadFrom(file)
			if err != nil {
				log.Println("Error reading metadata:", err)
				return
			}
			fmt.Println(metadata.Title())
			Info = &InfoMusic{
				Title:    metadata.Title(),
				Artist:   metadata.Artist(),
				Album:    metadata.Album(),
				Year:     metadata.Year(),
				Composer: metadata.Composer(),
				Genre:    metadata.Genre(),
				Lyrics:   metadata.Lyrics(),
			}
			buffer := make([]byte, 2048)
			for {
				n, err := file.Read(buffer)
				if err != nil {
					if err.Error() == "EOF" {
						break
					}
					log.Println("Error reading file:", err)
					break
				}
				client.Broadcast(buffer[:n])
				time.Sleep(50 * time.Millisecond)
			}
		}
	}
}
