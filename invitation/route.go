package invitation

import (
	"context"
	"main/model"
	"main/mongodb"
	"main/response"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Route(app *fiber.App) {
	route := app.Group("/invitation")
	route.Get("/:name", getInvitationDetail)
}

func collection(ctx *fiber.Ctx) *mongo.Collection {
	return mongodb.InitConnection(ctx).Collection("invitation-detail")
}

func getInvitationDetail(ctx *fiber.Ctx) error {
	name := ctx.Params("name")
	var result model.Invitation
	err := collection(ctx).FindOne(context.Background(), bson.M{"name": name}).Decode(&result)
	if err != nil {
		return err
	}
	return response.Success(ctx, result)

}
