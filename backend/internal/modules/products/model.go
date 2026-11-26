package products

import (
	"ArtiVana/internal/modules/dictionaries"
)

type FilterDataItem struct {
	dictionaries.ProductOption
	Count int64 `json:"count" example:"25"`
}

type FilterDataItemCharacteristics struct {
	dictionaries.ProductOptionChars
	Options []FilterDataItem `json:"options"`
}

type FilterDataPrice struct {
	BaseRange [2]int64 `json:"base_range" example:"0,500"`
	Range     [2]int64 `json:"range" example:"100,300"`
}

type FilterDataRating struct {
	Count int64 `json:"count" example:"4"`
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
	Count   int64      `json:"count" example:"120"`
	Filters FilterData `json:"filters"`
}

type CreateFilter struct {
	Brands           []int64            `json:"brands,omitempty" example:"1, 2, 3"`
	Subcategory      int64              `json:"subcategory" example:"5"`
	Subsubcategory   *int64             `json:"subsubcategory,omitempty" example:"12"`
	Subsubcategories []int64            `json:"subsubcategories,omitempty" example:"10, 11, 12"`
	MinPrice         *int64             `json:"min_price,omitempty" example:"100"`
	MaxPrice         *int64             `json:"max_price,omitempty" example:"500"`
	Rating           *int64             `json:"rating,omitempty" example:"1"`
	Characteristics  map[string][]int64 `json:"characteristics,omitempty"`
	Sort             *string            `json:"sort,omitempty" example:"price_asc"`
	Page             *int64             `json:"page,omitempty" example:"1"`
}

type FacetCount struct {
	ID    int64
	Count int64
}

type FilterFacets struct {
	Count            int64
	BasePriceMin     int64
	BasePriceMax     int64
	PriceMin         int64
	PriceMax         int64
	RatingCount      int64
	Brands           []FacetCount
	Subsubcategories []FacetCount
	Characteristics  map[string][]FacetCount
}
