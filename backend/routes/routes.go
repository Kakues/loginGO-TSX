package routes

import(
	"projectGO/controllers"
	"github.com/gofiber/fiber/v2"

)



func Setup(app *fiber.App){

	api := app.Group("/api")

	api.Post("/register", controllers.Register)
	api.Post("/login", controllers.Login)
	api.Get("/user", controllers.User)
	api.Post("/logout", controllers.Logout)
	api.Post("/forgot", controllers.Forgot)
	api.Post("/reset", controllers.Reset)
}