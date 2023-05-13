package controllers

import (
	"strings"
	"time"
	"github.com/gofiber/fiber/v2"
	"notes/models"
	"notes/initializers"
)

func CreateNoteHandler(c *fiber.Ctx) error {
	var payload *models.CreateNoteSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{ "status" : "Falha", "message" : err.Error() })
	}

	errors := models.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	now := time.Now()
	newNote := models.Note {
		Title: payload.Title,
		Content: payload.Content,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := initializers.DB.Create(&newNote)
	if result.Error != nil && strings.Contains(result.Error.Error(), "Duplicate entry") {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{ "status" : "Falha", "message" : "Título já está sendo utilizado! Por favor mude o título." })
	} else if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{ "status" : "Erro", "message" : result.Error.Error() })
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{ "status": "success", "data" : fiber.Map{"note" : newNote} })
}