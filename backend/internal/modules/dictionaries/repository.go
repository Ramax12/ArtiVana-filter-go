package dictionaries

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	FindAllBrands(ctx context.Context) ([]ProductOption, error)
	FindAllCharsGeneral(ctx context.Context) ([]ProductOptionChars, error)
	FindAllCharsExtra(ctx context.Context) ([]ProductOptionChars, error)
	FindAllCharsGeneralOptions(ctx context.Context) ([]ProductOptionCharsOption, error)
	FindAllCharsExtraOptions(ctx context.Context) ([]ProductOptionCharsOption, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) FindAllBrands(ctx context.Context) ([]ProductOption, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			id,
			name
		FROM brands
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var brands []ProductOption

	for rows.Next() {
		var brand ProductOption

		err := rows.Scan(
			&brand.ID,
			&brand.Name,
		)

		if err != nil {
			return nil, err
		}

		brands = append(brands, brand)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return brands, nil
}

func (r *repository) FindAllCharsGeneral(ctx context.Context) ([]ProductOptionChars, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			id,
			name,
			value
		FROM chars_general
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var chars []ProductOptionChars

	for rows.Next() {
		var char ProductOptionChars

		err := rows.Scan(
			&char.ID,
			&char.Name,
			&char.Value,
		)

		if err != nil {
			return nil, err
		}

		chars = append(chars, char)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return chars, nil
}

func (r *repository) FindAllCharsExtra(ctx context.Context) ([]ProductOptionChars, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			id,
			name,
			value
		FROM chars_extra
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var chars []ProductOptionChars

	for rows.Next() {
		var char ProductOptionChars

		err := rows.Scan(
			&char.ID,
			&char.Name,
			&char.Value,
		)

		if err != nil {
			return nil, err
		}

		chars = append(chars, char)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return chars, nil
}

func (r *repository) FindAllCharsGeneralOptions(ctx context.Context) ([]ProductOptionCharsOption, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			id,
			name,
			parent_id
		FROM chars_general_options
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var chars []ProductOptionCharsOption

	for rows.Next() {
		var char ProductOptionCharsOption

		err := rows.Scan(
			&char.ID,
			&char.Name,
			&char.ParentID,
		)

		if err != nil {
			return nil, err
		}

		chars = append(chars, char)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return chars, nil
}

func (r *repository) FindAllCharsExtraOptions(ctx context.Context) ([]ProductOptionCharsOption, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			id,
			name,
			parent_id
		FROM chars_extra_options
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var chars []ProductOptionCharsOption

	for rows.Next() {
		var char ProductOptionCharsOption

		err := rows.Scan(
			&char.ID,
			&char.Name,
			&char.ParentID,
		)

		if err != nil {
			return nil, err
		}

		chars = append(chars, char)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return chars, nil
}
