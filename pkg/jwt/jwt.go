package jwt

import "time"

type Manager struct {
	secret string
	ttl    time.Duration
}

func New(secret string, ttl time.Duration) *Manager {
	return &Manager{
		secret: secret,
		ttl:    ttl,
	}
}
