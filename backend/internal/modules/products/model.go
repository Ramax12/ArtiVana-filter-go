package products

import (
	"ArtiVana/internal/modules/dictionaries"
)

type FilterDataItem struct {
	dictionaries.ProductOption
	Count int `json:"count" example:"25"`
}

type FilterDataItemCharacteristics struct {
	dictionaries.ProductOptionChars
	Options []FilterDataItem `json:"options"`
}

type FilterDataPrice struct {
	BaseRange [2]int `json:"base_range" example:"0,500"`
	Range     [2]int `json:"range" example:"100,300"`
}

type FilterDataRating struct {
	Count int `json:"count" example:"4"`
}

type FilterDataStandard struct {
	Subsubcategories []FilterDataItem `json:"subsubcategories"`
	Brands           []FilterDataItem `json:"brands"`
	Price            FilterDataPrice  `json:"price"`
	Rating           FilterDataRating `json:"rating"`
}

type FilterData struct {
	Standard        FilterDataStandard              `json:"standard"`
	Characteristics []FilterDataItemCharacteristics `json:"characteristics"`
}

type FilterMeta struct {
	Count   int      `json:"count" example:"120"`
	Filters FilterData `json:"filters"`
}

type CreateFilter struct {
	Brands           []int            `json:"brands,omitempty" example:"1, 2, 3"`
	Subcategory      int              `json:"subcategory" example:"5"`
	Subsubcategory   *int             `json:"subsubcategory,omitempty" example:"12"`
	Subsubcategories []int            `json:"subsubcategories,omitempty" example:"10, 11, 12"`
	MinPrice         *int             `json:"min_price,omitempty" example:"100"`
	MaxPrice         *int             `json:"max_price,omitempty" example:"500"`
	Rating           *int             `json:"rating,omitempty" example:"1"`
	Characteristics  map[string][]int `json:"characteristics,omitempty"`
	Sort             *string            `json:"sort,omitempty" example:"price_asc"`
	Page             *int             `json:"page,omitempty" example:"1"`
}

type FacetCount struct {
	ID    int
	Count int
}

type FilterFacets struct {
	Count            int
	BasePriceMin     int
	BasePriceMax     int
	PriceMin         int
	PriceMax         int
	RatingCount      int
	Brands           []FacetCount
	Subsubcategories []FacetCount
	Characteristics  map[string][]FacetCount
}
