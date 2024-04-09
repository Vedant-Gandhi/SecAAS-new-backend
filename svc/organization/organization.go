package organization

import (
	"context"
	"secaas_backend/db/doc"
	"secaas_backend/model"
	"secaas_backend/svc/errors"

	"github.com/kamva/mgm/v3"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

type OrganizationSVC struct {
	logger *logrus.Logger
}

func New(logger *logrus.Logger) *OrganizationSVC {
	u := &OrganizationSVC{logger: logger}
	return u
}

func (o *OrganizationSVC) CreateNew(ctx context.Context, organization model.Organization) (data model.Organization, err error) {

	if organization.AsymmKey.Alg == "" || organization.AsymmKey.EncryptedPvtKey == "" || organization.AsymmKey.Public == "" {
		err = errors.ErrInvalidAsymmetricKey
		return
	}

	if organization.AdminEmail == "" {
		err = errors.ErrInvalidEmail
		return
	}
	docOrganization := &doc.Organization{
		Name:         organization.Name,
		BillingEmail: organization.BillingEmail,
		AdminEmail:   organization.AdminEmail,
		AsymmKey: doc.AsymmKey{
			Public:          organization.AsymmKey.Public,
			EncryptedPvtKey: organization.AsymmKey.EncryptedPvtKey,
			Alg:             organization.AsymmKey.Alg,
		},
	}

	err = mgm.Coll(docOrganization).CreateWithCtx(ctx, docOrganization)

	if err != nil {
		o.logger.WithError(err).Error("error while creating a new organization")
		err = errors.ErrUnknown
		return
	}

	data = o.MapDocToOrganization(docOrganization)

	return
}

func (o *OrganizationSVC) DeleteOrganization(ctx context.Context, organizationId model.OrganizationID) (deleted int, err error) {

	objId, err := primitive.ObjectIDFromHex(organizationId.String())

	if err != nil {
		o.logger.WithContext(ctx).WithError(err).Error("Object ID for organization delete is not valid")
		err = errors.ErrInvalidID
		return
	}

	filter := bson.M{
		"_id": objId,
	}

	docOrg := &doc.Organization{}

	res, err := mgm.Coll(docOrg).DeleteOne(ctx, filter)

	if err != nil {
		o.logger.WithContext(ctx).WithError(err).Error("error while deleting organization")
		err = errors.ErrUnknown
		return
	}

	deleted = int(res.DeletedCount)
	return
}

func (o *OrganizationSVC) MapDocToOrganization(docOrg *doc.Organization) model.Organization {
	org := model.Organization{
		ID:           model.OrganizationID(docOrg.ID.Hex()),
		CreatedAt:    docOrg.CreatedAt.String(),
		UpdatedAt:    docOrg.UpdatedAt.String(),
		Name:         docOrg.Name,
		BillingEmail: docOrg.BillingEmail,
		AdminEmail:   docOrg.AdminEmail,
		AsymmKey: model.AsymmKey{
			Public:          docOrg.AsymmKey.Public,
			EncryptedPvtKey: docOrg.AsymmKey.EncryptedPvtKey,
			Alg:             docOrg.AsymmKey.Alg,
		},
	}

	return org
}
