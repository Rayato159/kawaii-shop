package orders

import (
	"github.com/Rayato159/kawaii-shop/modules/entities"
	"github.com/Rayato159/kawaii-shop/modules/products"
)

type OrderFilter struct {
	Search    string `query:"search"` // user_id, address, contract, products {id, title, description}
	StartDate string `query:"start_date"`
	EndDate   string `query:"end_date"`
	*entities.PaginateReq
	*entities.SortReq
}

type Order struct {
	Id           string           `json:"id"`
	UserId       string           `json:"user_id"`
	TransterSlip *TransterSlip    `json:"image"`
	Products     []*ProductsOrder `json:"products"`
	Address      string           `json:"address"`
	Contract     string           `json:"contract"`
	Status       string           `json:"status"`
	TotalPaid    float64          `json:"total_paid"`
	CreatedAt    string           `json:"created_at"`
	UpdatedAt    string           `json:"updated_at"`
}

type TransterSlip struct {
	Id        string `json:"id"`
	FileName  string `json:"filename"`
	Url       string `json:"url"`
	CreatedAt string `json:"created_at"`
}

type ProductsOrder struct {
	Id      string            `json:"id"`
	Qty     int               `json:"qty"`
	Product *products.Product `json:"product"`
}
