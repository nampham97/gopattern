package repository

import (
	"GoPattern/models"
	"context"
	"database/sql"
	"fmt"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	query := "SELECT id, username, role FROM users"
	rows, err := r.db.QueryContext(ctx, query) // Sử dụng QueryContext với context
	if err != nil {
		return nil, fmt.Errorf("Error fetching users: %v", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Role); err != nil {
			return nil, fmt.Errorf("Error scanning row: %v", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error with rows: %v", err)
	}

	return users, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	query := "SELECT id, username, role FROM users WHERE id = $1"
	row := r.db.QueryRowContext(ctx, query, id)

	user := &models.User{}
	if err := row.Scan(&user.ID, &user.Username, &user.Role); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with ID %d not found", id)
		}
		return nil, fmt.Errorf("error scanning row for user with ID %d: %v", id, err)
	}

	return user, nil
}
