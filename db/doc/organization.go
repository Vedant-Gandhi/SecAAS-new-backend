package doc

import (
	"time"

	"github.com/kamva/mgm/v3"
)

type Organization struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string    `bson:"name,omitempty"`
	BillingEmail     string    `bson:"billingEmail,omitempty"`
	AdminEmail       string    `bson:"adminEmail,omitempty"`
	SymmKey          SymKey    `bson:"symKey,omitempty"`
	SoftDelete       bool      `bson:"softDelete,omitempty"`
	DeleteTimeStamp  time.Time `bson:"deleteTs,omitempty"`
}
