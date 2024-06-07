package handler

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

// ReviewLink represents a review link for a product
type ReviewLink struct {
	ProductName string `json:"productName"`
	Link        string `json:"link"`
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

// saveReviewLinksToJson saves the review links to a JSON file
func saveReviewLinksToJson(reviewLinks []ReviewLink, filename string) error {
	data, err := json.MarshalIndent(reviewLinks, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
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
