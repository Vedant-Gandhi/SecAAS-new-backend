package invite

import (
	"context"
	"secaas_backend/db/doc"
	"secaas_backend/model"
	"secaas_backend/svc/errors"
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
		ExpiresAt:      time.Now(),
		ToUserEmail:    data.ToUserEmail,
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
		{"createdAt", 1},
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
		{"updatedAt", 1},
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

func (u *InviteSVC) MapDocToInvite(userInvite *doc.Invite) model.Invite {
	user := model.Invite{
		ID:             model.InviteID(userInvite.ID.Hex()),
		CreatedAt:      userInvite.CreatedAt,
		ExpiresAt:      userInvite.ExpiresAt,
		UpdatedAt:      userInvite.UpdatedAt,
		FromUserEmail:  userInvite.FromUserEmail,
		ToUserEmail:    userInvite.ToUserEmail,
		OrganizationID: userInvite.OrganizationID,
	}

	return user

}
