package entities

type Images struct {
	Id       string `db:"id" json:"id"`
	FileName string `db:"filename" json:"filename"`
	Url      string `db:"url" json:"url"`
}

type PaginateReq struct {
	Page      int `query:"page"`
	Limit     int `query:"limit"`
	TotalPage int `query:"total_page"`
	TotalItem int `query:"total_item"`
}

type SortReq struct {
	OrderBy string `query:"order_by"`
	Sort    string `query:"sort"`
}

type PaginateRes struct {
	Data      any `json:"data"`
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	TotalPage int `json:"total_page"`
	TotalItem int `json:"total_item"`
}
