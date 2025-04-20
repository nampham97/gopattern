package repository

import (
	"GoPattern/internal/params"
	"GoPattern/models"
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	query := "SELECT id, username, role FROM users"
	rows, err := r.db.QueryContext(ctx, query) // Sử dụng QueryContext với context
	if err != nil {
		return nil, fmt.Errorf("error fetching users: %v", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Role); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error with rows: %v", err)
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

func (r *UserRepository) GetUsers(filter models.UserFilter) ([]models.User, error) {
	var users []models.User
	var args []interface{}
	conditions := []string{"1=1"} // mặc định luôn đúng

	if filter.Search != "" {
		args = append(args, "%"+filter.Search+"%")
		conditions = append(conditions, "(username ILIKE $"+strconv.Itoa(len(args))+")") //+" OR email ILIKE $"+strconv.Itoa(len(args))+")")
	}

	if filter.Role != "" {
		args = append(args, filter.Role)
		conditions = append(conditions, "role = $"+strconv.Itoa(len(args)))
	}

	offset := (filter.Page - 1) * filter.Limit
	args = append(args, filter.Limit, offset)
	fmt.Println("args:", args)
	fmt.Println("conditions:", conditions)
	query := `
		SELECT id, username, role 
		FROM users 
		WHERE ` + strings.Join(conditions, " AND ") + `
		ORDER BY id DESC
		LIMIT $` + strconv.Itoa(len(args)-1) + ` OFFSET $` + strconv.Itoa(len(args))
	log.Println("query:", query)
	err := r.db.Select(&users, query, args...)
	log.Println("error", err)
	return users, err
}

func (r *UserRepository) GetPaginatedUsers(page, limit int) ([]models.User, error) {
	offset := (page - 1) * limit
	query := `SELECT id, username, role FROM users ORDER BY id LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Role); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepository) FindUsers(q params.QueryParams) ([]models.User, error) {
	var users []models.User
	query := "SELECT * FROM users WHERE 1=1"
	args := []interface{}{}

	// Xử lý tìm kiếm
	if q.Search != "" {
		query += " AND (username ILIKE $1)"
		args = append(args, "%"+q.Search+"%")
	}

	// Xác thực và xử lý sắp xếp
	validSortBy := map[string]bool{
		"id":       true,
		"username": true,
		"role":     true,
	}
	if !validSortBy[q.SortBy] {
		return nil, fmt.Errorf("invalid sort field: %s", q.SortBy)
	}

	if q.Order != "ASC" && q.Order != "DESC" {
		return nil, fmt.Errorf("invalid order: %s", q.Order)
	}

	query += fmt.Sprintf(" ORDER BY %s %s", q.SortBy, q.Order)

	// Xử lý phân trang
	query += " LIMIT $2 OFFSET $3"
	args = append(args, q.Limit, (q.Page-1)*q.Limit)

	// Thực thi truy vấn
	err := r.db.Select(&users, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}

	return users, nil
}
