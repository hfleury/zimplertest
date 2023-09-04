package main

import (
	"encoding/json"
	"fmt"
	"github.com/hfleury/zimplertest/internal/web"
	"log"
	"sort"
)

func main() {
	url := "https://candystore.zimpler.net/#candystore-customers"

	webHandler := web.NewWebHandler()

	doc, err := webHandler.FetchDataFromWebsite(url)
	if err != nil {
		log.Fatal(err)
		return
	}

	targetClass := "top.customers summary"
	topCustomers, err := webHandler.ExtractDataFromTable(doc, targetClass)
	if err != nil {
		log.Fatal(err)
		return
	}

	sort.SliceStable(topCustomers, func(i, j int) bool {
		return topCustomers[i].TotalSnacks > topCustomers[j].TotalSnacks
	})

	jsonData, err := json.Marshal(topCustomers)
	if err != nil {
		log.Fatal("Error marshaling JSON:", err)
		return
	}

	fmt.Println(string(jsonData))
}
