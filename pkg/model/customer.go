package model

type Customer struct {
	Name          string `json:"name"`
	TotalSnacks   int    `json:"totalSnacks"`
	FavoriteSnack string `json:"favouriteSnack"`
}
