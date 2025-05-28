package role

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"idm/inner/database"
	"time"
)

type Role struct {
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

func (r *Repository) FindById(id int64) (*Role, error) {
	var role Role

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.db.GetContext(ctx, &role, "SELECT * FROM roles WHERE id = $1", id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, database.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &role, err
}

func (r *Repository) FindAll() ([]*Role, error) {
	var roles []*Role

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.db.SelectContext(ctx, &roles, "SELECT * FROM roles")

	return roles, err
}

func (r *Repository) FindByIds(ids []int64) ([]*Role, error) {
	var roles []*Role
	err := r.db.Select(&roles, "SELECT * FROM roles WHERE id = ANY($1)", pq.Array(ids))

	return roles, err
}

func (r *Repository) Create(role *Role) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx,
		"INSERT INTO roles (name) VALUES ($1) RETURNING id, created_at, updated_at",
		role.Name).Scan(&role.Id, &role.CreatedAt, &role.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Remove(id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx, "DELETE FROM roles WHERE id = $1", id)

	return err
}

func (r *Repository) RemoveByIds(ids []int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx, "DELETE FROM roles WHERE id = ANY($1)", pq.Array(ids))

	return err
}
