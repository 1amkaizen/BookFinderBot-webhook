package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	datauser "github.com/1amkaizen/BookFinder/user"
	"github.com/sirupsen/logrus"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

// ANSI escape codes for coloring
const (
	Reset = "\033[0m"
	Green = "\033[32m"
) // Product represents a product with multiple affiliate links
type Product struct {
	Nama  string            `json:"name"`
	Links map[string]string `json:"links"`
}

// ReviewLink represents a review link for a product
type ReviewLink struct {
	ProductName string `json:"productName"`
	Link        string `json:"link"`
}

// loadProductsFromTxt reads and parses the text file containing product data
func loadProductsFromTxt(filename string) ([]Product, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var products []Product
	var currentProduct Product
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			if currentProduct.Nama != "" {
				products = append(products, currentProduct)
				currentProduct = Product{}
			}
			continue
		}
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			if currentProduct.Links == nil {
				currentProduct.Links = make(map[string]string)
			}
			currentProduct.Links[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		} else {
			currentProduct.Nama = line
		}
	}
	if currentProduct.Nama != "" {
		products = append(products, currentProduct)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

// loadReviewLinksFromTxt reads and parses the text file containing product review links
func loadReviewLinksFromTxt(filename string) ([]ReviewLink, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var reviewLinks []ReviewLink
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		reviewLink := ReviewLink{
			ProductName: strings.TrimSpace(parts[0]),
			Link:        strings.TrimSpace(parts[1]),
		}
		reviewLinks = append(reviewLinks, reviewLink)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return reviewLinks, nil
}

// saveProductsToJson saves the products list to a JSON file
func saveProductsToJson(products []Product, filename string) error {
	data, err := json.MarshalIndent(products, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
}

// saveReviewLinksToJson saves the review links to a JSON file
func saveReviewLinksToJson(reviewLinks []ReviewLink, filename string) error {
	data, err := json.MarshalIndent(reviewLinks, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
}

// extractKeywords generates a list of keywords from the product name
func extractKeywords(name string) []string {
	words := strings.Fields(name)
	keywords := make([]string, len(words))
	for i, word := range words {
		keywords[i] = strings.ToLower(word)
	}
	return keywords
}

// findProducts searches for products based on keywords
func findProducts(products []Product, message string) []*Product {
	message = strings.ToLower(message)
	var matchingProducts []*Product
	for i := range products {
		keywords := extractKeywords(products[i].Nama)
		for _, keyword := range keywords {
			if strings.Contains(message, keyword) {
				matchingProducts = append(matchingProducts, &products[i])
				break // Break the inner loop once a match is found
			}
		}
	}
	return matchingProducts
}

// findReviewLinkByName searches for review link based on product name
func findReviewLinkByName(reviewLinks []ReviewLink, productName string) (string, bool) {
	for _, link := range reviewLinks {
		if link.ProductName == productName {
			return link.Link, true
		}
	}
	return "", false
}

func load() ([]Product, []ReviewLink, error) {
	// Load products from text file
	products, err := loadProductsFromTxt("products.txt")
	if err != nil {
		return nil, nil, fmt.Errorf("Gagal memuat produk: %v", err)
	}

	// Load review links from text file
	reviewLinks, err := loadReviewLinksFromTxt("review_links.txt")
	if err != nil {
		return nil, nil, fmt.Errorf("Gagal memuat link review: %v", err)
	}

	// Save products to JSON file
	err = saveProductsToJson(products, "products.json")
	if err != nil {
		return nil, nil, fmt.Errorf("Gagal menyimpan produk ke JSON: %v", err)
	}

	// Save review links to JSON file
	err = saveReviewLinksToJson(reviewLinks, "review_links.json")
	if err != nil {
		return nil, nil, fmt.Errorf("Gagal menyimpan link review ke JSON: %v", err)
	}

	// Load products from JSON file (optional, for consistency)
	data, err := ioutil.ReadFile("products.json")
	if err != nil {
		return nil, nil, fmt.Errorf("Gagal membaca produk dari JSON: %v", err)
	}
	err = json.Unmarshal(data, &products)
	if err != nil {
		return nil, nil, fmt.Errorf("Gagal mengurai produk dari JSON: %v", err)
	}

	// Load review links from JSON file (optional, for consistency)
	data, err = ioutil.ReadFile("review_links.json")
	if err != nil {
		return nil, nil, fmt.Errorf("Gagal membaca link review dari JSON: %v", err)
	}

	err = json.Unmarshal(data, &reviewLinks)
	if err != nil {
		return nil, nil, fmt.Errorf("Gagal mengurai link review dari JSON: %v", err)
	}

	return products, reviewLinks, nil
}

func webhook(app *fiber.App, bot *tgbotapi.BotAPI, products []Product, reviewLinks []ReviewLink) {
	var users []datauser.UserData
	// Handler untuk webhook
	app.Post("/webhook", func(c *fiber.Ctx) error {
		update := new(tgbotapi.Update)
		if err := c.BodyParser(update); err != nil {
			log.Println("Gagal memparsing update:", err)
			return err
		}

		if update.Message == nil {
			return nil
		}

		// Log informasi pengguna
		userInfo := update.Message.From
		logMessage := fmt.Sprintf(
			"%sUsername: %s%s\n%sUser ID: %d%s\n%sChat ID: %d%s\n%sMessage: %s%s",
			Green, userInfo.UserName, Reset,
			Green, userInfo.ID, Reset,
			Green, update.Message.Chat.ID, Reset,
			Green, update.Message.Text, Reset,
		)
		logrus.Info(logMessage)

		currenttime := time.Now()
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		// Dapatkan foto profil pengguna
		userProfilePhotos, err := bot.GetUserProfilePhotos(tgbotapi.UserProfilePhotosConfig{UserID: userInfo.ID})
		if err != nil {
			logrus.Error("Failed to get user profile photos:", err)
			return err
		}

		var profilePhotoURL string
		if len(userProfilePhotos.Photos) > 0 {
			photo := userProfilePhotos.Photos[0][0]
			fileID := photo.FileID
			file, err := bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
			if err != nil {
				logrus.Error("Failed to get file info:", err)
				return err
			}
			profilePhotoURL = fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", bot.Token, file.FilePath)
		}

		// Logika tanggapan terhadap pesan
		switch update.Message.Text {
		case "/start":
			// Tanggapan untuk perintah /start
			msg.Text = "📚 Selamat datang di BookFinderBot! Saya adalah bot pencari Ebook & Buku. Cari Ebook apa yang Anda butuhkan? Ketikkan judul atau topik yang Anda inginkan, dan saya akan mencarikannya untuk Anda."
			if _, err := bot.Send(msg); err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("Failed to send /start response")
			}
		case "/help":
			// Tanggapan untuk perintah /help
			msg.Text = `ℹ️ Gunakan bot ini untuk mencari Ebook & Buku. Anda cukup ketik judul atau topik yang ingin Anda cari, dan saya akan mencarikannya untuk Anda.

🔍 Contoh penggunaan:
Ketikkan "Belajar Python" untuk mencari Ebook atau Buku tentang pemrograman Python.
Ketikkan "Hacking" untuk mencari Ebook atau Buku tentang hacking.

📖 Anda juga bisa menggunakan perintah:
/ulasan [nama lengkap produk] untuk mendapatkan link ulasan produk tersebut.

⚠️ Perhatian: Judul harus sesuai, perhatikan huruf besar dan kecilnya agar mendapatkan link ulasan.

📘 Contoh penggunaan:
/ulasan Ilmu Hacking 
untuk mendapatkan link ulasan buku Ilmu Hacking.
📝 Catatan:
Kamu juga bisa memberikan ulasan di sini:
http://aigoretech.rf.gd/kirim-ulasan`
			if _, err := bot.Send(msg); err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("Failed to send /help response")
			}
		case "/ulasan":
			// Tanggapan untuk perintah /ulasan
			msg.Text = "⚠️ Mohon berikan judul lengkap buku untuk mendapatkan link ulasannya.\nContoh penggunaan: /ulasan Judul Buku"
			if _, err := bot.Send(msg); err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("Failed to send /ulasan response")
			}
		default:
			// Tanggapan untuk pesan lain
			// Panggil SaveUserDataToHTML untuk setiap pesan yang diterima
			userData := datauser.UserData{
				ID:              update.Message.Chat.ID,
				Username:        update.Message.Chat.UserName,
				FirstName:       update.Message.Chat.FirstName,
				LastName:        update.Message.Chat.LastName,
				Message:         update.Message.Text,
				Timestamp:       currenttime,
				ProfilePhotoURL: profilePhotoURL,
			}

			// Tanggapi pesan phone number
			if update.Message.Contact != nil {
				userData.PhoneNumber = update.Message.Contact.PhoneNumber
			}

			users = append(users, userData)
			err := datauser.SaveUserDataToHTML(users, "user_data.html")
			if err != nil {
				log.Println("Gagal menyimpan data pengguna:", err)
			}
			if strings.HasPrefix(update.Message.Text, "/ulasan ") {
				productName := strings.TrimPrefix(update.Message.Text, "/ulasan ")
				if link, found := findReviewLinkByName(reviewLinks, productName); found {
					msg.Text = "📘 Link ulasan untuk " + productName + ":\n" + link
				} else {
					msg.Text = "⚠️ Link ulasan untuk " + productName + " tidak ditemukan.\nKamu bisa memberikan ulasan di sini: http://aigoretech.rf.gd/kirim-ulasan"
				}
				bot.Send(msg)
			} else {
				matchingProducts := findProducts(products, update.Message.Text)
				if len(matchingProducts) > 0 {
					for _, product := range matchingProducts {
						// Buat pesan untuk setiap produk yang cocok
						msg.Text = "📖 Judul: " + product.Nama

						// Kirim satu link pratinjau dari produk pertama
						for linkName, linkURL := range product.Links {
							msg.Text += "\n🔗 [" + linkName + "](" + linkURL + ")"
							break
						}

						// Buat tombol inline untuk setiap produk
						var buttons []tgbotapi.InlineKeyboardButton
						for linkName, linkURL := range product.Links {
							button := tgbotapi.NewInlineKeyboardButtonURL(linkName, linkURL)
							buttons = append(buttons, button)
						}

						// Buat keyboard inline
						keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons)

						// Set keyboard pada pesan
						msg.ReplyMarkup = &keyboard

						// Kirim pesan
						bot.Send(msg)
					}
				} else {
					msg.Text = "⚠️ Produk tidak ditemukan."
					bot.Send(msg)
				}
			}
		}
		return nil
	})
}

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
	products, reviewLinks, err := load()
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
	webhook(app, bot, products, reviewLinks)

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
