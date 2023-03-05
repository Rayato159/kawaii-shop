package products

import "github.com/Rayato159/kawaii-shop/modules/entities"

type ProductFilter struct {
	Id     string `query:"id"`
	Search string `query:"search"` // Title & Description
	*entities.PaginateReq
	*entities.SortReq
}

type Product struct {
	Id          string             `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	CreatedAt   string             `json:"created_at"`
	UpdatedAt   string             `json:"updated_at"`
	Images      []*entities.Images `json:"images"`
}
