package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func Load() ([]Product, []ReviewLink, error) {
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
