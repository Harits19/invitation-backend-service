package repository

import (
	"context"
	"main/common/model"
	"main/common/mongodb"

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
