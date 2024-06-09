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

func extractKeywords(name string) []string {
	// Use jdokato/prose/v2 for tokenization
	doc, _ := prose.NewDocument(name)
	var keywords []string
	for _, tok := range doc.Tokens() {
		keywords = append(keywords, tok.Text)
	}
	return keywords
}

func extractEntities(message string) []string {
	doc, _ := prose.NewDocument(message)
	var entities []string
	for _, entity := range doc.Entities() {
		entities = append(entities, entity.Text)
	}
	return entities
}

func findProducts(products []Product, message string) []*Product {
	message = strings.ToLower(message)
	var matchingProducts []*Product

	if len(message) < 2 {
		return matchingProducts
	}

	// Tokenisasi pesan menggunakan NLP
	keywords := tokenize(message)

	// Simpan kata kunci yang unik
	uniqueKeywords := make(map[string]bool)
	for _, keyword := range keywords {
		uniqueKeywords[keyword] = true
	}

	// Cari produk yang cocok berdasarkan kata kunci dalam nama produk
	for i := range products {
		productName := strings.ToLower(products[i].Nama)
		if matchKeywords(productName, uniqueKeywords) {
			matchingProducts = append(matchingProducts, &products[i])
		}
	}

	return matchingProducts
}

// Tokenisasi menggunakan pemisah kata sederhana
func tokenize(message string) []string {
	return strings.Fields(message)
}

// Cocokkan kata kunci dalam nama produk
func matchKeywords(productName string, keywords map[string]bool) bool {
	for keyword := range keywords {
		if strings.Contains(productName, keyword) {
			return true
		}
	}
	return false
}
