package appinfo

type Category struct {
	Id    int    `db:"id" json:"id"`
	Title string `db:"title" json:"title"`
}

type CategoryFilter struct {
	Title string `query:"title"`
}
