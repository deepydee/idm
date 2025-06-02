package employee

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"idm/inner/database"
	"time"
)

type Employee struct {
	Id        int64     `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindById(id int64) (*Employee, error) {
	var employee Employee

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.db.GetContext(ctx, &employee, "SELECT * FROM employees WHERE id = $1", id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, database.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &employee, err
}

func (r *Repository) FindAll() ([]*Employee, error) {
	var employees []*Employee

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.db.SelectContext(ctx, &employees, "SELECT * FROM employees")

	return employees, err
}

func (r *Repository) FindByIds(ids []int64) ([]*Employee, error) {
	var employees []*Employee
	err := r.db.Select(&employees, "SELECT * FROM employees WHERE id = ANY($1)", pq.Array(ids))

	return employees, err
}

func (r *Repository) Create(employee *Employee) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx,
		"INSERT INTO employees (name) VALUES ($1) RETURNING id, created_at, updated_at",
		employee.Name,
	).Scan(&employee.Id, &employee.CreatedAt, &employee.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Remove(id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx, "DELETE FROM employees WHERE id = $1", id)

	return err
}

func (r *Repository) RemoveByIds(ids []int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx, "DELETE FROM employees WHERE id = ANY($1)", pq.Array(ids))

	return err
}
