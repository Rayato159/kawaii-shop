package orders

import (
	"github.com/Rayato159/kawaii-shop/modules/entities"
	"github.com/Rayato159/kawaii-shop/modules/products"
)

type OrderFilter struct {
	Search    string `query:"search"` // user_id, address, contract
	Status    string `query:"status"`
	StartDate string `query:"start_date"`
	EndDate   string `query:"end_date"`
	*entities.PaginateReq
	*entities.SortReq
}

type Order struct {
	Id           string           `db:"id" json:"id"`
	UserId       string           `db:"user_id" json:"user_id"`
	TransterSlip *TransterSlip    `db:"transfer_slip" json:"transfer_slip"`
	Products     []*ProductsOrder `json:"products"`
	Address      string           `db:"address" json:"address"`
	Contact      string           `db:"contact" json:"contact"`
	Status       string           `db:"status" json:"status"`
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
	Id      string            `db:"id" json:"id"`
	Qty     int               `db:"qty" json:"qty"`
	Product *products.Product `db:"product" json:"product"`
}

type UpdateOrderReq struct {
	OrderId      string        `db:"order_id" json:"order_id"`
	Status       string        `db:"status" json:"status"`
	TransterSlip *TransterSlip `db:"transfer_slip" json:"transfer_slip"`
}
