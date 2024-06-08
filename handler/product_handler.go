package handler

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/jdkato/prose/v2"
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

// matchKeywords memeriksa apakah nama produk cocok dengan setidaknya satu kata kunci
func matchKeywords(productName string, keywords []string) bool {
	for _, keyword := range keywords {
		if strings.Contains(strings.ToLower(productName), keyword) {
			return true
		}
	}
	return false
}

// tokenize menggunakan prose untuk tokenisasi NLP yang lebih canggih
func tokenize(message string) []string {
	doc, _ := prose.NewDocument(message)
	tokens := doc.Tokens()
	var words []string
	for _, token := range tokens {
		if token.Text != "untuk" && token.Text != "pemula" {
			words = append(words, token.Text)
		}
	}
	return words
}

// extractKeywords menggunakan Prose untuk mengekstrak kata-kata dari nama produk
func extractKeywords(name string) []string {
	doc, _ := prose.NewDocument(name)
	var keywords []string
	for _, tok := range doc.Tokens() {
		if tok.Tag == "Noun" { // hanya mengambil kata benda
			keywords = append(keywords, tok.Text)
		}
	}
	return keywords
}

// findProducts menggunakan Prose untuk memperkaya pencarian produk
func findProducts(products []Product, message string) []*Product {
	message = strings.ToLower(message)
	var matchingProducts []*Product

	// Pemeriksaan jika pesan hanya terdiri dari satu kata atau lebih
	if len(message) < 2 {
		return matchingProducts
	}

	// Ekstrak entitas dari pesan
	doc, _ := prose.NewDocument(message)
	var entities []string
	for _, ent := range doc.Entities() {
		entities = append(entities, ent.Text)
	}

	for i := range products {
		productName := strings.ToLower(products[i].Nama)
		// Gunakan Prose untuk ekstraksi kata-kata dari nama produk
		productKeywords := extractKeywords(productName)
		// Cocokkan kata kunci produk dengan entitas dalam pesan pengguna
		for _, keyword := range productKeywords {
			for _, entity := range entities {
				if keyword == entity {
					matchingProducts = append(matchingProducts, &products[i])
					break
				}
			}
		}
	}

	return matchingProducts
}
