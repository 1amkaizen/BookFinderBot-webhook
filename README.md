# BookFinderBot(webhook)


<div align="center">
    <img loading="lazy" width="30px" src="https://github.com/montasim/montasim/blob/main/media/icons/code.png " alt="code png" />
    <img loading="lazy" src="https://readme-typing-svg.demolab.com?font=Poppins&weight=700&size=30&duration=1&pause=1&color=add8e6&center=true&vCenter=true&repeat=false&width=395&height=29&lines=BookFinderBot" alt="hello nice to meet you svg" />
    <img loading="lazy" width="30px" src="https://github.com/montasim/montasim/blob/main/media/icons/layers.png " alt="layers png" />
</div>

<div align="center">
  <img src="https://media.giphy.com/media/mAgG12Pk85e1mc31HJ/giphy.gif" width="100px"/>
</div>



BookFinderBot(webhook) adalah bot pencari Ebook & Buku yang dibangun menggunakan Go, Fiber, dan Telegram Bot API. Bot ini memungkinkan pengguna untuk mencari Ebook dan Buku berdasarkan judul atau topik tertentu, serta memberikan link ulasan produk yang diminta.

> Ini adalah versi webhook dari [BookFinderBot](https://github.com/1amkaizen/BookFinderBot).
## Persyaratan

Sebelum menjalankan bot, pastikan Anda telah memenuhi persyaratan berikut:

1. Go (v1.16 atau lebih baru)
2. Fiber (v2.5.0 atau lebih baru)
3. Telegram Bot API (v5 atau lebih baru)
4. Akses ke Telegram Bot Token
5. File .env yang berisi konfigurasi bot

## Cara Mengatur Variabel Lingkungan

Untuk mengatur variabel lingkungan yang diperlukan, buatlah file bernama `.env` di direktori proyek Anda dan isi dengan informasi berikut:

```
TELEGRAM_BOT_TOKEN=TOKEN_ANDA_DISINI
ADDR=:3000
WEBHOOK_URL=https://webhookurl.app/webhook
```

Ganti `TOKEN_ANDA_DISINI` dengan token bot Telegram Anda yang diperoleh dari BotFather. Anda juga dapat mengubah port `ADDR` sesuai kebutuhan Anda.

## Cara Mendapatkan Token Bot Telegram

Untuk menggunakan Bot Telegram, Anda perlu membuat bot baru dan mendapatkan token dari BotFather, bot resmi untuk mengelola bot Telegram.

Berikut langkah-langkahnya:

1. Buka Telegram dan cari BotFather di pencarian.

2. Mulai percakapan dengan BotFather dengan menekan tombol "Mulai" atau mengirimkan pesan "/start".

3. Ketikkan perintah `/newbot` untuk membuat bot baru.

4. BotFather akan meminta Anda untuk memberikan nama bot baru. Berikan nama yang unik untuk bot Anda.

5. Setelah memberikan nama, BotFather akan meminta Anda memberikan username bot. Username bot harus diakhiri dengan kata "bot". Contoh: `book_finder_bot`.

6. Setelah Anda memberikan username, BotFather akan mengonfirmasi pembuatan bot dan memberikan token bot Anda.

7. Salin token tersebut dan tempelkan ke dalam file `.env` di variabel `TELEGRAM_BOT_TOKEN`.


## Cara Menjalankan Bot

1. Pastikan Anda telah mengatur variabel lingkungan di file `.env`.
2. Buka terminal dan navigasi ke direktori proyek Anda.
3. Jalankan perintah `go run main.go` untuk menjalankan bot.
4. Bot sekarang siap digunakan di Telegram.

## Contoh Penggunaan

Berikut adalah contoh penggunaan bot di Telegram:

1. `/start` - Memulai bot dan menampilkan pesan selamat datang.
2. `/help` - Menampilkan panduan penggunaan bot.
3. `/ulasan [judul lengkap produk]` - Mendapatkan link ulasan untuk produk yang diminta.

Contoh: `/ulasan Belajar Golang`

## Menyiapkan Data Produk dan Link Ulasan

1. Ganti isi file `products.txt` dengan produk-produk yang ingin Anda tampilkan dalam bot. Format setiap baris adalah `Nama Produk: https://linkproduk`.
2. Ganti isi file `link_reviews.txt` dengan link ulasan untuk setiap produk. Format setiap baris adalah `Nama Produk: https://linkulasan`.
3. Pastikan nama produk di `link_reviews.txt` cocok dengan nama produk di `products.txt`.


## Perintah yang Tersedia

- `/start`: Memulai percakapan dengan bot dan menampilkan pesan selamat datang.
- `/help`: Menampilkan daftar perintah yang tersedia beserta contoh penggunaannya.
- `/ulasan [nama lengkap produk]`: Mendapatkan link ulasan produk yang diinginkan.

## Mengkontribusi

Anda dapat berkontribusi pada pengembangan BookFinderBot dengan melakukan pull request ke repositori ini. Silakan buka issue untuk saran atau permintaan fitur.

## Mengirim Ulasan

Anda juga dapat memberikan ulasan langsung melalui [form ulasan kami](https://aigoretech.rf.gd/kirim-ulasan).

## Dokumentasi Tambahan

Untuk dokumentasi lebih lanjut tentang penggunaan dan cara deploy BookFinderBot, silakan lihat [dokumentasi lengkap](https://github.com/1amkaizen/BookFinderBot/wiki).

## Lisensi

[MIT](LICENSE)
