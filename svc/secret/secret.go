package secret

import (
	"context"
	"secaas_backend/db/doc"
	"secaas_backend/model"
	"secaas_backend/svc/errors"
	"secaas_backend/svc/user"

	"github.com/kamva/mgm/v3"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SecretsSVC struct {
	logger  *logrus.Logger
	userSvc *user.UserSVC
}

func New(logger *logrus.Logger, userSvc *user.UserSVC) *SecretsSVC {
	u := &SecretsSVC{logger: logger, userSvc: userSvc}
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
	objRefId, _ := primitive.ObjectIDFromHex(*data.ReferenceKey)
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
		ReferenceKey:   &objRefId,
		ExpiresAt:      data.ExpiresAt,
		OrganizationID: data.OrganizationID,
	}

	err = mgm.Coll(docSecret).Create(docSecret)
	if err != nil {
		s.logger.WithContext(ctx).WithError(err).Error("Error while creating secret")
		err = errors.ErrUnknown
		return
	}
	docRefKey := docSecret.ReferenceKey.Hex()
	sec = model.Secret{
		ID:            model.SecretID(docSecret.ID.Hex()),
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
		ReferenceKey:   &docRefKey,
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

func (s *SecretsSVC) GetAllUsersforSecretByAdmin(ctx context.Context, orgId model.OrganizationID, originalKeyID string, params model.PaginationParams) (data []model.Secret, err error) {

	secretDoc := &doc.Secret{}

	objRefKey, _ := primitive.ObjectIDFromHex(originalKeyID)

	filter := bson.M{
		"organizationId": orgId,
		"referenceKey":   objRefKey,
	}

	s.logger.Error(filter)

	findOptions := options.Find().SetLimit(int64(params.Limit)).SetSkip(int64(params.Skip)).SetSort(bson.D{
		{"updatedAt", 1},
	})

	cursor, err := mgm.Coll(secretDoc).Find(ctx, filter, findOptions)
	if err != nil {
		s.logger.WithContext(ctx).WithError(err).Error("Error while fetching users for a particular secret.")
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
	secretModel := model.Secret{
		ID:            model.SecretID(docSecret.ID.Hex()),
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
		ExpiresAt:      docSecret.ExpiresAt,
		OrganizationID: docSecret.OrganizationID,
	}

	if docSecret.ReferenceKey != nil {
		refKey := docSecret.ReferenceKey.Hex()
		secretModel.ReferenceKey = &refKey
	}

	return secretModel
}

func (s *SecretsSVC) ShareSecret(c context.Context, originalID model.SecretID, endUserEmail []model.SecretUser) (resultSet map[string]bool, err error) {
	logger := s.logger.WithContext(c).WithField("shareId", originalID.String())
	resultSet = make(map[string]bool)

	bsonId, err := primitive.ObjectIDFromHex(originalID.String())

	if err != nil {
		logger.WithError(err).Error("error while converting key to bson in secret share")
		err = errors.ErrInvalidID
		return
	}

	secretDoc := &doc.Secret{}
	insertDocs := []interface{}{}

	err = mgm.Coll(secretDoc).FindByID(bsonId, secretDoc)

	if err != nil {
		logger.WithError(err).Error("failed to get the doc by id")
		err = errors.ErrSecretNotFound
		return
	}

	// Loop over all the user emails.
	for _, userDoc := range endUserEmail {

		// Check if users exists or not.
		_, err := s.userSvc.GetByID(c, userDoc.ID)

		// If error occurs set the current index as false for processing and continue.
		if err != nil {
			logger.WithError(err).Error("error while fetching user data")
			err = nil
			resultSet[userDoc.ID.String()] = false
			continue
		}

		refKey := secretDoc.ID

		newInsertDoc := doc.Secret{
			Description:   secretDoc.Description,
			EncryptedData: secretDoc.EncryptedData,
			CreatorEmail:  secretDoc.CreatorEmail,
			ReferenceKey:  &refKey,
			User: doc.SecretUser{
				ID:   userDoc.ID.String(),
				Role: userDoc.Role,
			},
			Tags:           secretDoc.Tags,
			OrganizationID: secretDoc.OrganizationID,
			ExpiresAt:      secretDoc.ExpiresAt,
			Type:           secretDoc.Type,
		}
		insertDocs = append(insertDocs, newInsertDoc)
	}

	if len(insertDocs) == 0 {
		return
	}

	_, err = mgm.Coll(secretDoc).InsertMany(c, insertDocs)

	// If error is not empty then log it and reset all additions.
	if err != nil {
		logger.WithError(err).Error("failed to share secret among the user")
		err = errors.ErrSecretShareFailed
		return
	}

	// Set the non existing keys as entered.
	for _, userDoc := range endUserEmail {
		_, ok := resultSet[userDoc.ID.String()]

		if !ok {
			resultSet[userDoc.ID.String()] = true
		}
	}

	return
}
