package model

import "time"

type UserID string

func (u UserID) String() string {
	return string(u)
}

type Email string

func (e Email) String() string {
	return string(e)
}

type UserOrganization struct {
	ID      string `json:"id"`
	IsAdmin string `json:"isAdmin"`
	PvtKey  string `json:"pvtKey"`
}

type User struct {
	ID            UserID             `json:"id"`
	Name          string             `json:"name"`
	Email         Email              `json:"email"`
	PassHash      PassHash           `json:"passHash"`
	SymKey        SymKey             `json:"symKey"`
	AsymmKey      AsymmKey           `json:"asymmKey"`
	CreatedAt     time.Time          `json:"createdAt"`
	UpdatedAt     time.Time          `json:"updatedAt"`
	IsBlackListed bool               `json:"isBlackListed"`
	Organization  []UserOrganization `json:"organizations"`
}

type PassHash struct {
	Hash string `json:"hash"`
	Alg  string `json:"alg"`
}

type SymKey struct {
	EncryptedData string `json:"encryptedData"`
	Alg           string `json:"alg"`
}

type AsymmKey struct {
	Public          string `json:"public"`
	EncryptedPvtKey string `json:"encryptedPvtKey"`
	Alg             string `json:"alg"`
}
