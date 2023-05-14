package controllers

import (
	"strconv"
	"strings"
	"time"
	"github.com/gofiber/fiber/v2"
	"notes/models"
	"notes/initializers"
	"gorm.io/gorm"
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

func FindNoteById(c *fiber.Ctx) error {
	noteId := c.Params("noteId")

	var note models.Note

	result := initializers.DB.First(&note, "id = ?", noteId)
	
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{ "status" : "Falha", "message" : "Não existe anotação com esse ID" })
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{ "status": "Falha", "message" : err.Error() })
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{ "status" : "success", "data" : fiber.Map{"note" : note} })
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

func UpdateNote(c *fiber.Ctx) error {
	noteId := c.Params("noteId")

	var payload *models.UpdateNoteSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{ "status": "Falha", "message" : err.Error() })
	}

	var note models.Note
	result := initializers.DB.First(&note, "id = ?", noteId)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{ "status" : "Falha", "message" : "Não existe anotação com esse ID" })
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{ "status" : "Falha", "message" : err.Error() })
	}

	updates := make(map[string]interface{})
	if payload.Title != "" {
		updates["title"] = payload.Title
	}
	if payload.Content != "" {
		updates["content"] = payload.Content
	}

	updates["updated_at"] = time.Now()

	initializers.DB.Model(&note).Updates(updates)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{ "status" : "success", "data" : fiber.Map{"note" : note} })
}

func DeleteNote(c *fiber.Ctx) error {
	noteId := c.Params("noteId")
	
	result := initializers.DB.Delete(&models.Note{}, "id = ?", noteId)
	
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{ "status": "Falha", "message" : "Não foi encontrado anotação com esse ID." })
	} else if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{ "status" : "Erro", "message" : result.Error })
	}

	return c.SendStatus(fiber.StatusNoContent)
}