package main

import (
	"log"
	app "order-v2-microservice/internal/bootstrap"
	"order-v2-microservice/internal/routers"
	"os"
)

func main() {
	h := app.Bootstrap()
	router := routers.SetupRouter(h)

	port := os.Getenv("APP_PORT")
	log.Println("Server port:", port)
	log.Fatal(router.Start(":" + port))
}
