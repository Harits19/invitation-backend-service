package route

import (
	"encoding/json"
	"fmt"
	"main/common/bucket"
	"main/common/model"
	"main/common/response"
	"main/common/util"
	"main/invitation/repository"
	"mime/multipart"
	"net/http"
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gofiber/fiber/v2"
)

func Route(route fiber.Router) {
	route.Get("/:name", getInvitationDetail)
	route.Patch("/:name", updateInvitationDetail)
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

	bucket := bucket.New(fmt.Sprint(*invitation.Id))
	if err := bucket.CreateBucket(); err != nil {
		return response.Error(ctx, err, nil)
	}

	bucketFiles, err := bucket.GetListOfBucketFile()

	if err != nil {
		return response.Error(ctx, err, nil)
	}

	for key, files := range form.File {
		if len(files) == 0 {
			continue
		}

		err := setInvitationValue(bucket, bucketFiles, invitation, key, files[0])

		if err != nil {
			return response.Error(ctx, err, nil)

		}
	}

	err = repository.UpdateInvitationDetail(invitation)

	if err != nil {
		return response.Error(ctx, err, nil)

	}

	return response.Success(ctx, invitation)
}

func setInvitationValue(bucket bucket.Bucket, files *s3.ListObjectsV2Output, invitation model.Invitation, prefix string, file *multipart.FileHeader) error {
	prefix = util.TitleCase(prefix)
	prefix, index := util.GetRealKey(prefix)

	keys := strings.Split(prefix, ".")

	r := reflect.ValueOf(invitation)

	url, err := bucket.SaveToStorage(file, prefix, index)
	if err != nil {
		return err
	}

	for _, key := range keys {

		value := r.FieldByName(key)

		if value.Kind() == reflect.Ptr {
			value = value.Elem()
		}

		fmt.Println("key", key, value.Kind())

		if value.Kind() == reflect.Struct {
			r = value
			continue
		}

		if value.Kind() == reflect.String {
			value.Set(reflect.ValueOf(url))
			break
		}

		if value.Kind() == reflect.Slice {
			value = value.Index(index)
			if value.CanSet() {
				err := bucket.FindAndDelete(files, prefix, index)
				if err != nil {
					return err
				}
				value.Set(reflect.ValueOf(url))
			}
		}

	}

	return nil

}
