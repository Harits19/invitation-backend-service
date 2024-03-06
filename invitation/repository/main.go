package repository

import (
	"context"
	"main/model"
	"main/mongodb"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func collection(ctx *fiber.Ctx) *mongo.Collection {
	return mongodb.InitConnection(ctx).Collection("invitation-detail")
}

func GetInvitationDetailRepo(ctx *fiber.Ctx, name string) (*model.Invitation, error) {
	var result *model.Invitation

	err := collection(ctx).FindOne(context.Background(), bson.M{"name": name}).Decode(&result)

	return result, err
}

// func updateInvitationDetailRepo() error {

// }
