package doc

import "github.com/kamva/mgm/v3"

type Organization struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string   `bson:"name,omitempty"`
	BillingEmail     string   `bson:"billingEmail,omitempty"`
	AdminEmail       string   `bson:"adminEmail,omitempty"`
	AsymmKey         AsymmKey `bson:"asymmKey,omitempty"`
}
