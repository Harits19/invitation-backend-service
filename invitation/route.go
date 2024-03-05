package invitation

import (
	"context"
	"encoding/json"
	"fmt"
	"main/model"
	"main/mongodb"
	"main/response"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Route(route fiber.Router) {
	route.Get("/:name", getInvitationDetail)
	route.Put("/", updateInvitationDetail)
}

func collection(ctx *fiber.Ctx) *mongo.Collection {
	return mongodb.InitConnection(ctx).Collection("invitation-detail")
}

func getInvitationDetail(ctx *fiber.Ctx) error {
	name := ctx.Params("name")
	var result *model.Invitation
	err := collection(ctx).FindOne(context.Background(), bson.M{"name": name}).Decode(&result)
	if err != nil {
		return response.Error(ctx, err)
	}

	return response.Success(ctx, result)

}

func updateInvitationDetail(ctx *fiber.Ctx) error {
	name := ctx.Params("name")
	fmt.Println("updateInvitationDetail ", name)

	form, err := ctx.MultipartForm()

	if err != nil {
		return response.Error(ctx, err)
	}
	jsonValue := form.Value["json"][0]

	var invitation model.Invitation

	if err := json.Unmarshal([]byte(jsonValue), &invitation); err != nil {
		return response.Error(ctx, err)
	}

	if err := saveToLocal(ctx, invitation, form, "music"); err != nil {
		return response.Error(ctx, err)
	}

	return response.Success(ctx, invitation)
}

func saveToLocal(ctx *fiber.Ctx, invitation model.Invitation, form *multipart.Form, prefix string) error {
	file := form.File[prefix][0]

	folderPath := fmt.Sprintf("./%s", invitation.Id)

	files, err := os.ReadDir(folderPath)
	if err != nil {
		return err
	}

	for _, file := range files {

		fileName := file.Name()

		if strings.Contains(fileName, prefix) {
			os.Remove(fmt.Sprintf("%s/%s", folderPath, fileName))
		}

	}

	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		if err := os.Mkdir(folderPath, os.ModePerm); err != nil {
			return err
		}
	}

	fileName := strings.Split(file.Filename, ".")
	fileType := fileName[len(fileName)-1]

	id := time.Now().Unix()

	if err := ctx.SaveFile(file, fmt.Sprintf("%s/%s_%d.%s", folderPath, prefix, id, fileType)); err != nil {
		return err
	}
	return nil
}
