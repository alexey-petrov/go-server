package routes

import (
	"github.com/alexey-petrov/go-server/server/routes/todoRoutes"
	"github.com/alexey-petrov/go-server/server/routes/userRoutes"
	"github.com/gofiber/fiber/v2"
)


func SetRoutes(app *fiber.App) {

	initEndpoints(app)

	todoRoutes.TodoRoutes(app)
	userRoutes.UserRoutes(app)
}




func initEndpoints(app *fiber.App) {
	app.Get("api/healthcheck", helloHandler)
}


func helloHandler(c *fiber.Ctx) error {
	return c.SendString("Access Granted")
}