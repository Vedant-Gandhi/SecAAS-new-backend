package model

type OrganizationID string

func (o OrganizationID) String() string {
	return string(o)
}

type Organization struct {
	ID           OrganizationID `json:"id"`
	CreatedAt    string         `json:"createdAt"`
	UpdatedAt    string         `json:"updatedAt"`
	Name         string         `json:"name"`
	BillingEmail string         `json:"billingEmail"`
	AdminEmail   string         `json:"adminEmail"`
	AsymmKey     AsymmKey       `json:"asymmKey"`
}
