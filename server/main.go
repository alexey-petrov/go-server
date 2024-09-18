package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alexey-petrov/go-server/server/db"
	"github.com/alexey-petrov/go-server/server/routes"
	"github.com/alexey-petrov/go-server/server/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func establishdbConnection() {
	fmt.Println("Establishing Gorm DB connection")
	// Initialize DB connection
	db.InitDB()

	db.CreateTable()
	db.CreateJTITable()
	db.CreateSearchSettingsTable()
}

func main() {
	// Connect to the database
	establishdbConnection()

	app := fiber.New(fiber.Config{
		IdleTimeout: 5,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	app.Use(compress.New())

	routes.SetRoutes(app)
	utils.StartCronJobs()
	handleLogFatal(app)

	go func() {
		if error := app.Listen(":4000"); error != nil {
			log.Panic(error)
		}
	}()

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c // Block the main thread until a signal is received/interrupted

	app.Shutdown()
	fmt.Println("Shutting down the server")
}

func handleLogFatal(app *fiber.App) {
	log.Fatal(app.Listen(":4000"))

}
