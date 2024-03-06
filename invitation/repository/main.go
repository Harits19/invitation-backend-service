package repository

import (
	"context"
	"main/model"
	"main/mongodb"

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

// func updateInvitationDetail(ctx *fiber.Ctx) error {
// 	err := collection(ctx).FindOne(context.Background(), bson.M{"name": name}).Decode(&result)

// }
