package employee

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error creating transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("creating employee panic: %v", p)
			errTx := tx.Rollback()
			if errTx != nil {
				err = fmt.Errorf("creating employee: rolling back transaction errors: %w, %w", err, errTx)
			}
		} else if err != nil {
			errTx := tx.Rollback()
			if errTx != nil {
				err = fmt.Errorf("creating employee: rolling back transaction errors: %w, %w", err, errTx)
			}
		} else {
			errTx := tx.Commit()
			if errTx != nil {
				err = fmt.Errorf("creating employee: commiting transaction error: %w", errTx)
			}
		}
	}()

	var exists int
	err = tx.GetContext(ctx, &exists, "SELECT COUNT(*) FROM employees WHERE name = $1", employee.Name)
	if err != nil {
		return err
	}
	if exists > 0 {
		return errors.New("employee already exists")
	}

	err = tx.QueryRowContext(ctx,
		"INSERT INTO employees (name) VALUES ($1) RETURNING id, created_at, updated_at",
		employee.Name,
	).Scan(&employee.Id, &employee.CreatedAt, &employee.UpdatedAt)

	return err
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
