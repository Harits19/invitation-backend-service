package route

import (
	"encoding/json"
	"fmt"
	"main/common/model"
	"main/common/response"
	"main/common/util"
	"main/invitation/repository"

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

		paths, err := saveToLocal(ctx, invitation, form, key)

		for _, path := range paths {
			if err != nil {
				return response.Error(ctx, err)
			}

			key = util.TitleCase(key)
			setValue(key, invitation, path)
		}

	}
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

		if result.Kind() == reflect.Ptr {
			realResult := result.Elem()
			if realResult.Kind() == reflect.String {
				realResult.SetString(replace)
			} else if realResult.Kind() == reflect.Slice {
				newValue := reflect.Append(realResult, reflect.ValueOf(replace))
				realResult.Set(newValue)
			}

		}

	}
}

func saveToLocal(ctx *fiber.Ctx, invitation model.Invitation, form *multipart.Form, prefix string) ([]string, error) {
	filePathList := []string{}
	for index, file := range form.File[prefix] {

		folderPath := fmt.Sprintf("./assets/%s", *invitation.Id)
		newPrefix := fmt.Sprintf("%s%d", prefix, index)
		if _, err := os.Stat(folderPath); os.IsNotExist(err) {
			if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
				return filePathList, err
			}
		} else {
			err := removeCurrentFile(folderPath, newPrefix)
			if err != nil {
				return filePathList, err
			}
		}

		fileName := strings.Split(file.Filename, ".")
		fileType := fileName[len(fileName)-1]

		id := time.Now().Unix()

		filePath := fmt.Sprintf("%s/%s_%d.%s", folderPath, newPrefix, id, fileType)

		if err := ctx.SaveFile(file, filePath); err != nil {
			return filePathList, err
		}
		filePath = filePath[1:]
		filePathList = append(filePathList, filePath)
	}

	return filePathList, nil

}

func removeCurrentFile(folderPath string, prefix string) error {

	files, err := os.ReadDir(folderPath)
	if err != nil {
		return err
	}

	for _, file := range files {

		fileName := file.Name()
		if strings.Contains(fileName, prefix) {
			removeFileName := fmt.Sprintf("%s/%s", folderPath, fileName)
			err := os.Remove(removeFileName)
			if err != nil {
				return err
			}
		}

	}
	return nil
}
