package web

import (
	"fmt"
	"github.com/hfleury/zimplertest/pkg/model"
	"golang.org/x/net/html"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFetchDataFromWebsite(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "<html><body>Mocked server test</body></html>")
	}))
	defer ts.Close()

	wh := NewWebHandler()

	parsedHTML, err := wh.FetchDataFromWebsite(ts.URL)
	if err != nil {
		t.Errorf("Error fetching data: %v", err)
	}

	if parsedHTML == nil {
		t.Error("Parsed HTML is nil")
	}
}

func TestExtractDataFromTable(t *testing.T) {
	htmlCandy := `
		<table class="top.customers summary">
			<tr>
				<td x-total-candy="11">Aadya</td>
			</tr>
			<tr>
				<td x-total-candy="208">Annika</td>
			</tr>
		</table>`
	parsedHTML, _ := html.Parse(strings.NewReader(htmlCandy))

	wh := NewWebHandler()

	topCustomers, err := wh.ExtractDataFromTable(parsedHTML, "top.customers summary")
	if err != nil {
		t.Errorf("Error extracting data: %v", err)
	}

	expectedCustomers := []*model.Customer{
		{Name: "Aadya", TotalSnacks: 11},
		{Name: "Annika", TotalSnacks: 208},
	}

	for i, customer := range topCustomers {
		if customer.Name != expectedCustomers[i].Name ||
			customer.TotalSnacks != expectedCustomers[i].TotalSnacks {
			t.Errorf("Mismatch in extracted data for customer #%d", i)
		}
	}
}
