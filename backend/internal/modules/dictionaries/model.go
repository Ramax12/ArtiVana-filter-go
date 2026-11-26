package dictionaries

type ProductOption struct {
	ID   int64  `json:"id" example:"1"`
	Name string `json:"name" example:"Acrylic"`
}

type ProductOptionChars struct {
	ProductOption
	Value string `json:"value" example:"acrylic"`
}

type ProductOptionCharsOption struct {
	ProductOption
	ParentID int64 `json:"parent_id" example:"1"`
}
