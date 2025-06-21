package service

import (
	"context"
	"fmt"
	"rlaas/internal/database"
	"rlaas/internal/models"
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

func CreateProject(name, apiKey string) error {
	db, err := database.GetInstance()
	if err != nil {
		return fmt.Errorf("db error: %w", err)
	}

	query := `INSERT INTO projects (name, api_key) VALUES ($1, $2)`
	_, err = db.Pool().Exec(context.Background(), query, name, apiKey)
	return err
}
