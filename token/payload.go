package token

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrTokenInvalid = fmt.Errorf("token is invalid")
	ErrExpiredToken = fmt.Errorf("token has expired")
)

// Payload contains the payload data of the token
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// NewPayload creates a new payload with the username and duration
func NewPayload(username string, duration time.Duration) *Payload {
	return &Payload{
		ID:        uuid.New(),
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
}

func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
