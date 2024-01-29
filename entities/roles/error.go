package roles

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrInvalidUUID   = errors.New("invalid uuid")
	ErrBuildingQuery = errors.New("error building query")
)
