package network

import "github.com/google/uuid"

func NewSessionID() string {
	return uuid.NewString()
}
