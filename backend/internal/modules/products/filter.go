package products

import (
	"fmt"
	"strings"
)

type sqlFilter struct {
	where string
	args  []any
}

type PriceRange struct {
	Min float64
	Max float64
}

type excludeFilter struct {
	brands           bool
	subsubcategories bool
	characteristics  string
	price            bool
	rating           bool
}

func buildSQLFilter(filter CreateFilter) sqlFilter {
	var conditions []string
	var args []any

	add := func(condition string, value any) {
		args = append(args, value)
		conditions = append(
			conditions,
			fmt.Sprintf(condition, len(args)),
		)
	}

	conditions = append(conditions, "subcategory_id = $1")
	args = append(args, filter.Subcategory)

	if filter.Subsubcategory != nil {
		add("subsubcategory_id = $%d", *filter.Subsubcategory)
	}

	if len(filter.Brands) > 0 {
		add("brand_id = ANY($%d::bigint[])", filter.Brands)
	}

	if len(filter.Subsubcategories) > 0 {
		add("subsubcategory_id = ANY($%d::bigint[])", filter.Subsubcategories)
	}

	if filter.MinPrice != nil {
		add("final_price >= $%d", *filter.MinPrice)
	}

	if filter.MaxPrice != nil {
		add("final_price <= $%d", *filter.MaxPrice)
	}

	if filter.Rating != nil && *filter.Rating == 1 {
		conditions = append(conditions, "rating >= 4")
	}

	for key, values := range filter.Characteristics {
		if len(values) == 0 {
			continue
		}

		args = append(args, key+"_id")
		keyPlaceholder := len(args)

		args = append(args, values)
		valuesPlaceholder := len(args)

		conditions = append(
			conditions,
			fmt.Sprintf("(chars_general ->> $%d)::bigint = ANY($%d::bigint[])", keyPlaceholder, valuesPlaceholder),
		)
	}

	return sqlFilter{
		where: strings.Join(conditions, " AND "),
		args:  args,
	}
}

func appendFilterConditions(filter CreateFilter, args *[]any, exclude excludeFilter) []string {
	var conditions []string

	add := func(condition string, value any) {
		*args = append(*args, value)
		conditions = append(
			conditions,
			fmt.Sprintf(condition, len(*args)),
		)
	}

	add("p.subcategory_id = $%d", filter.Subcategory)

	if filter.Subsubcategory != nil &&
		!exclude.subsubcategories {
		add("p.subsubcategory_id = $%d", *filter.Subsubcategory)
	}

	if len(filter.Brands) > 0 && !exclude.brands {
		add("p.brand_id = ANY($%d::bigint[])", filter.Brands)
	}

	if len(filter.Subsubcategories) > 0 &&
		!exclude.subsubcategories {
		add(
			"p.subsubcategory_id = ANY($%d::bigint[])",
			filter.Subsubcategories,
		)
	}

	if !exclude.price {
		if filter.MinPrice != nil {
			add("p.final_price >= $%d", *filter.MinPrice)
		}

		if filter.MaxPrice != nil {
			add("p.final_price <= $%d", *filter.MaxPrice)
		}
	}

	if !exclude.rating &&
		filter.Rating != nil &&
		*filter.Rating == 1 {
		conditions = append(conditions, "p.rating >= 4")
	}

	for key, values := range filter.Characteristics {
		if len(values) == 0 {
			continue
		}

		if key == exclude.characteristics {
			continue
		}

		jsonKey := key + "_id"

		*args = append(*args, jsonKey)
		keyPos := len(*args)

		*args = append(*args, values)
		valuesPos := len(*args)

		conditions = append(
			conditions,
			fmt.Sprintf(
				"(p.chars_general ->> $%d)::bigint = ANY($%d::bigint[])",
				keyPos,
				valuesPos,
			),
		)
	}

	return conditions
}

func whereClause(conditions []string) string {
	if len(conditions) == 0 {
		return "TRUE"
	}

	return strings.Join(conditions, " AND ")
}
