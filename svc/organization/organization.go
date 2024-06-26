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

	if organization.SymmKey.Alg == "" || organization.SymmKey.EncryptedData == "" {
		err = errors.ErrInvalidSymmetricKey
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
		SymmKey: doc.SymKey{
			EncryptedData: organization.SymmKey.EncryptedData,
			Alg:           organization.SymmKey.Alg,
		},
	}

	err = mgm.Coll(docOrganization).CreateWithCtx(ctx, docOrganization)

	if err != nil {
		o.logger.WithError(err).Error("error while creating a new organization")
		err = errors.ErrUnknown
		return
	}

	data = o.MapDocToOrganization(docOrganization)

	// Add to user doc.
	userDoc := &doc.User{}

	userFilter := bson.M{
		"email": organization.AdminEmail,
	}

	userUpdate := bson.M{
		"$push": bson.M{
			"organizations": doc.UserOrganization{
				ID:      data.ID.String(),
				IsAdmin: true,
				PvtKey:  docOrganization.SymmKey.EncryptedData,
			},
		},
	}

	updateRes, err := mgm.Coll(userDoc).UpdateOne(ctx, userFilter, userUpdate)

	if err != nil {
		o.logger.WithContext(ctx).WithField("Admin Email", data.AdminEmail).WithError(err).Error("Failed to add the organization entry in the admin.")
		err = nil
	}

	if updateRes.ModifiedCount == 0 {
		o.logger.WithContext(ctx).WithField("Admin Email", data.AdminEmail).Print("Could not update the email with organization.")
	} else {
		o.logger.WithContext(ctx).WithField("Admin Email", data.AdminEmail).Print("Added organization to the email.")
	}

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

	// Delete all the users related to that organization.
	docUser := &doc.User{}

	userFilter := bson.M{
		"organizations.id": organizationId,
	}

	userUpdate := bson.M{
		"$pull": bson.M{
			"organizations": bson.M{
				"id": organizationId,
			},
		},
	}

	userRes, userErr := mgm.Coll(docUser).UpdateMany(ctx, userFilter, userUpdate)

	if userErr != nil {
		o.logger.WithContext(ctx).WithError(err).Error("Failed to delete the users associated with organization.")
	}

	o.logger.WithContext(ctx).WithField("userMatched", userRes.MatchedCount).WithField("userUpdated", userRes.ModifiedCount).Debug("removed organization from users after its deletion.")

	return
}

func (o *OrganizationSVC) GetOrganizationList(ctx context.Context, organizationIds []model.OrganizationID) (orgs []model.Organization, err error) {

	organizationObjIds := []primitive.ObjectID{}

	for _, orgId := range organizationIds {
		if orgId != "" {
			objId, err := primitive.ObjectIDFromHex(orgId.String())

			if err != nil {
				o.logger.WithContext(ctx).WithField("Organization ID", orgId).WithError(err).Error("Object ID for organization get is not valid")
				continue
			}

			organizationObjIds = append(organizationObjIds, objId)
		}

	}

	filter := bson.M{
		"_id": bson.M{
			"$in": organizationObjIds,
		},
	}

	docOrg := &doc.Organization{}

	res, err := mgm.Coll(docOrg).Find(ctx, filter)

	if err != nil {
		o.logger.WithContext(ctx).WithError(err).Error("error while fetching list of organizations")
		err = errors.ErrUnknown
		return
	}

	defer res.Close(ctx)

	mappedOrgs := []model.Organization{}

	for res.Next(ctx) {
		var curDoc doc.Organization

		err = res.Decode(&curDoc)

		if err != nil {
			o.logger.WithContext(ctx).WithError(err).Error("error while decoding organisation from doc in get organisations")
			err = nil
			continue
		}

		mappedOrgs = append(mappedOrgs, o.MapDocToOrganization(&curDoc))
	}

	orgs = mappedOrgs

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
		SymmKey: model.SymKey{
			EncryptedData: docOrg.SymmKey.EncryptedData,
			Alg:           docOrg.SymmKey.Alg,
		},
	}

	return org
}
