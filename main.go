package main

import (
	"log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"notes/initializers"
	"notes/controllers"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatalln("Falha ao carregar variáveis de ambiente! \n", err.Error())
	}
	initializers.ConnectDB(&config)
}

func main()  {
	app := fiber.New()
	micro := fiber.New()

	app.Mount("/api", micro)
	app.Use(logger.New())
	app.Use(cors.New(cors.Config {
		AllowOrigins: "http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PATCH, DELETE",
		AllowCredentials: true,
	}))

	micro.Route("/notes", func(router fiber.Router) {
		router.Post("/", controllers.CreateNoteHandler)
	})

	micro.Get("/healthchecker", healthChecker)

	log.Fatal(app.Listen(":8000"))
}

func healthChecker(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"message": "Welcome to Golang, Fiber and GORM!",
	})
}