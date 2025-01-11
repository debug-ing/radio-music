package main

import (
	"flag"
	"log"

	"github.com/debug-ing/radio-music/config"
	"github.com/debug-ing/radio-music/internal"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	configPath := flag.String("config", "config.toml", "config file")
	flag.Parse()
	config := config.LoadConfig(*configPath)
	client := internal.NewClient()
	go internal.StartStream(config.App.Folder, client)
	r := gin.Default()
	r.GET("/radio", client.HandleClientGin)
	r.GET("/status", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"music": internal.CurrentMusic,
			"info":  internal.Info,
		})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	log.Println("Radio server is running on http://localhost:8080/radio")
	if err := r.Run(":" + config.App.Port); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
