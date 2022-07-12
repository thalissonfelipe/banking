package testdata

import "errors"

var (
	ErrRepositoryFailsToFetch = errors.New("repository fails to fetch")
	ErrRepositoryFailsToSave  = errors.New("repository fails to save")
	ErrUsecaseFails           = errors.New("usecase error")
	ErrHashFails              = errors.New("hash error")
)
