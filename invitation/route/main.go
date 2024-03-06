package route

import (
	"encoding/json"
	"fmt"
	"main/invitation/repository"
	"main/model"
	"main/response"
	"mime/multipart"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Route(route fiber.Router) {
	route.Get("/:name", getInvitationDetail)
	route.Put("/", updateInvitationDetail)
}

func getInvitationDetail(ctx *fiber.Ctx) error {
	name := ctx.Params("name")
	var result *model.Invitation
	result, err := repository.GetInvitationDetailByName(name)
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

	for key := range form.File {
		if err := saveToLocal(ctx, invitation, form, key); err != nil {
			return response.Error(ctx, err)
		}

	}

	setValue("Music", invitation, "aaaa")

	fmt.Println("set string", *invitation.Music)

	if err := repository.UpdateInvitationDetail(invitation); err != nil {
		return response.Error(ctx, err)
	}

	return response.Success(ctx, invitation)
}

func setValue(key string, source interface{}, replace string) {

	reflectKeys := strings.Split(key, ".")
	result := reflect.ValueOf(source)

	for _, reflectKey := range reflectKeys {

		result = reflect.Indirect(result).FieldByName(reflectKey)
		fmt.Println("reflectKey", reflectKey, " result ", result.String(), " result.Kind()", result.Kind())

		if result.Kind() == reflect.Ptr {
			realResult := result.Elem()
			fmt.Println("real value", realResult.String())
			if realResult.Kind() == reflect.String {
				realResult.SetString(replace)
			}
		}

	}
}

func saveToLocal(ctx *fiber.Ctx, invitation model.Invitation, form *multipart.Form, prefix string) error {
	file := form.File[prefix][0]

	folderPath := fmt.Sprintf("./assets/%s", *invitation.Id)

	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
			return err
		}
	} else {
		err := removeCurrentFile(folderPath, prefix)
		if err != nil {
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

func removeCurrentFile(folderPath string, prefix string) error {

	files, err := os.ReadDir(folderPath)
	if err != nil {
		return err
	}

	for _, file := range files {

		fileName := file.Name()
		fmt.Println("fileName", fileName)
		fmt.Println("prefix", prefix)
		if strings.Contains(fileName, prefix) {
			removeFileName := fmt.Sprintf("%s/%s", folderPath, fileName)
			fmt.Println("removeFileName", removeFileName)
			err := os.Remove(removeFileName)
			if err != nil {
				return err
			}
		}

	}
	return nil
}
