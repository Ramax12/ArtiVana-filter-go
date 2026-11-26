package products

import (
	"ArtiVana/internal/modules/categories"
	"ArtiVana/internal/modules/dictionaries"
	productModel "ArtiVana/internal/modules/product"
	"context"
	"sort"
)

type Service struct {
	productRepository      productModel.Repository
	productsRepository     Repository
	dictionariesRepository dictionaries.Repository
	categoriesRepository   categories.Repository
	assembler              *productModel.ProductAssembler
}

func NewService(
	productRepository productModel.Repository,
	productsRepository Repository,
	dictionariesRepository dictionaries.Repository,
	categoriesRepository categories.Repository,
) *Service {
	return &Service{
		productRepository:      productRepository,
		productsRepository:     productsRepository,
		dictionariesRepository: dictionariesRepository,
		categoriesRepository:   categoriesRepository,
		assembler:              productModel.NewProductAssembler(),
	}
}

func (s *Service) GetProducts(
	ctx context.Context,
	ids []int64,
	subcategoryID int64,
	subsubcategoryID int64,
) (
	[]productModel.ProductResponse,
	error,
) {
	var products []productModel.ProductResponse

	categoriesAll, subcategoriesAll, subsubcategoriesAll, err := s.categoriesRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		product, err := s.productRepository.GetByID(ctx, id)
		if err != nil {
			return nil, err
		}

		if subcategoryID&product.SubcategoryID != subcategoryID {
			return nil, nil
		}
		if subsubcategoryID&product.SubsubcategoryID != subsubcategoryID {
			return nil, nil
		}

		brands, err := s.dictionariesRepository.FindAllBrands(ctx)
		if err != nil {
			return nil, err
		}
		productBrand, err := s.assembler.BuildBrand(*product, brands)
		if err != nil {
			return nil, err
		}

		productsGroup := make([]productModel.ProductsGroupItem, 0)
		if product.GroupID != nil {
			groupsProduct, err := s.productRepository.FindByGroupID(
				ctx,
				*product.GroupID,
				id,
			)
			if err != nil {
				return nil, err
			}

			productsGroup = s.assembler.BuildGroup(groupsProduct, productBrand, id)
		}

		productCategory, productSubcategory, productSubsubcategory, err := s.assembler.BuildCategory(
			*product,
			categoriesAll,
			subcategoriesAll,
			subsubcategoriesAll,
		)
		if err != nil {
			return nil, err
		}

		charsGeneral, err := s.dictionariesRepository.FindAllCharsGeneral(ctx)
		if err != nil {
			return nil, err
		}
		charsGeneralOptions, err := s.dictionariesRepository.FindAllCharsGeneralOptions(ctx)
		if err != nil {
			return nil, err
		}
		productCharsGeneral := s.assembler.BuildCharacteristics(product.CharsGeneral, charsGeneral, charsGeneralOptions)

		charsExtra, err := s.dictionariesRepository.FindAllCharsExtra(ctx)
		if err != nil {
			return nil, err
		}
		charsExtraOptions, err := s.dictionariesRepository.FindAllCharsExtraOptions(ctx)
		if err != nil {
			return nil, err
		}

		var productCharsExtra []productModel.ProductCharacteristic
		if product.CharsExtra != nil {
			productCharsExtra = s.assembler.BuildCharacteristics(*product.CharsExtra, charsExtra, charsExtraOptions)
		}

		products = append(
			products,
			s.assembler.BuildDetail(
				*product,
				productBrand,
				productsGroup,
				productCategory,
				productSubcategory,
				productSubsubcategory,
				productCharsGeneral,
				productCharsExtra,
			),
		)
	}

	return products, nil
}

func (s *Service) GetFilterProducts(ctx context.Context, filter CreateFilter) ([]productModel.ProductCutDown, error) {
	result := []productModel.ProductCutDown{}

	products, err := s.productsRepository.FindFiltered(ctx, filter)
	if err != nil {
		return nil, err
	}

	categoriesAll, subcategoriesAll, subsubcategoriesAll, err := s.categoriesRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, product := range products {
		productCategory, productSubcategory, productSubsubcategory, err := s.assembler.BuildCategory(
			product,
			categoriesAll,
			subcategoriesAll,
			subsubcategoriesAll,
		)
		if err != nil {
			return nil, err
		}

		result = append(
			result,
			s.assembler.BuildProductCutDown(
				product,
				productCategory,
				productSubcategory,
				productSubsubcategory,
			),
		)
	}

	return result, nil
}

