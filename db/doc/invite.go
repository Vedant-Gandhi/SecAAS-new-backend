package doc

import (
	"time"

	"github.com/kamva/mgm/v3"
)

type Invite struct {
	mgm.DefaultModel `bson:",inline"`
	ExpiresAt        time.Time `bson:"expiresAt,omitempty"`
	FromUserEmail    string    `bson:"fromUserEmail"`
	ToUserEmail      string    `bson:"toUserEmail"`
	OrganizationID   string    `bson:"organizationId"`
	SymKey           SymKey    `bson:"symKey`
}
