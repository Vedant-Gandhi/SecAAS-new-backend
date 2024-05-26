package doc

import (
	"time"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SecretUser struct {
	ID   string `bson:"id,omitempty"`
	Role string `bson:"role,omitempty"`
}

type Secret struct {
	mgm.DefaultModel `bson:",inline"`
	EncryptedData    string              `bson:"encryptedData,omitempty"`
	User             SecretUser          `bson:"user,omitempty"`
	Name             string              `bson:"name,omitempty"`
	Description      string              `bson:"description,omitempty"`
	Tags             []string            `bson:"tags,omitempty"`
	CreatorEmail     string              `bson:"creatorEmail,omitempty"`
	Type             string              `bson:"type,omitempty"`
	ReferenceKey     *primitive.ObjectID `bson:"referenceKey,omitempty"`
	OrganizationID   string              `bson:"organizationId"`
	ExpiresAt        time.Time           `bson:"expiresAt,omitempty"`
}
