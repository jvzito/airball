package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/jvzito/airball/internal/models"
)

type UserRepo struct{ db *sqlx.DB }

func NewUserRepo(db *sqlx.DB) *UserRepo { return &UserRepo{db: db} }

func (r *UserRepo) Create(ctx context.Context, email, passwordHash string) (*models.User, error) {
	u := &models.User{}
	err := r.db.QueryRowxContext(ctx,
		`INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, email, password_hash, created_at`,
		email, passwordHash,
	).StructScan(u)
	return u, err
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	u := &models.User{}
	err := r.db.QueryRowxContext(ctx,
		`SELECT id, email, password_hash, created_at FROM users WHERE email=$1`,
		email,
	).StructScan(u)
	return u, err
}
