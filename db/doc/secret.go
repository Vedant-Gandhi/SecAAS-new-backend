package doc

import (
	"time"

	"github.com/kamva/mgm/v3"
)

type SecretUser struct {
	ID   string `bson:"id,omitempty"`
	Role string `bson:"role,omitempty"`
}

type Secret struct {
	mgm.DefaultModel `bson:",inline"`
	EncryptedData    string     `bson:"encryptedData,omitempty"`
	User             SecretUser `bson:"user,omitempty"`
	Name             string     `bson:"name,omitempty"`
	Description      string     `bson:"description,omitempty"`
	Tags             []string   `bson:"tags,omitempty"`
	CreatorEmail     string     `bson:"creatorEmail,omitempty"`
	Type             string     `bson:"type,omitempty"`
	ReferenceKey     string     `bson:"referenceKey,omitempty"`
	ExpiresAt        time.Time  `bson:"expiresAt,omitempty"`
}
