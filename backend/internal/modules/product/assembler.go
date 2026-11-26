package product

import (
	"ArtiVana/internal/modules/categories"
	"ArtiVana/internal/modules/dictionaries"
	"fmt"
	"strings"
)

type ProductAssembler struct{}

func NewProductAssembler() *ProductAssembler {
	return &ProductAssembler{}
}

func (a *ProductAssembler) BuildBrand(product Product, brands []dictionaries.ProductOption) (dictionaries.ProductOption, error) {
	for _, brand := range brands {
		if brand.ID == product.BrandID {
			return brand, nil
		}
	}

	return dictionaries.ProductOption{}, fmt.Errorf(
		"brand %d not found",
		product.BrandID,
	)
}

func (a *ProductAssembler) BuildGroup(
	groupsProduct []ProductByGroupID,
	productBrand dictionaries.ProductOption,
	id int64,
) []ProductsGroupItem {
	productsGroup := make([]ProductsGroupItem, 0)

	for _, p := range groupsProduct {
		if p.ID == id {
			continue
		}

		model := ""
		if p.Model != nil {
			model = *p.Model
		}

		image := ""
		if len(p.Images) > 0 {
			image = p.Images[0]
		}

		productsGroup = append(
			productsGroup,
			ProductsGroupItem{
				Slug:  p.Slug,
				Name:  strings.TrimSpace(productBrand.Name + " " + model),
				Image: image,
			},
		)
	}

	return productsGroup
}

func (a *ProductAssembler) BuildCategory(
	product Product,
	categoriesAll []categories.Category,
	subcategoriesAll []categories.Subcategory,
	subsubcategoriesAll []categories.Subsubcategory,
) (
	ProductOptionWithSlug,
	ProductOptionWithSlug,
	ProductOptionWithSlug,
	error,
) {
	var productCategory, productSubcategory, productSubsubcategory ProductOptionWithSlug

	categoryFound := false
	for _, category := range categoriesAll {
		if category.ID == product.CategoryID {
			productCategory = ProductOptionWithSlug{
				ID:   category.ID,
				Name: category.Name,
				Slug: category.Slug,
			}
			categoryFound = true
			break
		}
	}

	if !categoryFound {
		return productCategory, productSubcategory, productSubsubcategory,
			fmt.Errorf("category %d not found", product.CategoryID)
	}

	subcategoryFound := false
	for _, subCategory := range subcategoriesAll {
		if subCategory.ID == product.SubcategoryID {
			productSubcategory = ProductOptionWithSlug{
				ID:   subCategory.ID,
				Name: subCategory.Name,
				Slug: subCategory.Slug,
			}
			subcategoryFound = true
			break
		}
	}

	if !subcategoryFound {
		return productCategory, productSubcategory, productSubsubcategory,
			fmt.Errorf("subcategory %d not found", product.SubcategoryID)
	}

	subsubcategoryFound := false
	for _, subsubCategory := range subsubcategoriesAll {
		if subsubCategory.ID == product.SubsubcategoryID {
			productSubsubcategory = ProductOptionWithSlug{
				ID:   subsubCategory.ID,
				Name: subsubCategory.Name,
				Slug: subsubCategory.Slug,
			}
			subsubcategoryFound = true
			break
		}
	}

	if !subsubcategoryFound {
		return productCategory, productSubcategory, productSubsubcategory,
			fmt.Errorf("subsubcategory %d not found", product.SubsubcategoryID)
	}
	return productCategory, productSubcategory, productSubsubcategory, nil
}

func (a *ProductAssembler) BuildCharacteristics(
	productChars map[string]int64,
	chars []dictionaries.ProductOptionChars,
	charsOptions []dictionaries.ProductOptionCharsOption,
) []ProductCharacteristic {
	charsMap := make(
		map[string]dictionaries.ProductOptionChars,
		len(chars),
	)
	for _, char := range chars {
		charsMap[char.Value] = char
	}

	charsOptionsMap := make(
		map[int64]dictionaries.ProductOptionCharsOption,
		len(charsOptions),
	)
	for _, option := range charsOptions {
		charsOptionsMap[option.ID] = option
	}

	result := make([]ProductCharacteristic, 0, len(productChars))

	for key, value := range productChars {
		char, ok := charsMap[strings.TrimSuffix(key, "_id")]
		if !ok {
			continue
		}

		option, ok := charsOptionsMap[value]
		if !ok {
			continue
		}

		result = append(result, ProductCharacteristic{
			ID:    char.ID,
			Name:  char.Name,
			Value: char.Value,
			Option: dictionaries.ProductOption{
				ID:   option.ID,
				Name: option.Name,
			},
		})
	}

	return result
}

func (a *ProductAssembler) BuildDetail(
	product Product,
	productBrand dictionaries.ProductOption,
	productsGroup []ProductsGroupItem,
	productCategory ProductOptionWithSlug,
	productSubcategory ProductOptionWithSlug,
	productSubsubcategory ProductOptionWithSlug,
	productCharsGeneral []ProductCharacteristic,
	productCharsExtra []ProductCharacteristic,
) ProductResponse {
	response := ProductResponse{
		ID:              product.ID,
		Slug:            product.Slug,
		GroupID:         product.GroupID,
		ProductsGroup:   productsGroup,
		CreatedAt:       product.CreatedAt,
		Name:            product.Name,
		Article:         product.Article,
		Model:           product.Model,
		OriginalPrice:   product.OriginalPrice,
		FinalPrice:      product.FinalPrice,
		Stock:           product.Stock,
		Images:          product.Images,
		Rating:          product.Rating,
		Package:         product.Package,
		Description:     product.Description,
		ShippingOptions: product.ShippingOptions,
		Tags:            product.Tags,
		Brand:           productBrand,
		Category:        productCategory,
		Subcategory:     productSubcategory,
		Subsubcategory:  productSubsubcategory,
		CharsGeneral:    productCharsGeneral,
		CharsExtra:      productCharsExtra,
	}

	return response
}

func (a *ProductAssembler) BuildProductCutDown(
	product Product,
	productCategory ProductOptionWithSlug,
	productSubcategory ProductOptionWithSlug,
	productSubsubcategory ProductOptionWithSlug,
) ProductCutDown {
	response := ProductCutDown{
		ID:             product.ID,
		Slug:           product.Slug,
		Name:           product.Name,
		Article:        product.Article,
		OriginalPrice:  product.OriginalPrice,
		FinalPrice:     product.FinalPrice,
		Stock:          product.Stock,
		Images:         product.Images,
		Tags:           product.Tags,
		Category:       productCategory,
		Subcategory:    productSubcategory,
		Subsubcategory: productSubsubcategory,
	}

	return response
}
