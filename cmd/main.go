package main

import (
	"log"

	"github.com/debug-ing/radio-music/config"
	"github.com/debug-ing/radio-music/internal"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	config := config.LoadConfig()
	client := internal.NewClient()
	go internal.StartStream(config.App.Folder, client)
	r := gin.Default()
	r.GET("/radio", client.HandleClientGin)
	r.GET("/status", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name_music": internal.CurrentMusic,
		})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	log.Println("Radio server is running on http://localhost:8080/radio")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
