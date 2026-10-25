package repository

import "errors"

var (
	ErrProjectNotFound    = errors.New("project not found")
	ErrProjectDuplicate   = errors.New("project name already exists")
	ErrProviderNotFound   = errors.New("provider not found")
	ErrModelRouteNotFound = errors.New("model route not found")
	ErrAPIKeyNotFound     = errors.New("api key not found")
)
