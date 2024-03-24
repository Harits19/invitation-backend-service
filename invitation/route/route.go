package route

import (
	"encoding/json"
	"errors"
	"fmt"
	"main/common/model"
	"main/common/response"
	"main/common/s3"
	"main/common/util"
	"main/invitation/repository"
	"net/http"
	"strconv"

	"mime/multipart"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Route(route fiber.Router) {
	route.Get("/:name", getInvitationDetail)
	route.Put("/", updateInvitationDetail)
	route.Post("/", createInvitation)
}

func createInvitation(ctx *fiber.Ctx) (err error) {
	body := ctx.Body()
	var invitation model.Invitation

	if err = json.Unmarshal(body, &invitation); err != nil {
		return response.Error(ctx, err, nil)
	}

	util.Stringify(invitation)
	err = repository.CreateInvitation(invitation)

	if err != nil {
		return response.Error(ctx, err, nil)
	}

	return response.Success(ctx, "success create invitation ")

}

func getInvitationDetail(ctx *fiber.Ctx) error {
	name := ctx.Params("name")
	var result *model.Invitation
	result, err := repository.GetInvitationDetailByName(name)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			notFound := http.StatusNotFound
			return response.Error(ctx, err, &notFound)
		}
		return response.Error(ctx, err, nil)
	}

	return response.Success(ctx, result)

}

func updateInvitationDetail(ctx *fiber.Ctx) error {
	name := ctx.Params("name")
	fmt.Println("updateInvitationDetail ", name)

	form, err := ctx.MultipartForm()

	if err != nil {
		return response.Error(ctx, err, nil)
	}
	jsonValue := form.Value["json"][0]

	var invitation model.Invitation

	if err := json.Unmarshal([]byte(jsonValue), &invitation); err != nil {
		return response.Error(ctx, err, nil)
	}

	bucket := s3.New(fmt.Sprint(*invitation.Id))
	if err := bucket.CreateBucket(); err != nil {
		return response.Error(ctx, err, nil)
	}

	for prefix, file := range form.File {
		saveToStorage(bucket, file, prefix)
	}
	return response.Success(ctx, invitation)
}

func saveToStorage(bucket s3.Bucket, fileHeader []*multipart.FileHeader, prefix string) (string, error) {
	if len(fileHeader) == 0 {
		return "", errors.New("file header length = 0")
	}

	file := fileHeader[0]
	fileName := strings.Split(file.Filename, ".")
	fileExtension := fileName[len(fileName)-1]

	prefix, index := getRealKey(prefix)

	newFileName := fmt.Sprintf("%s/%s%d", "assets", prefix, index)

	uniqueId := time.Now().Unix()

	filePath := fmt.Sprintf("%s_%d.%s", newFileName, uniqueId, fileExtension)

	url, err := bucket.UploadFile(filePath, *file)

	return *url, err
}

func getRealKey(key string) (string, int) {
	splitKey := strings.Split(key, ".")
	lastIndex := len(splitKey) - 1
	lastKey := splitKey[lastIndex]

	lasKeyIndex, err := strconv.Atoi(lastKey)

	if err != nil {
		return key, 0
	}

	newKey := strings.Join(removeIndex(splitKey, lastIndex), "")

	return newKey, lasKeyIndex
}

func removeIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}
