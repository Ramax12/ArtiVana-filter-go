package categories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	FindAll(ctx context.Context) ([]Category, []Subcategory, []Subsubcategory, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) FindAll(ctx context.Context) ([]Category, []Subcategory, []Subsubcategory, error) {
	categoryRows, err := r.db.Query(ctx, `
		SELECT
			id,
			slug,
			name,
			image
		FROM categories
	`)

	if err != nil {
		return nil, nil, nil, err
	}

	defer categoryRows.Close()

	var categories []Category

	for categoryRows.Next() {
		var category Category

		err := categoryRows.Scan(
			&category.ID,
			&category.Slug,
			&category.Name,
			&category.Image,
		)

		if err != nil {
			return nil, nil, nil, err
		}

		categories = append(categories, category)
	}

	if err := categoryRows.Err(); err != nil {
		return nil, nil, nil, err
	}

	subcategoryRows, err := r.db.Query(ctx, `
		SELECT
			id,
			slug,
			name,
			image,
			category_id
		FROM subcategories
	`)

	if err != nil {
		return nil, nil, nil, err
	}

	defer subcategoryRows.Close()

	var subcategories []Subcategory

	for subcategoryRows.Next() {
		var subcategory Subcategory

		err := subcategoryRows.Scan(
			&subcategory.ID,
			&subcategory.Slug,
			&subcategory.Name,
			&subcategory.Image,
			&subcategory.CategoryID,
		)

		if err != nil {
			return nil, nil, nil, err
		}

		subcategories = append(subcategories, subcategory)
	}

	if err := subcategoryRows.Err(); err != nil {
		return nil, nil, nil, err
	}

	subsubcategoryRows, err := r.db.Query(ctx, `
		SELECT
			id,
			slug,
			name,
			image,
			subcategory_id
		FROM subsubcategories
	`)

	if err != nil {
		return nil, nil, nil, err
	}

	defer subsubcategoryRows.Close()

	var subsubcategories []Subsubcategory

	for subsubcategoryRows.Next() {
		var subsubcategory Subsubcategory

		err := subsubcategoryRows.Scan(
			&subsubcategory.ID,
			&subsubcategory.Slug,
			&subsubcategory.Name,
			&subsubcategory.Image,
			&subsubcategory.SubcategoryID,
		)

		if err != nil {
			return nil, nil, nil, err
		}

		subsubcategories = append(subsubcategories, subsubcategory)
	}

	if err := subsubcategoryRows.Err(); err != nil {
		return nil, nil, nil, err
	}

	return categories, subcategories, subsubcategories, nil
}
