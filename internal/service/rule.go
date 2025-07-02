package service

import (
	"context"
	"fmt"
	"github.com/AliRizaAynaci/rlaas/internal/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

var dbPool *pgxpool.Pool

// Connect initializes the database connection
func Connect(pool *pgxpool.Pool) {
	dbPool = pool
}

// RateLimitRule represents a rate limiting rule configuration
type RateLimitRule struct {
	ID            int    `json:"id"`
	ProjectID     int    `json:"project_id"`
	Endpoint      string `json:"endpoint"`
	Strategy      string `json:"strategy"`
	KeyBy         string `json:"key_by"`
	LimitCount    int    `json:"limit_count"`
	WindowSeconds int    `json:"window_seconds"`
}

// Project represents a registered project in the system
type Project struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	ApiKey string `json:"api_key"`
}

// AddRule creates a new rate limit rule in the database
func AddRule(rule *RateLimitRule) error {
	db, err := database.GetInstance()
	if err != nil {
		return err
	}

	query := `
		INSERT INTO rate_limit_rules (project_id, endpoint, strategy, key_by, limit_count, window_seconds)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	err = db.Pool().QueryRow(
		context.Background(),
		query,
		rule.ProjectID,
		rule.Endpoint,
		rule.Strategy,
		rule.KeyBy,
		rule.LimitCount,
		rule.WindowSeconds,
	).Scan(&rule.ID)

	if err != nil {
		return fmt.Errorf("failed to insert rule: %v", err)
	}

	return nil
}

// GetRulesByProjectID retrieves all rate limit rules for a given project ID
func GetRulesByProjectID(projectID int) ([]RateLimitRule, error) {
	db, err := database.GetInstance()
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, project_id, endpoint, strategy, key_by, limit_count, window_seconds 
		FROM rate_limit_rules 
		WHERE project_id = $1
		ORDER BY id DESC`

	rows, err := db.Pool().Query(context.Background(), query, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to query rules: %v", err)
	}
	defer rows.Close()

	var rules []RateLimitRule
	for rows.Next() {
		var rule RateLimitRule
		err := rows.Scan(
			&rule.ID,
			&rule.ProjectID,
			&rule.Endpoint,
			&rule.Strategy,
			&rule.KeyBy,
			&rule.LimitCount,
			&rule.WindowSeconds,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rule: %v", err)
		}
		rules = append(rules, rule)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rules: %v", err)
	}

	return rules, nil
}

// DeleteRule removes a rate limit rule from the database
func DeleteRule(ruleID, projectID int) error {
	db, err := database.GetInstance()
	if err != nil {
		return err
	}

	query := `DELETE FROM rate_limit_rules WHERE id = $1 AND project_id = $2`
	result, err := db.Pool().Exec(context.Background(), query, ruleID, projectID)
	if err != nil {
		return fmt.Errorf("failed to delete rule: %v", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("rule not found or not authorized")
	}

	return nil
}

// UpdateRule modifies an existing rate limit rule in the database
func UpdateRule(rule *RateLimitRule) error {
	db, err := database.GetInstance()
	if err != nil {
		return err
	}

	query := `
		UPDATE rate_limit_rules 
		SET endpoint = $1, strategy = $2, key_by = $3, limit_count = $4, window_seconds = $5
		WHERE id = $6 AND project_id = $7
		RETURNING id`

	err = db.Pool().QueryRow(
		context.Background(),
		query,
		rule.Endpoint,
		rule.Strategy,
		rule.KeyBy,
		rule.LimitCount,
		rule.WindowSeconds,
		rule.ID,
		rule.ProjectID,
	).Scan(&rule.ID)

	if err != nil {
		return fmt.Errorf("failed to update rule: %v", err)
	}

	return nil
}
