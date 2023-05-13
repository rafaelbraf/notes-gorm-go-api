package main

import (
	"log"
	"github.com/gofiber/fiber/v2"
	"notes/initializers"
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
	app.Get("/api/healthchecker", healthChecker)

	log.Fatal(app.Listen(":8000"))
}

func healthChecker(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"message": "Welcome to Golang, Fiber and GORM!",
	})
}