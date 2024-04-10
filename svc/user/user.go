package user

import (
	"context"
	"secaas_backend/db/doc"
	"secaas_backend/model"
	"secaas_backend/svc/errors"
	"strings"

	"github.com/kamva/mgm/v3"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserSVC struct {
	logger *logrus.Logger
}

func New(logger *logrus.Logger) *UserSVC {
	u := &UserSVC{logger: logger}
	return u
}

func (u *UserSVC) GetByEmail(ctx context.Context, email model.Email) (user model.User, err error) {
	log := u.logger.WithContext(ctx)

	if email == "" {
		log.Error("invalid email")
		err = errors.ErrInvalidEmail
		return
	}

	userDoc := &doc.User{}

	coll := mgm.Coll(userDoc)

	filter := bson.M{
		"email": email.String(),
	}

	err = coll.First(filter, userDoc)

	if err != nil {

		if strings.Contains(err.Error(), "no documents") {
			log.WithError(err).Error("User not found.")
			err = errors.ErrUserNotFound
			return
		}
		log.WithError(err).Error("Unknown error occured when finding user by email.")
		err = errors.ErrUnknown
		return
	}

	user = u.MapDocToUser(userDoc)

	return
}

func (u *UserSVC) CreateUser(ctx context.Context, user model.User) (data model.User, err error) {

	if user.PassHash.Hash == "" || user.PassHash.Alg == "" {
		err = errors.ErrInvalidPassHash
		return
	}

	if user.Email == "" {
		err = errors.ErrInvalidEmail
		return
	}
	docUser := &doc.User{
		Name:  user.Name,
		Email: user.Email,
		PassHash: doc.PassHash{
			Hash: user.PassHash.Hash,
			Alg:  user.PassHash.Alg,
		},
		SymKey: doc.SymKey{
			EncryptedData: user.SymKey.EncryptedData,
			Alg:           user.SymKey.Alg,
		},
		AsymmKey: doc.AsymmKey{
			EncryptedPvtKey: user.AsymmKey.EncryptedPvtKey,
			Public:          user.AsymmKey.Public,
			Alg:             user.AsymmKey.Alg,
		},
		IsBlackListed: false,
		Organization:  []doc.UserOrganization{},
	}

	err = mgm.Coll(docUser).CreateWithCtx(ctx, docUser)

	if err != nil {
		u.logger.WithError(err).Error("error while creating a new user")
		err = errors.ErrUnknown
		return
	}

	data = u.MapDocToUser(docUser)

	return
}

func (u *UserSVC) GetByID(ctx context.Context, id model.UserID) (user model.User, err error) {
	log := u.logger.WithContext(ctx)

	if id == "" {
		log.Error("invalid email")
		err = errors.ErrInvalidEmail
		return
	}

	objId, _ := primitive.ObjectIDFromHex(id.String())

	userDoc := &doc.User{}

	coll := mgm.Coll(userDoc)

	filter := bson.M{
		"_id": objId,
	}

	err = coll.First(filter, userDoc)

	if err != nil {

		if strings.Contains(err.Error(), "no documents") {
			log.WithError(err).Error("User not found.")
			err = errors.ErrUserNotFound
			return
		}
		log.WithError(err).Error("Unknown error occured when finding user by id.")
		err = errors.ErrUnknown
		return
	}

	user = u.MapDocToUser(userDoc)

	return
}

func (u *UserSVC) MapDocToUser(userDoc *doc.User) model.User {
	user := model.User{
		ID:            model.UserID(userDoc.ID.Hex()),
		Name:          userDoc.Name,
		Email:         userDoc.Email,
		PassHash:      model.PassHash(userDoc.PassHash),
		SymKey:        model.SymKey(userDoc.SymKey),
		AsymmKey:      model.AsymmKey(userDoc.AsymmKey),
		CreatedAt:     userDoc.CreatedAt,
		UpdatedAt:     userDoc.UpdatedAt,
		IsBlackListed: userDoc.IsBlackListed,
	}

	if 0 < len(userDoc.Organization) {
		modelOrgs := []model.UserOrganization{}

		for _, org := range userDoc.Organization {

			modelOrg := model.UserOrganization{
				ID: org.ID,
			}

			modelOrgs = append(modelOrgs, modelOrg)

		}

		user.Organization = modelOrgs
	}

	return user

}
