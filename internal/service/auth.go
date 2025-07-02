package service

import (
	"context"
	"fmt"
	"github.com/AliRizaAynaci/rlaas/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CreateUser saves a new user or returns existing one by GoogleID
func CreateUser(ctx context.Context, pool *pgxpool.Pool, googleID, email string) (*models.User, error) {
	var u models.User
	// Try Fetch
	err := pool.QueryRow(ctx,
		`SELECT id, google_id, email, created_at FROM users WHERE google_id=$1`, googleID).
		Scan(&u.ID, &u.GoogleID, &u.Email, &u.CreatedAt)
	if err == nil {
		return &u, nil
	}
	// Insert
	err = pool.QueryRow(ctx,
		`INSERT INTO users (google_id, email) VALUES ($1,$2) RETURNING id, google_id, email, created_at`,
		googleID, email).
		Scan(&u.ID, &u.GoogleID, &u.Email, &u.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("CreateUser: %w", err)
	}
	return &u, nil
}

// GetUserByGoogleID returns a user by its GoogleID.
func GetUserByGoogleID(ctx context.Context, pool *pgxpool.Pool, googleID string) (*models.User, error) {
	var u models.User
	err := pool.QueryRow(ctx,
		`SELECT id, google_id, email, created_at FROM users WHERE google_id=$1`, googleID).
		Scan(&u.ID, &u.GoogleID, &u.Email, &u.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("GetUserByGoogleID: %w", err)
	}
	return &u, nil
}
