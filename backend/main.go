package main

import (
	"fmt"
	"log"

	"backend/Configs"
	"backend/controllers"
	"github.com/gin-gonic/gin"
)

func init() {
	Configs.IntializeDotEnv()
	Configs.IntializeMongodb()
	Configs.IntializeGcpStorage()
}

func main() {
	app := gin.Default()
	app.POST("/Videos", controllers.StoreVideo)
	app.GET("/Videos", controllers.GetAnnotatedVideo)

	err := app.Run()
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	fmt.Println("Server running on port:8080")
}
