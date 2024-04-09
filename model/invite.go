package model

import "time"

type InviteID string

func (i InviteID) String() string {
	return string(i)
}

type Invite struct {
	ID             InviteID  `json:"id"`
	CreatedAt      time.Time `json:"createdAt"`
	ExpiresAt      time.Time `json:"expiresAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	FromUserEmail  string    `json:"fromUserEmail"`
	ToUserEmail    string    `json:"toUserEmail"`
	OrganizationID string    `json:"organizationId"`
}
