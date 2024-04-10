package model

import "time"

type SecretUser struct {
	ID   UserID `json:"id"`
	Role string `json:"role"`
}

type Secret struct {
	ID            string     `json:"id"`
	EncryptedData string     `json:"encryptedData"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	User          SecretUser `json:"user"`
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	Tags          []string   `json:"tags"`
	CreatorEmail  string     `json:"creatorEmail"`
	Type          string     `json:"type"`
	ReferenceKey  string     `json:"referenceKey"`
	ExpiresAt     time.Time  `json:"expiresAt,omitempty"`
}
