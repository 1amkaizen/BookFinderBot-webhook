package handler

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/fiber/v2"
)

func Webhook(app *fiber.App, bot *tgbotapi.BotAPI, products []Product, reviewLinks []ReviewLink) {
	app.Post("/webhook", func(c *fiber.Ctx) error {
		return handleWebhook(c, bot, products, reviewLinks)
	})
}

func handleWebhook(c *fiber.Ctx, bot *tgbotapi.BotAPI, products []Product, reviewLinks []ReviewLink) error {
	update := new(tgbotapi.Update)
	if err := c.BodyParser(update); err != nil {
		log.Println("Gagal memparsing update:", err)
		return err
	}

	if update.Message == nil {
		return nil
	}

	handleMessage(update, bot, products, reviewLinks)
	return nil
}
