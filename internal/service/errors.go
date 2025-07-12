package service

import "errors"

var (
	ErrProjectNotFound  = errors.New("project not found for given API key")
	ErrEndpointNotOwned = errors.New("endpoint does not belong to this project")
)
