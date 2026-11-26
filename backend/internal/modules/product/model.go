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
	ID      int64    `json:"id" example:"1"`
	Slug    string   `json:"slug" example:"neewer-bi-color-660-led-video-light-142816"`
	Name    string   `json:"name" example:"Neewer Bi-Color 660 LED Video Light"`
	Model   *string  `json:"model" example:"K2-90" extensions:"x-nullable"`
	Images  []string `json:"images" example:"/img/image1.jpg, /img/image2.jpg"`
	GroupID int64    `json:"group_id" example:"1"`
}

type Product struct {
	ID               int64
	Slug             string
	GroupID          *int64
	CreatedAt        time.Time
	CategoryID       int64
	SubcategoryID    int64
	SubsubcategoryID int64
	Name             string
	Article          int64
	BrandID          int64
	Model            *string
	OriginalPrice    *int64
	FinalPrice       int64
	Stock            int64
	Images           []string
	Rating           float64
	CharsGeneral     map[string]int64
	CharsExtra       *map[string]int64
	Package          []string
	Description      string
	ShippingOptions  []DeliveryOption
	Tags             []string
}

type ProductCutDown struct {
	ID             int64                 `json:"id" example:"1"`
	Slug           string                `json:"slug" example:"neewer-bi-color-660-led-video-light-142816"`
	Category       ProductOptionWithSlug `json:"category"`
	Subcategory    ProductOptionWithSlug `json:"subcategory"`
	Subsubcategory ProductOptionWithSlug `json:"subsubcategory"`
	Name           string                `json:"name" example:"Neewer Bi-Color 660 LED Video Light"`
	Article        int64                 `json:"article" example:"123456"`
	OriginalPrice  *int64                `json:"original_price" example:"25000" extensions:"x-nullable"`
	FinalPrice     int64                 `json:"final_price" example:"24000"`
	Stock          int64                 `json:"stock" example:"50"`
	Images         []string              `json:"images" example:"/img/image1.jpg, /img/image2.jpg"`
	Tags           []string              `json:"tags" example:"New Arrival, Best Seller" extensions:"x-nullable"`
}

type ProductResponse struct {
	ID              int64                      `json:"id" example:"1"`
	Slug            string                     `json:"slug" example:"neewer-bi-color-660-led-video-light-142816"`
	GroupID         *int64                     `json:"group_id" example:"1" extensions:"x-nullable"`
	ProductsGroup   []ProductsGroupItem        `json:"products_group" extension:"x-nullable"`
	CreatedAt       time.Time                  `json:"created_at" example:"2025-10-11T14:30:00.000Z"`
	Category        ProductOptionWithSlug      `json:"category"`
	Subcategory     ProductOptionWithSlug      `json:"subcategory"`
	Subsubcategory  ProductOptionWithSlug      `json:"subsubcategory"`
	Name            string                     `json:"name" example:"Neewer Bi-Color 660 LED Video Light"`
	Article         int64                      `json:"article" example:"123456"`
	Brand           dictionaries.ProductOption `json:"brand"`
	Model           *string                    `json:"model" example:"K2-90" extensions:"x-nullable"`
	OriginalPrice   *int64                     `json:"original_price" example:"24000" extensions:"x-nullable"`
	FinalPrice      int64                      `json:"final_price" example:"24000"`
	Stock           int64                      `json:"stock" example:"50"`
	Images          []string                   `json:"images" example:"/img/image1.jpg, /img/image2.jpg"`
	Rating          float64                    `json:"rating" example:"4.8"`
	CharsGeneral    []ProductCharacteristic    `json:"chars_general"`
	CharsExtra      []ProductCharacteristic    `json:"chars_extra" extensions:"x-nullable"`
	Package         []string                   `json:"package" example:"2 x Bi-color LED Video" extensions:"x-nullable"`
	Description     string                     `json:"description" example:"Bi Color Dimmable LED Studio Light with 330 warm white and 330 cool white LEDs"`
	ShippingOptions []DeliveryOption           `json:"shipping_options" enums:"shipping,pickup,delivery" example:"shipping, pickup"`
	Tags            []string                   `json:"tags" example:"Budget Pick" extensions:"x-nullable"`
}
