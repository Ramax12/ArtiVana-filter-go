package product

import (
	"ArtiVana/internal/modules/dictionaries"
	"time"
)

type DeliveryOption string

const (
	deliveryOptionShipping DeliveryOption = "shipping"
	deliveryOptionPickup   DeliveryOption = "pickup"
	deliveryOptionDelivery DeliveryOption = "delivery"
)

type ProductOptionWithSlug struct {
	dictionaries.ProductOption
	Slug string `json:"slug" example:"acrylic"`
}

type ProductCharacteristic struct {
	dictionaries.ProductOptionChars
	Option dictionaries.ProductOption `json:"option"`
}

type ProductsGroupItem struct {
	Slug  string `json:"slug" example:"neewer-bi-color-660-led-video-light-142816"`
	Name  string `json:"name" example:"Neewer Bi-Color 660 LED Video Light"`
	Image string `json:"image" example:"/img/image1.jpg"`
}

type ProductByGroupID struct {
	ID      int    `json:"id" example:"1"`
	Slug    string   `json:"slug" example:"neewer-bi-color-660-led-video-light-142816"`
	Name    string   `json:"name" example:"Neewer Bi-Color 660 LED Video Light"`
	Model   *string  `json:"model" example:"K2-90" extensions:"x-nullable"`
	Images  []string `json:"images" example:"/img/image1.jpg, /img/image2.jpg"`
	GroupID int    `json:"group_id" example:"1"`
}

type Product struct {
	ID               int
	Slug             string
	GroupID          *int
	CreatedAt        time.Time
	CategoryID       int
	SubcategoryID    int
	SubsubcategoryID int
	Name             string
	Article          int
	BrandID          int
	Model            *string
	OriginalPrice    *int
	FinalPrice       int
	Stock            int
	Images           []string
	Rating           float64
	CharsGeneral     map[string]int
	CharsExtra       *map[string]int
	Package          []string
	Description      string
	ShippingOptions  []DeliveryOption
	Tags             []string
}

type ProductCutDown struct {
	ID             int                 `json:"id" example:"1"`
	Slug           string                `json:"slug" example:"neewer-bi-color-660-led-video-light-142816"`
	Category       ProductOptionWithSlug `json:"category"`
	Subcategory    ProductOptionWithSlug `json:"subcategory"`
	Subsubcategory ProductOptionWithSlug `json:"subsubcategory"`
	Name           string                `json:"name" example:"Neewer Bi-Color 660 LED Video Light"`
	Article        int                 `json:"article" example:"123456"`
	OriginalPrice  *int                `json:"original_price" example:"25000" extensions:"x-nullable"`
	FinalPrice     int                 `json:"final_price" example:"24000"`
	Stock          int                 `json:"stock" example:"50"`
	Images         []string              `json:"images" example:"/img/image1.jpg, /img/image2.jpg"`
	Tags           []string              `json:"tags" example:"New Arrival, Best Seller" extensions:"x-nullable"`
}

type ProductResponse struct {
	ID              int                      `json:"id" example:"1"`
	Slug            string                     `json:"slug" example:"neewer-bi-color-660-led-video-light-142816"`
	GroupID         *int                     `json:"group_id" example:"1" extensions:"x-nullable"`
	ProductsGroup   []ProductsGroupItem        `json:"products_group" extension:"x-nullable"`
	CreatedAt       time.Time                  `json:"created_at" example:"2025-10-11T14:30:00.000Z"`
	Category        ProductOptionWithSlug      `json:"category"`
	Subcategory     ProductOptionWithSlug      `json:"subcategory"`
	Subsubcategory  ProductOptionWithSlug      `json:"subsubcategory"`
	Name            string                     `json:"name" example:"Neewer Bi-Color 660 LED Video Light"`
	Article         int                      `json:"article" example:"123456"`
	Brand           dictionaries.ProductOption `json:"brand"`
	Model           *string                    `json:"model" example:"K2-90" extensions:"x-nullable"`
	OriginalPrice   *int                     `json:"original_price" example:"24000" extensions:"x-nullable"`
	FinalPrice      int                      `json:"final_price" example:"24000"`
	Stock           int                      `json:"stock" example:"50"`
	Images          []string                   `json:"images" example:"/img/image1.jpg, /img/image2.jpg"`
	Rating          float64                    `json:"rating" example:"4.8"`
	CharsGeneral    []ProductCharacteristic    `json:"chars_general"`
	CharsExtra      []ProductCharacteristic    `json:"chars_extra" extensions:"x-nullable"`
	Package         []string                   `json:"package" example:"2 x Bi-color LED Video" extensions:"x-nullable"`
	Description     string                     `json:"description" example:"Bi Color Dimmable LED Studio Light with 330 warm white and 330 cool white LEDs"`
	ShippingOptions []DeliveryOption           `json:"shipping_options" enums:"shipping,pickup,delivery" example:"shipping, pickup"`
	Tags            []string                   `json:"tags" example:"Budget Pick" extensions:"x-nullable"`
}
