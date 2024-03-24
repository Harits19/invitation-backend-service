package route

import (
	"encoding/json"
	"fmt"
	"main/common/model"
	"main/common/response"
	"main/common/s3"
	"main/common/util"
	"main/invitation/repository"
	"net/http"
	"reflect"
	"strings"

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

	for key, files := range form.File {
		if len(files) == 0 {
			continue
		}
		file := files[0]
		fmt.Println(key, file.Filename)

		setInvitationValue(invitation, key)
	}

	return response.Success(ctx, invitation)
}

func setInvitationValue(invitation model.Invitation, prefix string) {
	prefix = util.TitleCase(prefix)

	prefix, index := util.GetRealKey(prefix)
	keys := strings.Split(prefix, ".")

	r := reflect.ValueOf(invitation)

	for _, key := range keys {
		value := r.FieldByName(key)

		if value.Kind() == reflect.Ptr {
			value = value.Elem()
		}

		if value.Kind() == reflect.Struct {
			r = value
			continue
		}
		fmt.Println("key", key, value.Kind())

		if value.Kind() == reflect.String {
			fmt.Println("set value key", prefix, value.String())
			break
		}

		if value.Kind() == reflect.Slice {
			if value.Len() == 0 {
				continue
			}
			value = value.Index(index)

			if value.Kind() == reflect.String {
				fmt.Println("set value key", prefix, value.String())
				break
			}

		}

	}

}
