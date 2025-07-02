package service

import (
	"context"
	"fmt"
	"github.com/AliRizaAynaci/rlaas/internal/database"
	"github.com/AliRizaAynaci/rlaas/internal/models"
)

func GetProjectByAPIKey(apiKey string) (*models.Project, error) {
	db, err := database.GetInstance()
	if err != nil {
		return nil, fmt.Errorf("db error: %w", err)
	}

	query := `SELECT id, name, api_key FROM projects WHERE api_key = $1`
	var project models.Project

	err = db.Pool().QueryRow(context.Background(), query, apiKey).Scan(
		&project.ID,
		&project.Name,
		&project.ApiKey,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	return &project, nil
}

func CreateProject(ctx context.Context, userID int, name, apiKey string) (*models.Project, error) {
	// DB servisini al
	dbSvc, err := database.GetInstance()
	if err != nil {
		return nil, fmt.Errorf("db error: %w", err)
	}

	const q = `
      INSERT INTO projects (user_id, name, api_key)
      VALUES ($1, $2, $3)
      RETURNING id, user_id, name, api_key, created_at
    `
	var p models.Project
	err = dbSvc.Pool().QueryRow(ctx, q,
		userID,
		name,
		apiKey,
	).Scan(
		&p.ID,
		&p.UserID,
		&p.Name,
		&p.ApiKey,
		&p.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("CreateProject: %w", err)
	}
	return &p, nil
}
