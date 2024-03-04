package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitConnection(ctx *fiber.Ctx) *mongo.Database {
	fmt.Println("init connection")
	context, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	uri := "mongodb://root:example@localhost:27017/"
	dbName := "invitation"
	client, err := mongo.Connect(context, options.Client().ApplyURI(uri))
	if err != nil {
		ctx.JSON(fiber.Map{"error": err})
		return nil
	}
	InvitationDb := client.Database(dbName)
	return InvitationDb
}
