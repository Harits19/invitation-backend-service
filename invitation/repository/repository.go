package repository

import (
	"context"
	"errors"
	"main/common/model"
	"main/common/mongodb"
	"main/common/util"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func collection() *mongo.Collection {
	return mongodb.DB.Collection("invitation-detail")
}

func GetInvitationDetailByName(name string) (*model.Invitation, error) {
	var result *model.Invitation

	err := collection().FindOne(context.Background(), bson.M{"name": name}).Decode(&result)

	return result, err
}

func UpdateInvitationDetail(invitation model.Invitation) error {
	invitation.Id = nil
	_, err := collection().UpdateOne(context.Background(), bson.M{"name": invitation.Name}, bson.M{"$set": invitation})

	return err
}

func CreateInvitation(invitation model.Invitation) error {
	var resultFind *model.Invitation
	collection().FindOne(context.Background(), bson.M{"name": invitation.Name}).Decode(&resultFind)

	util.Log(resultFind)

	if resultFind != nil {
		return errors.New("invitation name already exist")
	}
	_, err := collection().InsertOne(context.Background(), invitation)
	return err
}
