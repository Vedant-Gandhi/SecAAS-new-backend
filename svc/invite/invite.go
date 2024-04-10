package invite

import (
	"context"
	"secaas_backend/db/doc"
	"secaas_backend/model"
	"secaas_backend/svc/errors"
	"strings"
	"time"

	"github.com/kamva/mgm/v3"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type InviteSVC struct {
	logger *logrus.Logger
}

func New(logger *logrus.Logger) *InviteSVC {
	u := &InviteSVC{logger: logger}
	return u
}

func (s *InviteSVC) CreateInvite(ctx context.Context, data model.Invite) (model.Invite, error) {

	docInvite := &doc.Invite{
		FromUserEmail:  data.FromUserEmail,
		OrganizationID: data.OrganizationID,
		ExpiresAt:      data.ExpiresAt,
		ToUserEmail:    data.ToUserEmail,
		SymKey: doc.SymKey{
			EncryptedData: data.SymKey.EncryptedData,
			Alg:           data.SymKey.Alg,
		},
	}

	err := mgm.Coll(docInvite).CreateWithCtx(ctx, docInvite)

	if err != nil {
		s.logger.WithError(err).Error("error while creating a new user")
		err = errors.ErrUnknown
		return model.Invite{}, err
	}

	newInvite := s.MapDocToInvite(docInvite)

	return newInvite, nil
}

func (s *InviteSVC) GetInvitesByOrganization(ctx context.Context, orgId string, params model.PaginationParams, onlyActive bool) ([]model.Invite, error) {
	filter := bson.M{
		"organizationId": orgId,
	}

	// Get only active ones
	if onlyActive {
		filter["expiresAt"] = bson.M{"$gt": time.Now()}
	}

	inviteDoc := &doc.Invite{}

	coll := mgm.Coll(inviteDoc)

	findOptions := options.Find().SetSkip(int64(params.Skip)).SetLimit(int64(params.Limit)).SetSort(bson.D{
		{"createdAt", -1},
	})

	rawDocs, err := coll.Find(ctx, filter, findOptions)

	if err != nil {
		s.logger.WithContext(ctx).WithError(err).Error("error while fetching list of invites")
		err = errors.ErrUnknown
		return []model.Invite{}, err
	}

	defer rawDocs.Close(ctx)

	var docs []model.Invite = make([]model.Invite, 0)

	for rawDocs.Next(ctx) {
		var curDoc doc.Invite

		err = rawDocs.Decode(&curDoc)

		if err != nil {
			s.logger.WithContext(ctx).WithError(err).Error("error while decoding invite doc.")
			continue
		}

		modelInvite := s.MapDocToInvite(&curDoc)
		docs = append(docs, modelInvite)

	}

	s.DeleteInvite(ctx, inviteDoc.ID.Hex())

	return docs, nil
}

func (s *InviteSVC) GetInvitesForUser(ctx context.Context, receiverEmail string, params model.PaginationParams, onlyActive bool) ([]model.Invite, error) {
	filter := bson.M{
		"toUserEmail": receiverEmail,
	}

	// Get only active ones
	if onlyActive {
		filter["expiresAt"] = bson.M{"$gt": time.Now()}
	}

	inviteDoc := &doc.Invite{}

	coll := mgm.Coll(inviteDoc)

	findOptions := options.Find().SetSkip(int64(params.Skip)).SetLimit(int64(params.Limit)).SetSort(bson.D{
		{"updatedAt", -1},
	})

	rawDocs, err := coll.Find(ctx, filter, findOptions)

	if err != nil {
		s.logger.WithContext(ctx).WithError(err).Error("error while fetching list of invites for user.")
		err = errors.ErrUnknown
		return []model.Invite{}, err
	}

	defer rawDocs.Close(ctx)

	var docs []model.Invite = make([]model.Invite, 0)

	for rawDocs.Next(ctx) {
		var curDoc doc.Invite

		err = rawDocs.Decode(&curDoc)

		if err != nil {
			s.logger.WithContext(ctx).WithError(err).Error("error while decoding invite doc for user.")
			err = nil
			continue
		}

		modelInvite := s.MapDocToInvite(&curDoc)
		docs = append(docs, modelInvite)

	}

	return docs, nil
}

func (s *InviteSVC) DeleteInvite(ctx context.Context, inviteId string) (int, error) {
	objId, err := primitive.ObjectIDFromHex(inviteId)

	if err != nil {
		s.logger.WithContext(ctx).WithError(err).Error("Invalid invite id found")
		err = errors.ErrInvalidID
		return 0, err
	}

	inviteDoc := doc.Invite{}

	filter := bson.M{
		"_id": objId,
	}

	res, err := mgm.Coll(&inviteDoc).DeleteOne(ctx, filter)

	if err != nil {
		s.logger.WithContext(ctx).WithError(err).Error("Error while deleteing invitation")
		err = errors.ErrUnknown
		return 0, err
	}

	return int(res.DeletedCount), nil

}

func (i *InviteSVC) AcceptInvite(ctx context.Context, inviteId string) (err error) {
	docInvite := &doc.Invite{}

	objId, err := primitive.ObjectIDFromHex(inviteId)

	if err != nil {
		i.logger.WithContext(ctx).WithError(err).Error("failed to convert invite id to object id")
		err = errors.ErrInvalidID
		return
	}

	filter := bson.M{
		"_id":       objId,
		"expiresAt": bson.M{"$gt": time.Now()},
	}

	err = mgm.Coll(docInvite).First(filter, docInvite)

	if err != nil {

		if strings.Contains(err.Error(), "no documents") {
			i.logger.WithContext(ctx).WithError(err).Error("Invite not found.")
			err = errors.ErrInviteNotFound
			return
		}
		i.logger.WithContext(ctx).WithError(err).Error("Error while accepting invitation")
		err = errors.ErrUnknown
		return
	}

	userDoc := &doc.User{}

	userFilter := bson.M{
		"email": docInvite.ToUserEmail,
	}

	userUpdate := bson.M{
		"$push": bson.M{
			"organizations": doc.UserOrganization{
				ID:      docInvite.OrganizationID,
				IsAdmin: false,
				PvtKey:  docInvite.SymKey.EncryptedData,
			},
		},
	}

	updateRes, err := mgm.Coll(userDoc).UpdateOne(ctx, userFilter, userUpdate)

	if err != nil {
		i.logger.WithContext(ctx).WithField("To Email", docInvite.ToUserEmail).WithError(err).Error("Failed to add the organization entry in the admin when accepting invite.")
		err = nil
	}

	if updateRes.ModifiedCount == 0 {
		i.logger.WithContext(ctx).WithField("To Email", docInvite.ToUserEmail).Print("Could not update the email with organization when accepting invite.")
	} else {
		i.logger.WithContext(ctx).WithField("To Email", docInvite.ToUserEmail).Print("Added organization to the email when accepting invite.")
	}

	return

}

func (u *InviteSVC) MapDocToInvite(userInvite *doc.Invite) model.Invite {
	user := model.Invite{
		ID:               model.InviteID(userInvite.ID.Hex()),
		CreatedAt:        userInvite.CreatedAt,
		ExpiresAt:        userInvite.ExpiresAt,
		UpdatedAt:        userInvite.UpdatedAt,
		FromUserEmail:    userInvite.FromUserEmail,
		ToUserEmail:      userInvite.ToUserEmail,
		OrganizationID:   userInvite.OrganizationID,
		OrganizationName: userInvite.OrganizationName,
		SymKey:           model.SymKey(userInvite.SymKey),
	}

	return user

}
