package service

import (
	"context"
	"fmt"

	"github.com/AliRizaAynaci/rlaas/internal/database"
	"github.com/AliRizaAynaci/rlaas/internal/models"
)

func GetUserByID(ctx context.Context, userID int) (*models.User, error) {
	dbSvc, err := database.GetInstance()
	if err != nil {
		return nil, fmt.Errorf("db instance: %w", err)
	}

	const q = `
    SELECT id, google_id, email, created_at
    FROM users
    WHERE id = $1
  `
	var u models.User
	err = dbSvc.Pool().QueryRow(ctx, q, userID).
		Scan(&u.ID, &u.GoogleID, &u.Email, &u.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return &u, nil
}
