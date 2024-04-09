package doc

import (
	"secaas_backend/model"
	"time"

	"github.com/kamva/mgm/v3"
)

type UserOrganization struct {
	ID string `bson:"id,omitempty"`
}

type User struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string             `bson:"name,omitempty"`
	Email            model.Email        `bson:"email,omitempty"`
	PassHash         PassHash           `bson:"passHash,omitempty"`
	SymKey           SymKey             `bson:"symKey,omitempty"`
	AsymmKey         AsymmKey           `bson:"asymmKey,omitempty"`
	CreatedAt        time.Time          `bson:"createdAt,omitempty"`
	UpdatedAt        time.Time          `bson:"updatedAt,omitempty"`
	IsBlackListed    bool               `bson:"isBlackListed,omitempty"`
	Organization     []UserOrganization `bson:"organizations,omitempty"`
}

type PassHash struct {
	Hash string `bson:"hash,omitempty"`
	Alg  string `bson:"alg,omitempty"`
}

type SymKey struct {
	EncryptedData string `bson:"encryptedData,omitempty"`
	Alg           string `bson:"alg,omitempty"`
}

type AsymmKey struct {
	Public          string `bson:"public,omitempty"`
	EncryptedPvtKey string `bson:"encryptedPvtKey,omitempty"`
	Alg             string `bson:"alg,omitempty"`
}
