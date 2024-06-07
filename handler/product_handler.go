package handler

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"unicode"
)

// ANSI escape codes for coloring
const (
	Reset = "\033[0m"
	Green = "\033[32m"
)

// Product represents a product with multiple affiliate links
type Product struct {
	Nama  string            `json:"name"`
	Links map[string]string `json:"links"`
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

// saveProductsToJson saves the products list to a JSON file
func saveProductsToJson(products []Product, filename string) error {
	data, err := json.MarshalIndent(products, "", "  ")
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

func findProducts(products []Product, message string) []*Product {
	message = strings.ToLower(message)
	var matchingProducts []*Product

	// Split pesan pengguna menjadi kata-kata individual
	keywords := tokenize(message)

	// Membuat map untuk menyimpan produk yang sudah ditemukan
	foundProducts := make(map[string]bool)

	for i := range products {
		productName := strings.ToLower(products[i].Nama)
		if matchKeywords(productName, keywords) && !foundProducts[productName] {
			matchingProducts = append(matchingProducts, &products[i])
			foundProducts[productName] = true // Tandai produk sebagai sudah ditemukan
		}
	}

	return matchingProducts
}

// tokenize membagi pesan menjadi kata-kata individual dengan mengambil tanda baca ke dalam pertimbangan
func tokenize(message string) []string {
	var tokens []string
	builder := strings.Builder{}

	for _, char := range message {
		if unicode.IsLetter(char) {
			builder.WriteRune(char)
		} else {
			if builder.Len() > 0 {
				tokens = append(tokens, builder.String())
				builder.Reset()
			}
		}
	}

	if builder.Len() > 0 {
		tokens = append(tokens, builder.String())
	}

	return tokens
}

// matchKeywords memeriksa apakah nama produk cocok dengan setidaknya satu kata kunci
func matchKeywords(productName string, keywords []string) bool {
	for _, keyword := range keywords {
		if !strings.Contains(productName, keyword) {
			return false
		}
	}
	return true
}
