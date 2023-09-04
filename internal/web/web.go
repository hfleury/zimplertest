package web

import (
	"fmt"
	"github.com/hfleury/zimplertest/pkg/model"
	"golang.org/x/net/html"
	"net/http"
	"strconv"
	"strings"
)

type WebHandler struct {
}

func NewWebHandler() *WebHandler {
	return &WebHandler{}
}

// FetchDataFromWebsite fetches data from a specified website and returns the parsed HTML.
func (wh *WebHandler) FetchDataFromWebsite(url string) (*html.Node, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making the request: %v", err)
	}
	defer response.Body.Close()

	parsedHTML, err := html.Parse(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error parsing HTML: %v", err)
	}

	return parsedHTML, nil
}

// ExtractDataFromTable extracts data from the specified table in the HTML document.
func (wh *WebHandler) ExtractDataFromTable(doc *html.Node, targetClass string) ([]*model.Customer, error) {
	var topCustomers []*model.Customer

	var extract func(*html.Node)
	extract = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "table" {
			for _, attr := range n.Attr {
				if attr.Key == "class" && strings.Contains(attr.Val, targetClass) {
					wh.extractFromTable(n, &topCustomers)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extract(c)
		}
	}

	extract(doc)

	return topCustomers, nil
}

func (wh *WebHandler) extractFromTable(tableNode *html.Node, topCustomers *[]*model.Customer) {
	var extract func(*html.Node)
	var extractedData *model.Customer

	extract = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "tr" {
			extractedData = &model.Customer{}
		}

		if n.Type == html.ElementNode && n.Data == "td" {
			for _, attr := range n.Attr {
				if attr.Key == "x-total-candy" {
					totalCandy, err := strconv.Atoi(attr.Val)
					if err != nil {
						fmt.Printf("Error converting total candy %v", err)
					}
					extractedData.TotalSnacks = totalCandy
				}
			}
			if n.FirstChild != nil {
				if extractedData.Name == "" {
					extractedData.Name = n.FirstChild.Data
				} else {
					// Check if this is the second td element and store it as FavoriteSnack
					if extractedData.FavoriteSnack == "" {
						extractedData.FavoriteSnack = n.FirstChild.Data
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extract(c)
		}

		if n.Type == html.ElementNode && n.Data == "tr" && extractedData.Name != "" && extractedData.FavoriteSnack != "" {
			*topCustomers = append(*topCustomers, extractedData)
		}
	}

	extract(tableNode)
}
