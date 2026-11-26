package products

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

func parseint(value string) (int, error) {
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q: %w", value, err)
	}

	return parsed, nil
}

func parseCharacteristics(query url.Values) (map[string][]int, error) {
	characteristics := make(map[string][]int)

	const prefix = "characteristics["

	for key, values := range query {
		if !strings.HasPrefix(key, prefix) ||
			!strings.HasSuffix(key, "]") {
			continue
		}

		characteristicID := strings.TrimSuffix(
			strings.TrimPrefix(key, prefix),
			"]",
		)

		if characteristicID == "" {
			return nil, fmt.Errorf("empty characteristic ID")
		}

		for _, value := range values {
			id, err := parseint(value)
			if err != nil {
				return nil, fmt.Errorf(
					"characteristic %s: %w",
					characteristicID,
					err,
				)
			}

			characteristics[characteristicID] = append(
				characteristics[characteristicID],
				id,
			)
		}
	}

	return characteristics, nil
}

func parseCreateFilter(query url.Values) (CreateFilter, error) {
	sortVal := query.Get("sort")

	filter := CreateFilter{
		Sort: &sortVal,
	}

	var err error

	for _, value := range query["brands"] {
		id, err := parseint(value)
		if err != nil {
			return filter, fmt.Errorf("brands: %w", err)
		}

		filter.Brands = append(filter.Brands, id)
	}

	if value := query.Get("subcategory"); value != "" {
		filter.Subcategory, err = parseint(value)
		if err != nil {
			return filter, fmt.Errorf("subcategory: %w", err)
		}
	}

	if value := query.Get("subsubcategory"); value != "" {
		subsubcategory, err := parseint(value)
		if err != nil {
			return filter, fmt.Errorf("subsubcategory: %w", err)
		}
		filter.Subsubcategory = &subsubcategory
	}

	for _, value := range query["subsubcategories"] {
		id, err := parseint(value)
		if err != nil {
			return filter, fmt.Errorf("subsubcategories: %w", err)
		}

		filter.Subsubcategories = append(filter.Subsubcategories, id)
	}

	if value := query.Get("min_price"); value != "" {
		minPrice, err := parseint(value)
		if err != nil {
			return filter, fmt.Errorf("min_price: %w", err)
		}

		filter.MinPrice = &minPrice
	}

	if value := query.Get("max_price"); value != "" {
		maxPrice, err := parseint(value)
		if err != nil {
			return filter, fmt.Errorf("max_price: %w", err)
		}

		filter.MaxPrice = &maxPrice
	}

	if value := query.Get("rating"); value != "" {
		rating, err := parseint(value)
		if err != nil {
			return filter, fmt.Errorf("rating: %w", err)
		}

		filter.Rating = &rating
	}

	filter.Characteristics, err = parseCharacteristics(query)
	if err != nil {
		return filter, fmt.Errorf("characteristics: %w", err)
	}

	if value := query.Get("page"); value != "" {
		page, err := parseint(value)
		if err != nil {
			return filter, fmt.Errorf("page: %w", err)
		}

		filter.Page = &page
	}

	return filter, nil
}

func findFacetCount(facets []FacetCount, id int) int {
	for _, facet := range facets {
		if facet.ID == id {
			return facet.Count
		}
	}

	return 0
}

func addBasePriceArg(args *[]any, value any) string {
	*args = append(*args, value)
	return fmt.Sprintf("$%d", len(*args))
}