func (s *Service) buildFilterData() FilterData {
	return FilterData{
		Standard: FilterDataStandard{
			Subsubcategories: make([]FilterDataItem, 0),
			Brands:           make([]FilterDataItem, 0),
			Price: FilterDataPrice{
				BaseRange: [2]int64{},
				Range:     [2]int64{},
			},
			Rating: FilterDataRating{Count: 0},
		},
		Characteristics: make([]FilterDataItemCharacteristics, 0),
	}
}

func (s *Service) GetFilterMeta(ctx context.Context, filter CreateFilter) (FilterMeta, error) {
	facets, err := s.productsRepository.GetFilterFacets(ctx, filter)
	if err != nil {
		return FilterMeta{}, err
	}

	subsubcategories, err := s.productsRepository.FindUsedSubsubcategories(
		ctx,
		filter.Subcategory,
	)
	if err != nil {
		return FilterMeta{}, err
	}

	brands, err := s.productsRepository.FindUsedBrands(
		ctx,
		filter.Subcategory,
		filter.Subsubcategory,
	)
	if err != nil {
		return FilterMeta{}, err
	}

	characteristics, err := s.productsRepository.FindUsedCharsGeneral(
		ctx,
		filter.Subcategory,
		filter.Subsubcategory,
	)
	if err != nil {
		return FilterMeta{}, err
	}

	characteristicOptions, err := s.dictionariesRepository.FindAllCharsGeneralOptions(ctx)
	if err != nil {
		return FilterMeta{}, err
	}

	usedOptions, err := s.productsRepository.FindUsedCharsGeneralOptionIDs(
		ctx,
		filter.Subcategory,
		filter.Subsubcategory,
	)
	if err != nil {
		return FilterMeta{}, err
	}

	filterData := s.buildFilterData()

	for _, subsubcategory := range subsubcategories {
		filterData.Standard.Subsubcategories = append(
			filterData.Standard.Subsubcategories,
			FilterDataItem{
				ID:    subsubcategory.ID,
				Name:  subsubcategory.Name,
				Count: findFacetCount(facets.Subsubcategories, subsubcategory.ID),
			},
		)
	}

	sort.SliceStable(filterData.Standard.Subsubcategories, func(i, j int) bool {
		return filterData.Standard.Subsubcategories[i].Count >
			filterData.Standard.Subsubcategories[j].Count
	})

	for _, brand := range brands {
		filterData.Standard.Brands = append(
			filterData.Standard.Brands,
			FilterDataItem{
				ProductOption: brand,
				Count:         findFacetCount(facets.Brands, brand.ID),
			},
		)
	}

	sort.SliceStable(filterData.Standard.Brands, func(i, j int) bool {
		return filterData.Standard.Brands[i].Count >
			filterData.Standard.Brands[j].Count
	})

	for _, characteristic := range characteristics {
		item := FilterDataItemCharacteristics{
			ProductOptionChars: characteristic,
			Options:            make([]FilterDataItem, 0),
		}

		availableIDs := usedOptions[characteristic.Value]

		for _, option := range characteristicOptions {
			if option.ParentID != characteristic.ID {
				continue
			}

			if _, exists := availableIDs[option.ID]; !exists {
				continue
			}

			item.Options = append(item.Options, FilterDataItem{
				ProductOption: option.ProductOption,
				Count: findFacetCount(
					facets.Characteristics[characteristic.Value],
					option.ID,
				),
			})
		}

		sort.SliceStable(item.Options, func(i, j int) bool {
			return item.Options[i].Count > item.Options[j].Count
		})

		filterData.Characteristics = append(
			filterData.Characteristics,
			item,
		)
	}

	filterData.Standard.Price = FilterDataPrice{
		BaseRange: [2]int64{
			facets.BasePriceMin,
			facets.BasePriceMax,
		},
		Range: [2]int64{
			facets.PriceMin,
			facets.PriceMax,
		},
	}

	filterData.Standard.Rating.Count = facets.RatingCount

	return FilterMeta{
		Count:   facets.Count,
		Filters: filterData,
	}, nil
}
