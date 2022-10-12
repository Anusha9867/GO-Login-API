package routes

import (
	"login/controllers"

	"github.com/gofiber/fiber/v2"
)

// func Setup(router fiber.Router) {
// 	router.Get("/", func(c *fiber.Ctx) error {
// 		return c.SendString("respond with a resource")
// 	})
// }

func Setup(app *fiber.App) {

	app.Post("/api/register", controllers.Register)
	app.Post("/api/login", controllers.Login)
	app.Delete("/api/delete/:id", controllers.Delete)
	app.Post("/api/logout", controllers.Logout)
	app.Get("/api/user", controllers.User)
	app.Put("/api/update/:id", controllers.Update)
	app.Get("/api/getusers/:id", controllers.GetUser)

}
