package secret

import (
	"context"
	"secaas_backend/db/doc"
	"secaas_backend/model"
	"secaas_backend/svc/errors"

	"github.com/kamva/mgm/v3"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SecretsSVC struct {
	logger *logrus.Logger
}

func New(logger *logrus.Logger) *SecretsSVC {
	u := &SecretsSVC{logger: logger}
	return u
}
func (s *SecretsSVC) GetListForUser(ctx context.Context, userId model.UserID, organizationId string, params model.PaginationParams) (sec model.Secret, err error) {

	if userId == "" {
		s.logger.WithContext(ctx).Error("invalid user id to get secret list for user")
		err = errors.ErrInvalidID
		return
	}

	if organizationId == "" {
		s.logger.WithContext(ctx).Error("invalid organization id to get secret list for user")
		err = errors.ErrInvalidOrganizationID
		return
	}

	return
}

func (s *SecretsSVC) Create(ctx context.Context, data model.Secret) (sec model.Secret, err error) {
	objRefId, _ := primitive.ObjectIDFromHex(data.ReferenceKey)
	docSecret := &doc.Secret{
		EncryptedData: data.EncryptedData,
		User: doc.SecretUser{
			ID:   string(data.User.ID),
			Role: data.User.Role,
		},
		Name:           data.Name,
		Description:    data.Description,
		Tags:           data.Tags,
		CreatorEmail:   data.CreatorEmail,
		Type:           data.Type,
		ReferenceKey:   objRefId,
		ExpiresAt:      data.ExpiresAt,
		OrganizationID: data.OrganizationID,
	}

	err = mgm.Coll(docSecret).Create(docSecret)
	if err != nil {
		s.logger.WithContext(ctx).WithError(err).Error("Error while creatin secret")
		err = errors.ErrUnknown
		return
	}

	sec = model.Secret{
		ID:            docSecret.ID.Hex(),
		EncryptedData: docSecret.EncryptedData,
		CreatedAt:     docSecret.CreatedAt,
		UpdatedAt:     docSecret.UpdatedAt,
		Name:          docSecret.Name,
		User: model.SecretUser{
			ID:   model.UserID(docSecret.User.ID),
			Role: docSecret.User.Role,
		},
		Description:    docSecret.Description,
		CreatorEmail:   docSecret.CreatorEmail,
		Tags:           docSecret.Tags,
		Type:           docSecret.Type,
		ReferenceKey:   docSecret.ReferenceKey.Hex(),
		ExpiresAt:      docSecret.ExpiresAt,
		OrganizationID: docSecret.OrganizationID,
	}

	return
}

func (s *SecretsSVC) GetAllSecretsforUser(ctx context.Context, userId model.UserID, params model.PaginationParams) (data []model.Secret, err error) {

	secretDoc := &doc.Secret{}

	filter := bson.M{
		"user.id": userId,
	}

	findOptions := options.Find().SetLimit(int64(params.Limit)).SetSkip(int64(params.Skip)).SetSort(bson.D{
		{"updatedAt", -1},
	})

	cursor, err := mgm.Coll(secretDoc).Find(ctx, filter, findOptions)

	for cursor.Next(ctx) {
		var curDoc doc.Secret

		err := cursor.Decode(curDoc)

		if err != nil {
			s.logger.WithContext(ctx).WithError(err).Error("error while decoding secret document")
			continue
		}

		modelSecret := s.MapDocToModelSecret(curDoc)

		data = append(data, modelSecret)
	}
	return
}

func (s *SecretsSVC) GetAllSecretsforOrganization(ctx context.Context, orgId model.OrganizationID, params model.PaginationParams) (data []model.Secret, err error) {

	secretDoc := &doc.Secret{}

	filter := bson.M{
		"organizationId": orgId,
	}

	findOptions := options.Find().SetLimit(int64(params.Limit)).SetSkip(int64(params.Skip)).SetSort(bson.D{
		{"updatedAt", -1},
	})

	cursor, err := mgm.Coll(secretDoc).Find(ctx, filter, findOptions)

	for cursor.Next(ctx) {
		var curDoc doc.Secret

		err := cursor.Decode(&curDoc)

		if err != nil {
			s.logger.WithContext(ctx).WithError(err).Error("error while decoding secret document")
			continue
		}

		modelSecret := s.MapDocToModelSecret(curDoc)

		data = append(data, modelSecret)
	}
	return
}

func (s *SecretsSVC) GetAllSecretsforaUserInOrganization(ctx context.Context, userId model.UserID, orgId model.OrganizationID, params model.PaginationParams) (data []model.Secret, err error) {

	secretDoc := &doc.Secret{}

	filter := bson.M{
		"organizationId": orgId,
		"user.id":        userId,
	}

	findOptions := options.Find().SetLimit(int64(params.Limit)).SetSkip(int64(params.Skip)).SetSort(bson.D{
		{"updatedAt", 1},
	})

	cursor, err := mgm.Coll(secretDoc).Find(ctx, filter, findOptions)
	if err != nil {
		s.logger.WithContext(ctx).WithError(err).Error("Error while fetching organization for user secrets")
		err = errors.ErrUnknown
		return
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var curDoc doc.Secret

		err := cursor.Decode(&curDoc)

		if err != nil {
			s.logger.WithContext(ctx).WithError(err).Error("error while decoding secret document")
			continue
		}

		modelSecret := s.MapDocToModelSecret(curDoc)

		data = append(data, modelSecret)
	}
	return
}

func (s *SecretsSVC) MapDocToModelSecret(docSecret doc.Secret) model.Secret {
	return model.Secret{
		ID:            docSecret.ID.Hex(),
		EncryptedData: docSecret.EncryptedData,
		CreatedAt:     docSecret.CreatedAt,
		UpdatedAt:     docSecret.UpdatedAt,
		Name:          docSecret.Name,
		User: model.SecretUser{
			ID:   model.UserID(docSecret.User.ID),
			Role: docSecret.User.Role,
		},
		Description:    docSecret.Description,
		CreatorEmail:   docSecret.CreatorEmail,
		Tags:           docSecret.Tags,
		Type:           docSecret.Type,
		ReferenceKey:   docSecret.ReferenceKey.Hex(),
		ExpiresAt:      docSecret.ExpiresAt,
		OrganizationID: docSecret.OrganizationID,
	}
}
