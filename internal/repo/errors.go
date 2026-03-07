package repo

import "errors"

var (
	ErrPlayerNotFound      = errors.New("player not found")
	ErrInsufficientChips   = errors.New("insufficient chips")
	ErrRoundAlreadySettled = errors.New("round already settled")
	ErrStorageBusy         = errors.New("storage busy")
	ErrStorageCorrupted    = errors.New("storage corrupted")
)
