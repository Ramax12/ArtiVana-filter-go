package products

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"ArtiVana/internal/modules/categories"
	"ArtiVana/internal/modules/dictionaries"
	productModel "ArtiVana/internal/modules/product"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type facetRow struct {
	Kind  string
	Key   string
	ID    int64
	Count int64
}

type Repository interface {
	FindAll(ctx context.Context) ([]productModel.Product, error)
	FindFiltered(ctx context.Context, filter CreateFilter) ([]productModel.Product, error)
	FindUsedSubsubcategories(ctx context.Context, subcategoryID int64) ([]categories.Subsubcategory, error)
	FindUsedBrands(ctx context.Context, subcategoryID int64, subsubcategoryID *int64) ([]dictionaries.ProductOption, error)
	FindUsedCharsGeneral(ctx context.Context, subcategoryID int64, subsubcategoryID *int64) ([]dictionaries.ProductOptionChars, error)
	FindUsedCharsGeneralOptionIDs(ctx context.Context, subcategoryID int64, subsubcategoryID *int64) (map[string]map[int64]struct{}, error)
	GetFilterFacets(ctx context.Context, filter CreateFilter) (FilterFacets, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{
		db: db,
	}
}

func buildOrderBy(sort string) string {
	switch sort {
	case "rating_desc":
		return " ORDER BY rating DESC, id ASC"
	case "price_asc":
		return " ORDER BY final_price ASC, id ASC"
	case "price_desc":
		return " ORDER BY final_price DESC, id ASC"
	case "name_asc":
		return " ORDER BY name ASC, id ASC"
	case "name_desc":
		return " ORDER BY name DESC, id ASC"
	default:
		return " ORDER BY id ASC"
	}
}

func scanProducts(rows pgx.Rows) ([]productModel.Product, error) {
	defer rows.Close()

	products := []productModel.Product{}

	for rows.Next() {
		var product productModel.Product
		var shippingOptions []string

		err := rows.Scan(
			&product.ID,
			&product.Slug,
			&product.GroupID,
			&product.CreatedAt,
			&product.CategoryID,
			&product.SubcategoryID,
			&product.SubsubcategoryID,
			&product.Name,
			&product.Article,
			&product.BrandID,
			&product.Model,
			&product.OriginalPrice,
			&product.FinalPrice,
			&product.Stock,
			&product.Images,
			&product.Rating,
			&product.CharsGeneral,
			&product.CharsExtra,
			&product.Package,
			&product.Description,
			&shippingOptions,
			&product.Tags,
		)
		if err != nil {
			return nil, err
		}

		product.ShippingOptions = make(
			[]productModel.DeliveryOption,
			len(shippingOptions),
		)

		for i, option := range shippingOptions {
			product.ShippingOptions[i] = productModel.DeliveryOption(option)
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

const productColumns = `
	id, slug, group_id, created_at, category_id, subcategory_id, subsubcategory_id, name, article, brand_id, model, original_price, 
	final_price, stock, images, rating, chars_general, chars_extra, package, description, shipping_options::text[], tags`

func (r *repository) FindAll(ctx context.Context) ([]productModel.Product, error) {
	rows, err := r.db.Query(ctx,
		fmt.Sprintf(`
			SELECT %s
			FROM products
		`, productColumns),
	)

	if err != nil {
		return nil, err
	}

	return scanProducts(rows)
}

func (r *repository) FindFiltered(ctx context.Context, filter CreateFilter) ([]productModel.Product, error) {
	sqlFilter := buildSQLFilter(filter)

	query := fmt.Sprintf(`
		SELECT %s
		FROM products
		WHERE %s
	`, productColumns, sqlFilter.where)

	if filter.Sort != nil {
		query += buildOrderBy(*filter.Sort)
	}

	const perPage = 24
	var page int64 = 1
	if filter.Page != nil {
		page = *filter.Page
	}

	sqlFilter.args = append(sqlFilter.args, perPage)
	limitPosition := len(sqlFilter.args)

	sqlFilter.args = append(sqlFilter.args, (page-1)*perPage)
	offsetPosition := len(sqlFilter.args)

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", limitPosition, offsetPosition)

	rows, err := r.db.Query(ctx, query, sqlFilter.args...)
	if err != nil {
		return nil, fmt.Errorf("query filtered products: %w", err)
	}

	return scanProducts(rows)
}

func (r *repository) FindUsedSubsubcategories(ctx context.Context, subcategoryID int64) ([]categories.Subsubcategory, error) {
	query := `
		SELECT DISTINCT
			s.id,
			s.slug,
			s.name,
			s.image,
			s.subcategory_id
		FROM products p
		JOIN subsubcategories s ON s.id = p.subsubcategory_id
		WHERE p.subcategory_id = $1
	`

	args := []any{subcategoryID}

	query += ` ORDER BY s.name`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("FindUsedSubsubcategories: failed to execute query: %w", err)
	}
	defer rows.Close()

	subsubcategories := make([]categories.Subsubcategory, 0)

	for rows.Next() {
		var subsubcategory categories.Subsubcategory

		if err := rows.Scan(
			&subsubcategory.ID,
			&subsubcategory.Slug,
			&subsubcategory.Name,
			&subsubcategory.Image,
			&subsubcategory.SubcategoryID,
		); err != nil {
			return nil, fmt.Errorf("FindUsedSubsubcategories: failed to scan subsubcategory: %w", err)
		}

		subsubcategories = append(subsubcategories, subsubcategory)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("FindUsedSubsubcategories: failed to iterate rows: %w", err)
	}

	return subsubcategories, nil
}

func (r *repository) FindUsedBrands(
	ctx context.Context,
	subcategoryID int64,
	subsubcategoryID *int64,
) ([]dictionaries.ProductOption, error) {
	query := `
		SELECT DISTINCT
			b.id,
			b.name
		FROM products p
		JOIN brands b ON b.id = p.brand_id
		WHERE p.subcategory_id = $1
	`

	args := []any{subcategoryID}

	if subsubcategoryID != nil {
		args = append(args, *subsubcategoryID)
		query += fmt.Sprintf(" AND p.subsubcategory_id = $%d", len(args))
	}

	query += ` ORDER BY b.name`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("FindUsedBrands: failed to execute query: %w", err)
	}
	defer rows.Close()

	brands := make([]dictionaries.ProductOption, 0)

	for rows.Next() {
		var brand dictionaries.ProductOption

		if err := rows.Scan(
			&brand.ID,
			&brand.Name,
		); err != nil {
			return nil, fmt.Errorf("FindUsedBrands: failed to scan brand: %w", err)
		}

		brands = append(brands, brand)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("FindUsedBrands: failed to iterate rows: %w", err)
	}

	return brands, nil
}

func (r *repository) FindUsedCharsGeneral(
	ctx context.Context,
	subcategoryID int64,
	subsubcategoryID *int64,
) ([]dictionaries.ProductOptionChars, error) {
	query := `
		SELECT DISTINCT
				c.id,
				c.name,
				c.value
		FROM products p
		CROSS JOIN LATERAL jsonb_object_keys(
				COALESCE(p.chars_general, '{}'::jsonb)
		) AS k(key)
		JOIN chars_general c
				ON k.key = c.value || '_id'
		WHERE p.subcategory_id = $1
	`

	args := []any{subcategoryID}

	if subsubcategoryID != nil {
		args = append(args, *subsubcategoryID)
		query += fmt.Sprintf(" AND p.subsubcategory_id = $%d", len(args))
	}

	query += ` ORDER BY c.id`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("FindUsedCharsGeneral: failed to execute query: %w", err)
	}
	defer rows.Close()

	chars := make([]dictionaries.ProductOptionChars, 0)

	for rows.Next() {
		var char dictionaries.ProductOptionChars

		if err := rows.Scan(
			&char.ID,
			&char.Name,
			&char.Value,
		); err != nil {
			return nil, fmt.Errorf("FindUsedCharsGeneral: failed to scan characteristic: %w", err)
		}

		chars = append(chars, char)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("FindUsedCharsGeneral: failed to iterate rows: %w", err)
	}

	return chars, nil
}

func (r *repository) FindUsedCharsGeneralOptionIDs(
	ctx context.Context,
	subcategoryID int64,
	subsubcategoryID *int64,
) (map[string]map[int64]struct{}, error) {
	query := `
		SELECT DISTINCT
			j.key,
			j.value::bigint AS option_id
		FROM products p
		CROSS JOIN LATERAL jsonb_each_text(
			COALESCE(p.chars_general, '{}'::jsonb)
		) AS j(key, value)
		WHERE p.subcategory_id = $1
			AND j.key LIKE '%\_id' ESCAPE '\'
	`

	args := []any{subcategoryID}

	if subsubcategoryID != nil {
		args = append(args, *subsubcategoryID)
		query += fmt.Sprintf(
			" AND p.subsubcategory_id = $%d",
			len(args),
		)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf(
			"FindUsedCharsGeneralOptionIDs: failed to execute query: %w",
			err,
		)
	}
	defer rows.Close()

	result := make(map[string]map[int64]struct{})

	for rows.Next() {
		var key string
		var optionID int64

		if err := rows.Scan(&key, &optionID); err != nil {
			return nil, fmt.Errorf(
				"FindUsedCharsGeneralOptionIDs: failed to scan option: %w",
				err,
			)
		}

		key = strings.TrimSuffix(key, "_id")

		if result[key] == nil {
			result[key] = make(map[int64]struct{})
		}

		result[key][optionID] = struct{}{}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"FindUsedCharsGeneralOptionIDs: failed to iterate rows: %w",
			err,
		)
	}

	return result, nil
}

