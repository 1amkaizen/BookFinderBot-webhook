package handler

import (
	"fmt"
	"log"
	"strings"
	"time"

	datauser "github.com/1amkaizen/BookFinderBot/user"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

func handleMessage(update *tgbotapi.Update, bot *tgbotapi.BotAPI, products []Product, reviewLinks []ReviewLink) {
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
	var botResponse string

	profilePhotoURL := getProfilePhotoURL(bot, userInfo.ID)

	switch update.Message.Text {
	case "/start", "/help", "/ulasan":
		botResponse = processCommand(update.Message.Text)
		msg.Text = botResponse
	default:
		handleGeneralMessage(update, bot, products, reviewLinks, &msg, &botResponse)
	}

	if msg.Text != "" {
		if _, err := bot.Send(msg); err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("Failed to send message")
		}
	}

	saveUserData(update, botResponse, currenttime, profilePhotoURL)
}

func processCommand(command string) string {
	switch command {
	case "/start":
		return "📚 Selamat datang di BookFinderBot! Saya adalah bot pencari Ebook & Buku. Cari Ebook apa yang Anda butuhkan? Ketikkan judul atau topik yang Anda inginkan, dan saya akan mencarikannya untuk Anda."
	case "/help":
		return `ℹ️ Gunakan bot ini untuk mencari Ebook & Buku. Anda cukup ketik judul atau topik yang ingin Anda cari, dan saya akan mencarikannya untuk Anda.

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
	case "/ulasan":
		return "⚠️ Mohon berikan judul lengkap buku untuk mendapatkan link ulasannya.\nContoh penggunaan: /ulasan Judul Buku"
	}
	return ""
}

func handleGeneralMessage(update *tgbotapi.Update, bot *tgbotapi.BotAPI, products []Product, reviewLinks []ReviewLink, msg *tgbotapi.MessageConfig, botResponse *string) {
	if strings.HasPrefix(update.Message.Text, "/ulasan ") {
		productName := strings.TrimPrefix(update.Message.Text, "/ulasan ")
		if link, found := findReviewLinkByName(reviewLinks, productName); found {
			*botResponse = "📘 Link ulasan untuk " + productName + ":\n" + link
		} else {
			*botResponse = "⚠️ Link ulasan untuk " + productName + " tidak ditemukan.\nKamu bisa memberikan ulasan di sini: http://aigoretech.rf.gd/kirim-ulasan"
		}
		msg.Text = *botResponse
	} else {
		matchingProducts := findProducts(products, update.Message.Text)
		if len(matchingProducts) > 0 {
			for _, product := range matchingProducts {
				*botResponse = "📖 Judul: " + product.Nama
				for linkName, linkURL := range product.Links {
					*botResponse += "\n🔗 [" + linkName + "](" + linkURL + ")"
					break
				}

				var buttons []tgbotapi.InlineKeyboardButton
				for linkName, linkURL := range product.Links {
					button := tgbotapi.NewInlineKeyboardButtonURL(linkName, linkURL)
					buttons = append(buttons, button)
				}

				keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons)
				msg.ReplyMarkup = &keyboard

				if *botResponse != "" {
					msg.Text = *botResponse
					if _, err := bot.Send(msg); err != nil {
						logrus.WithFields(logrus.Fields{
							"error": err,
						}).Error("Failed to send product message")
					}
				} else {
					logrus.Error("Bot response is empty, not sending message")
				}
			}
		} else {
			*botResponse = "⚠️ Produk tidak ditemukan."
			msg.Text = *botResponse
		}
	}
}

func getProfilePhotoURL(bot *tgbotapi.BotAPI, userID int64) string {
	userProfilePhotos, err := bot.GetUserProfilePhotos(tgbotapi.UserProfilePhotosConfig{UserID: userID})
	if err != nil {
		logrus.Error("Failed to get user profile photos:", err)
		return ""
	}

	if len(userProfilePhotos.Photos) > 0 {
		photo := userProfilePhotos.Photos[0][0]
		fileID := photo.FileID
		file, err := bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
		if err != nil {
			logrus.Error("Failed to get file info:", err)
			return ""
		}
		return fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", bot.Token, file.FilePath)
	}

	return ""
}

func saveUserData(update *tgbotapi.Update, botResponse string, currenttime time.Time, profilePhotoURL string) {
	filename := "user_data.json"

	users, err := datauser.LoadUserData(filename)
	if err != nil {
		log.Println("Gagal memuat data pengguna:", err)
	}

	message := datauser.Message{
		Content:   update.Message.Text,
		Sender:    "user",
		Timestamp: currenttime,
	}

	updatedUsers, err := datauser.AddUserMessage(users, update.Message.Chat.ID, message)
	if err != nil {
		// Jika pengguna baru, tambahkan data pengguna baru
		newUser := datauser.UserData{
			ID:              update.Message.Chat.ID,
			Username:        update.Message.From.UserName,
			FirstName:       update.Message.From.FirstName,
			LastName:        update.Message.From.LastName,
			ProfilePhotoURL: profilePhotoURL,
			Messages:        []datauser.Message{message},
		}
		updatedUsers = append(users, newUser)
	}

	botMessage := datauser.Message{
		Content:   botResponse,
		Sender:    "bot",
		Timestamp: currenttime,
	}

	updatedUsers, err = datauser.AddUserMessage(updatedUsers, update.Message.Chat.ID, botMessage)
	if err != nil {
		log.Println("Gagal menambahkan pesan bot ke data pengguna:", err)
	}

	err = datauser.SaveUserData(filename, updatedUsers)
	if err != nil {
		log.Println("Gagal menyimpan data pengguna:", err)
	}

	err = datauser.SaveUserDataToHTML(updatedUsers, "user_data.html")
	if err != nil {
		log.Println("Gagal menyimpan data pengguna ke HTML:", err)
	}
}
