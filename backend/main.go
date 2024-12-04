package main

import (
    "projectGO/database"
	"projectGO/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
    database.Connect()
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		//se tiver como false o JS nao cosnegue enviar o cookie para o front 
		AllowOrigins: "http://localhost:3001",
		AllowCredentials: true,
		
	}))
	routes.Setup(app)
	app.Listen(":3000")
}