var placeholderRegexp = regexp.MustCompile(`\$(\d+)`)

func shiftPlaceholders(query string, offset int) string {
	return placeholderRegexp.ReplaceAllStringFunc(query, func(match string) string {
		n, _ := strconv.Atoi(match[1:])
		return fmt.Sprintf("$%d", n+offset)
	})
}

func (r *repository) GetFilterFacets(ctx context.Context, filter CreateFilter) (FilterFacets, error) {
	var result FilterFacets

	result.Characteristics = make(map[string][]FacetCount)

	var branches []string
	var args []any

	addBranch := func(branch string, branchArgs []any) {
		branch = shiftPlaceholders(branch, len(args))
		branches = append(branches, branch)
		args = append(args, branchArgs...)
	}

	var statsArgs []any

	statsConditions := appendFilterConditions(filter, &statsArgs, excludeFilter{})

	statsBranch := `
    SELECT
			'stats'::text AS kind,
			''::text AS facet_key,
			0::bigint AS option_id,
			COUNT(*)::bigint AS item_count
    FROM products p
    WHERE ` + whereClause(statsConditions)

	addBranch(statsBranch, statsArgs)

	var brandArgs []any

	brandConditions := appendFilterConditions(filter, &brandArgs, excludeFilter{brands: true})

	brandBranch := `
    SELECT
			'brand'::text,
			''::text,
			p.brand_id::bigint,
			COUNT(*)::bigint
    FROM products p
    WHERE ` + whereClause(brandConditions) + `
    GROUP BY p.brand_id
	`

	addBranch(brandBranch, brandArgs)

	var subcatArgs []any

	subcatConditions := appendFilterConditions(filter, &subcatArgs, excludeFilter{subsubcategories: true})

	subcatBranch := `
    SELECT
			'subsubcategory'::text,
			''::text,
			p.subsubcategory_id::bigint,
			COUNT(*)::bigint
    FROM products p
    WHERE ` + whereClause(subcatConditions) + `
    GROUP BY p.subsubcategory_id
	`

	addBranch(subcatBranch, subcatArgs)

	var priceArgs []any

	priceConditions := appendFilterConditions(filter, &priceArgs, excludeFilter{})

	priceBranch := `
		SELECT
			'price_min'::text,
			''::text,
			COALESCE(MIN(p.final_price), 0)::bigint,
			0::bigint
		FROM products p
		WHERE ` + whereClause(priceConditions) + `

		UNION ALL

		SELECT
			'price_max'::text,
			''::text,
			COALESCE(MAX(p.final_price), 0)::bigint,
			0::bigint
		FROM products p
		WHERE ` + whereClause(priceConditions)

	addBranch(priceBranch, priceArgs)

	var basePriceArgs []any

	basePriceConditions := []string{
		"p.subcategory_id = " + addBasePriceArg(&basePriceArgs, filter.Subcategory),
	}

	if filter.Subsubcategory != nil {
		basePriceConditions = append(
			basePriceConditions,
			"p.subsubcategory_id = "+addBasePriceArg(&basePriceArgs, *filter.Subsubcategory),
		)
	}

	basePriceWhere := whereClause(basePriceConditions)

	basePriceBranch := `
		SELECT
			'base_price_min'::text,
			''::text,
			COALESCE(MIN(p.final_price), 0)::bigint,
			0::bigint
		FROM products p
		WHERE ` + basePriceWhere + `

		UNION ALL

		SELECT
			'base_price_max'::text,
			''::text,
			COALESCE(MAX(p.final_price), 0)::bigint,
			0::bigint
		FROM products p
		WHERE ` + basePriceWhere

	addBranch(basePriceBranch, basePriceArgs)

	var ratingArgs []any

	ratingConditions := appendFilterConditions(filter, &ratingArgs, excludeFilter{})

	ratingBranch := `
		SELECT
			'rating'::text,
			''::text,
			0::bigint,
			COUNT(*) FILTER (
				WHERE p.rating >= 4
			)::bigint
		FROM products p
		WHERE ` + whereClause(ratingConditions)

	addBranch(ratingBranch, ratingArgs)

	usedCharacteristics, err := r.FindUsedCharsGeneral(
		ctx,
		filter.Subcategory,
		filter.Subsubcategory,
	)
	if err != nil {
		return result, fmt.Errorf(
			"find used characteristics for facets: %w",
			err,
		)
	}

	for _, char := range usedCharacteristics {
		characteristic := char.Value

		var charArgs []any

		conditions := appendFilterConditions(
			filter,
			&charArgs,
			excludeFilter{characteristics: characteristic},
		)

		charArgs = append(charArgs, characteristic)
		facetKeyPosition := len(charArgs)

		charArgs = append(charArgs, characteristic+"_id")
		jsonKeyPosition := len(charArgs)

		charBranch := fmt.Sprintf(`
			SELECT
				'characteristic'::text,
				$%d::text,
				(p.chars_general ->> $%d)::bigint,
				COUNT(*)::bigint
			FROM products p
			WHERE %s
				AND p.chars_general ? $%d
			GROUP BY (p.chars_general ->> $%d)::bigint
    `,
			facetKeyPosition,
			jsonKeyPosition,
			whereClause(conditions),
			jsonKeyPosition,
			jsonKeyPosition,
		)

		addBranch(charBranch, charArgs)
	}

	query := `
    WITH facet_data AS (` + strings.Join(branches, ` UNION ALL `) + `)
    SELECT kind, facet_key, option_id, item_count
    FROM facet_data
	`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return result, fmt.Errorf("get filter facets: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var row facetRow

		if err := rows.Scan(
			&row.Kind,
			&row.Key,
			&row.ID,
			&row.Count,
		); err != nil {
			return result, fmt.Errorf("scan filter facet: %w", err)
		}

		switch row.Kind {
		case "stats":
			result.Count = row.Count
		case "brand":
			result.Brands = append(result.Brands, FacetCount{
				ID:    row.ID,
				Count: row.Count,
			})
		case "subsubcategory":
			result.Subsubcategories = append(
				result.Subsubcategories,
				FacetCount{
					ID:    row.ID,
					Count: row.Count,
				},
			)
		case "price_min":
			result.PriceMin = row.ID
		case "price_max":
			result.PriceMax = row.ID
		case "base_price_min":
			result.BasePriceMin = row.ID
		case "base_price_max":
			result.BasePriceMax = row.ID
		case "rating":
			result.RatingCount = row.Count
		case "characteristic":
			result.Characteristics[row.Key] = append(
				result.Characteristics[row.Key],
				FacetCount{
					ID:    row.ID,
					Count: row.Count,
				},
			)
		}
	}

	if err := rows.Err(); err != nil {
		return result, fmt.Errorf("iterate filter facets: %w", err)
	}

	return result, nil
}
