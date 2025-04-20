package repository

import (
	"GoPattern/models"
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type ProvincesRepository struct {
	db *sqlx.DB
}

func NewProvinceRepository(db *sqlx.DB) *ProvincesRepository {
	return &ProvincesRepository{db: db}
}

func (r *ProvincesRepository) GetProvinces(ctx context.Context) ([]models.Provinces, error) {
	query := "SELECT id, name FROM provinces"  // Log the context for debugging
	rows, err := r.db.QueryContext(ctx, query) // Sử dụng QueryContext với context
	if err != nil {
		return nil, fmt.Errorf("error fetching province: %v", err)
	}
	defer rows.Close()

	var provinces []models.Provinces
	for rows.Next() {
		var province models.Provinces
		if err := rows.Scan(&province.ID, &province.Name); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		provinces = append(provinces, province)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error with rows: %v", err)
	}

	return provinces, nil
}
