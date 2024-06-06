package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/1amkaizen/BookFinderBot/handler"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	// Setup logrus
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.DebugLevel)

	// Load ENV variables
	err := godotenv.Load()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Failed to load .env file")
	}

	// Panggil fungsi load untuk mendapatkan produk dan review
	products, reviewLinks, err := handler.Load()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Failed to load data")
	}

	// Inisialisasi bot Telegram
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		logrus.Panic("TELEGRAM_BOT_TOKEN is not set")
	}
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		logrus.Panic(err)
	}

	// Mendapatkan URL webhook dari environment variables
	webhookURL := os.Getenv("WEBHOOK_URL")
	if webhookURL == "" {
		logrus.Fatal("WEBHOOK_URL is not set")
	}

	// Membuat payload untuk pengaturan webhook
	payload := struct {
		URL string `json:"url"`
	}{
		URL: webhookURL,
	}

	// Mengubah payload menjadi format JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Failed to convert payload to JSON")
	}

	// Mengirim permintaan HTTP POST untuk mengatur webhook
	resp, err := http.Post("https://api.telegram.org/bot"+botToken+"/setWebhook", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Failed to send set webhook request")
	}
	defer resp.Body.Close()

	// Mengecek kode status respons
	if resp.StatusCode != http.StatusOK {
		logrus.Fatalf("Failed to set webhook. Status code: %d", resp.StatusCode)
	}

	logrus.Info("Webhook successfully set")

	bot.Debug = true

	// Inisialisasi GoFiber
	app := fiber.New()

	// Panggil fungsi webhook dengan menyediakan app, bot, products, dan reviewLinks
	handler.Webhook(app, bot, products, reviewLinks)

	// Endpoint untuk melayani file HTML
	app.Get("/html", func(c *fiber.Ctx) error {
		return c.SendFile("user_data.html")
	})

	// Tentukan alamat dan port
	addr := ":3000"
	if envAddr := os.Getenv("ADDR"); envAddr != "" {
		addr = envAddr
	}

	// Jalankan server
	logrus.Fatal(app.Listen(addr))
}
