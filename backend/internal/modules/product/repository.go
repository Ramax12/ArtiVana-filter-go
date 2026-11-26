package product

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	GetByID(ctx context.Context, id int64) (*Product, error)
	FindByGroupID(ctx context.Context, groupID int64, excludeID int64) ([]ProductByGroupID, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetByID(
	ctx context.Context,
	id int64,
) (*Product, error) {
	var product Product
	var shippingOptions []string

	err := r.db.QueryRow(ctx, `
		SELECT
			id,
			slug,
			group_id,
			created_at,
			category_id,
			subcategory_id,
			subsubcategory_id,
			name,
			article,
			brand_id,
			model,
			original_price,
			final_price,
			stock,
			images,
			rating,
			chars_general,
			chars_extra,
			package,
			description,
			shipping_options::text[],
			tags
		FROM products
		WHERE id = $1;
	`, id).Scan(
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
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	product.ShippingOptions = make([]DeliveryOption, len(shippingOptions))
	for i, option := range shippingOptions {
		product.ShippingOptions[i] = DeliveryOption(option)
	}

	return &product, nil
}

func (r *repository) FindByGroupID(
	ctx context.Context,
	groupID int64,
	excludeID int64,
) ([]ProductByGroupID, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			id,
			slug,
			name,
			model,
			images,
			group_id
		FROM products
		WHERE group_id = $1
			AND id != $2
		ORDER BY id;
	`, groupID, excludeID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []ProductByGroupID

	for rows.Next() {
		var product ProductByGroupID

		err := rows.Scan(
			&product.ID,
			&product.Slug,
			&product.Name,
			&product.Model,
			&product.Images,
			&product.GroupID,
		)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
