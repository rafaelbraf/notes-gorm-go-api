package controllers

import (
	"strconv"
	"strings"
	"time"
	"github.com/gofiber/fiber/v2"
	"notes/models"
	"notes/initializers"
)

func FindNotes(c *fiber.Ctx) error {
	var page = c.Query("page", "1")
	var limit = c.Query("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var notes []models.Note
	results := initializers.DB.Limit(intLimit).Offset(offset).Find(&notes)
	if results.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{ "status": "Erro", "message": results.Error })
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{ "status" : "success", "results" : len(notes), "notes" : notes })
}

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