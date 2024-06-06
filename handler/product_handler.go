package handler

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
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
