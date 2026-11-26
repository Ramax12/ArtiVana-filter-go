package categories

type Category struct {
	ID    int64  `json:"id"`
	Slug  string `json:"slug"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

type Subcategory struct {
	ID         int64  `json:"id"`
	Slug       string `json:"slug"`
	Name       string `json:"name"`
	Image      string `json:"image"`
	CategoryID int64  `json:"category_id"`
}

type Subsubcategory struct {
	ID            int64  `json:"id"`
	Slug          string `json:"slug"`
	Name          string `json:"name"`
	Image         string `json:"image"`
	SubcategoryID int64  `json:"subcategory_id"`
}

type CategoryResponse struct {
	ID            int64                 `json:"id" example:"1"`
	Slug          string                `json:"slug" example:"oil-paints"`
	Name          string                `json:"name" example:"Oil Paints"`
	Image         string                `json:"image" example:"/img/image.jpg"`
	Subcategories []SubcategoryResponse `json:"subcategories"`
}

type SubcategoryResponse struct {
	ID               int64                    `json:"id" example:"1"`
	Slug             string                   `json:"slug" example:"paints"`
	Name             string                   `json:"name" example:"Paints"`
	Image            string                   `json:"image" example:"/img/image.jpg"`
	Subsubcategories []SubsubcategoryResponse `json:"subsubcategories"`
	CategoryID       int64                    `json:"category_id" example:"1"`
}

type SubsubcategoryResponse struct {
	ID            int64  `json:"id" example:"1"`
	Slug          string `json:"slug" example:"art-supplies"`
	Name          string `json:"name" example:"Art Supplies"`
	Image         string `json:"image" example:"/img/image.jpg"`
	SubcategoryID int64  `json:"subcategory_id" example:"1"`
}
